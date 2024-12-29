// Package async provides the business logic for our async tasks.
package async

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid/publisher"
	"github.com/RMI/pacta/async/parsed"
	"github.com/RMI/pacta/blob"
	"github.com/RMI/pacta/pacta"
	"github.com/RMI/pacta/task"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapfield"
)

type Config struct {
	Blob   Blob
	PubSub *publisher.Client
	Logger *zap.Logger

	BenchmarkDir string
	PACTADataDir string
}

func (c *Config) validate() error {
	if c.Blob == nil {
		return errors.New("no blob client given")
	}
	if c.PubSub == nil {
		return errors.New("no pub/sub client given")
	}
	if c.Logger == nil {
		return errors.New("no logger given")
	}

	if c.BenchmarkDir == "" {
		return errors.New("no benchmark dir specified")
	}
	if c.PACTADataDir == "" {
		return errors.New("no PACTA data dir specified")
	}
	return nil
}

type Blob interface {
	ReadBlob(ctx context.Context, uri string) (io.ReadCloser, error)
	WriteBlob(ctx context.Context, uri string, r io.Reader) error
	Scheme() blob.Scheme
}

type Handler struct {
	blob   Blob
	pubsub *publisher.Client
	logger *zap.Logger

	// Mounted directories with data needed for report generation.
	benchmarkDir, pactaDataDir string
}

func New(cfg *Config) (*Handler, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config given: %w", err)
	}

	return &Handler{
		blob:         cfg.Blob,
		pubsub:       cfg.PubSub,
		logger:       cfg.Logger,
		benchmarkDir: cfg.BenchmarkDir,
		pactaDataDir: cfg.PACTADataDir,
	}, nil
}

// TODO: Send a notification when parsing fails.
func (h *Handler) ParsePortfolio(ctx context.Context, taskID task.ID, req *task.ParsePortfolioRequest, destPortfolioContainer string) error {
	// Make the directories we require first. We use these instead of
	// /mnt/{input,output} because the base image (quite reasonably) uses a non-root
	// user, so we can't be creating directories in the root filesystem all willy
	// nilly.
	inputDir := filepath.Join("/", "home", "portfolio-parser", "input")
	outputDir := filepath.Join("/", "home", "portfolio-parser", "output")

	if err := os.MkdirAll(inputDir, 0700); err != nil {
		return fmt.Errorf("failed to create input dir to store input CSVs: %w", err)
	}
	if err := os.MkdirAll(outputDir, 0700); err != nil {
		return fmt.Errorf("failed to create output dir to store output CSVs: %w", err)
	}

	// Load the portfolio from blob storage, place it in /mnt/inputs, where
	// the `process_portfolios.R` script expects it to be.
	localCSVToBlob := make(map[string]pacta.BlobURI)
	for _, srcURI := range req.BlobURIs {
		id := uuid.New().String()
		// TODO: Probably set the CSV extension in the signed upload URL instead.
		fn := fmt.Sprintf("%s.csv", id)
		destPath := filepath.Join(inputDir, fn)
		if err := h.downloadBlob(ctx, string(srcURI), destPath); err != nil {
			return fmt.Errorf("failed to download raw portfolio blob: %w", err)
		}
		localCSVToBlob[fn] = srcURI
	}

	cmd := exec.CommandContext(ctx,
		"/usr/local/bin/Rscript",
		"-e", "logger::log_threshold(Sys.getenv('LOG_LEVEL', 'INFO'));workflow.portfolio.parsing::process_directory('"+inputDir+"', '"+outputDir+"')",
	)

	// We don't expect log output to be particularly large, it's fine to write them to an in-memory buffer.
	// TODO(#185): Find a good place to put these in storage, such that it can be correlated with the input file(s)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run process_portfolios script: %w", err)
	}

	// After successful execution, the API contract is that there should be a 'processed_portfolios.json' file in the output directory.
	outManifestPath := filepath.Join(outputDir, "processed_portfolios.json")
	omf, err := os.Open(outManifestPath)
	if err != nil {
		return fmt.Errorf("failed to open output processed_portfolios.json file: %w", err)
	}
	defer omf.Close()

	var sourceFiles []parsed.SourceFile
	if err := json.NewDecoder(omf).Decode(&sourceFiles); err != nil {
		return fmt.Errorf("failed to decode processed_portfolios.json as JSON: %w", err)
	}

	// We keep track of the outputs we processed, and then check if there are any files in the output directory that we weren't expecting.
	knownOutputFiles := map[string]bool{
		"processed_portfolios.json": true,
	}

	// NOTE: This code could benefit from some concurrency, but I'm opting not to prematurely optimize.
	var out []*task.ParsePortfolioResponseItem
	for _, sf := range sourceFiles {
		sourceURI, ok := localCSVToBlob[sf.InputFilename]
		if !ok {
			return fmt.Errorf("parse output mentioned input file %q, which wasn't found in our input -> blob URI map %+v", sf.InputFilename, localCSVToBlob)
		}

		// TODO(#187): There's lots of metadata associated with the input files (e.g.
		// sf.Errors, sf.GroupCols, etc), we should likely store that info somewhere.

		for _, p := range sf.Portfolios {
			outPath := filepath.Join(outputDir, p.OutputFilename)

			// We generate a fresh UUID here for uploading the file to blob storage, so that
			// we don't depend on the R code generating truly unique UUIDs.
			uploadName := fmt.Sprintf("%s.csv", uuid.New().String())

			blobURI := pacta.BlobURI(blob.Join(h.blob.Scheme(), destPortfolioContainer, uploadName))

			if err := h.uploadBlob(ctx, outPath, string(blobURI)); err != nil {
				return fmt.Errorf("failed to copy parsed portfolio from %q to %q: %w", p, blobURI, err)
			}
			h.logger.Info("uploaded output CSV to blob storage", zap.Any("portfolio", p), zapfield.Str("blob_uri", blobURI))

			extension := filepath.Ext(p.OutputFilename)
			fileType, err := pacta.ParseFileType(extension)
			if err != nil {
				return fmt.Errorf("failed to parse file type from file name %q: %w", p.OutputFilename, err)
			}
			if fileType != pacta.FileType_CSV {
				return fmt.Errorf("output portfolio %q was not of type CSV, was %q", p.OutputFilename, fileType)
			}

			knownOutputFiles[p.OutputFilename] = true
			out = append(out, &task.ParsePortfolioResponseItem{
				Source: sourceURI,
				Blob: pacta.Blob{
					FileName: p.OutputFilename,
					FileType: fileType,
					BlobURI:  blobURI,
				},
				Portfolio: p,
			})
		}
	}

	// Now that we're done uploading files, check the output directory and make sure
	// there aren't any unaccounted for files.
	dirEntries, err := os.ReadDir(outputDir)
	if err != nil {
		return fmt.Errorf("failed to read the output directory: %w", err)
	}
	for _, de := range dirEntries {
		if !knownOutputFiles[de.Name()] {
			h.logger.Error("output directory contained files not present in the generated 'processed_portfolios.json' manifest", zap.String("filename", de.Name()))
		}
	}

	events := []publisher.Event{
		{
			Data: task.ParsePortfolioResponse{
				TaskID:  taskID,
				Request: req,
				Outputs: out,
			},
			DataVersion: to.Ptr("1.0"),
			EventType:   to.Ptr("parsed-portfolio"),
			EventTime:   to.Ptr(time.Now()),
			ID:          to.Ptr(string(taskID)),
			Subject:     to.Ptr(string(taskID)),
		},
	}

	if _, err := h.pubsub.PublishEvents(ctx, events, nil); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	h.logger.Info("parsed portfolio", zap.String("task_id", string(taskID)))

	return nil
}

func (h *Handler) CreateAudit(ctx context.Context, taskID task.ID, req *task.CreateAuditRequest, auditContainer string) error {
	if n := len(req.BlobURIs); n != 1 {
		return fmt.Errorf("expected exactly one blob URI as input, got %d", n)
	}
	blobURI := req.BlobURIs[0]

	//  We use this instead of /mnt/... because the base image (quite
	// reasonably) uses a non-root user, so we can't be creating directories in the
	// root filesystem all willy nilly.
	baseDir := filepath.Join("/", "home", "workflow-pacta-webapp")

	// We don't use the benchmark or PACTA data here, it's just convenient to use the same directory-creating harness for everything.
	auditEnv, err := initEnv(h.benchmarkDir, h.pactaDataDir, baseDir, "create-audit")
	if err != nil {
		return fmt.Errorf("failed to init report env: %w", err)
	}

	// Load the parsed portfolio from blob storage, place it in our PORFOLIO_DIR,
	// where the `run_pacta.R` script expects it to be.
	fileNameWithExt := filepath.Base(string(blobURI))
	if !strings.HasSuffix(fileNameWithExt, ".csv") {
		return fmt.Errorf("given blob wasn't a CSV-formatted portfolio, %q", fileNameWithExt)
	}
	destPath := filepath.Join(auditEnv.pathForDir(PortfoliosDir), fileNameWithExt)
	if err := h.downloadBlob(ctx, string(blobURI), destPath); err != nil {
		return fmt.Errorf("failed to download processed portfolio blob: %w", err)
	}

	inp := AuditInput{
		Portfolio: AuditInputPortfolio{
			Files:        []string{fileNameWithExt},
			HoldingsDate: "2023-12-31",   // TODO(#206)
			Name:         "FooPortfolio", // TODO(#206)
		},
		Inherit: "GENERAL_2023Q4", // TODO(#206): Should this be configurable
	}

	var inpJSON bytes.Buffer
	if err := json.NewEncoder(&inpJSON).Encode(inp); err != nil {
		return fmt.Errorf("failed to encode audit input as JSON: %w", err)
	}

	cmd := exec.CommandContext(ctx, "/workflow.pacta.webapp/inst/extdata/scripts/run_audit.sh", inpJSON.String())
	cmd.Env = append(cmd.Env, auditEnv.asEnvVars()...)
	cmd.Env = append(cmd.Env,
		"LOG_LEVEL=DEBUG",
		"HOME=/root", /* Required by pandoc */
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run pacta test CLI: %w", err)
	}

	var artifacts []*task.AnalysisArtifact
	uploadDir := func(dir string) error {
		aas, err := h.uploadDirectory(ctx, dir, auditContainer, req.AnalysisID)
		if err != nil {
			return fmt.Errorf("failed to upload report directory: %w", err)
		}
		artifacts = append(artifacts, aas...)
		return nil
	}

	for _, outDir := range auditEnv.outputDirs() {
		if err := uploadDir(outDir); err != nil {
			return fmt.Errorf("failed to upload artifacts %q: %w", outDir, err)
		}
	}

	events := []publisher.Event{
		{
			Data: task.CreateAuditResponse{
				TaskID:    taskID,
				Request:   req,
				Artifacts: artifacts,
			},
			DataVersion: to.Ptr("1.0"),
			EventType:   to.Ptr("created-audit"),
			EventTime:   to.Ptr(time.Now()),
			ID:          to.Ptr(string(taskID)),
			Subject:     to.Ptr(string(taskID)),
		},
	}

	if _, err := h.pubsub.PublishEvents(ctx, events, nil); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	h.logger.Info("created report", zap.String("task_id", string(taskID)))

	return nil
}

type ReportInput struct {
	Portfolio ReportInputPortfolio `json:"portfolio"`
	Inherit   string               `json:"inherit"`
}

type ReportInputPortfolio struct {
	Files        []string `json:"files"`
	HoldingsDate string   `json:"holdingsDate"`
	Name         string   `json:"name"`
}

type AuditInput struct {
	Portfolio AuditInputPortfolio `json:"portfolio"`
	Inherit   string              `json:"inherit"`
}

type AuditInputPortfolio struct {
	Files        []string `json:"files"`
	HoldingsDate string   `json:"holdingsDate"`
	Name         string   `json:"name"`
}

type DashboardInput struct {
	Portfolio DashboardInputPortfolio `json:"portfolio"`
	Inherit   string                  `json:"inherit"`
}

type DashboardInputPortfolio struct {
	Files        []string `json:"files"`
	HoldingsDate string   `json:"holdingsDate"`
	Name         string   `json:"name"`
}

type TaskEnv struct {
	rootDir string

	// These are mounted in from externally.
	benchmarksDir string
	pactaDataDir  string
}

func initEnv(benchmarkDir, pactaDataDir, baseDir, taskName string) (*TaskEnv, error) {
	// Make sure the base directory exists first.
	if err := os.MkdirAll(baseDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create base input dir: %w", err)
	}
	// We create temp subdirectories, because while this code currently executes in
	// a new container for each invocation, that might not always be the case.
	rootDir, err := os.MkdirTemp(baseDir, taskName)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir for input CSVs: %w", err)
	}

	re := &TaskEnv{
		rootDir:       rootDir,
		benchmarksDir: benchmarkDir,
		pactaDataDir:  pactaDataDir,
	}

	if err := re.makeDirectories(); err != nil {
		return nil, fmt.Errorf("failed to create directories: %w", err)
	}

	return re, nil
}

type ReportDir string

const (
	PortfoliosDir = ReportDir("portfolios")
	RealEstateDir = ReportDir("real-estate")
	ScoreCardDir  = ReportDir("score-card")
	SurveyDir     = ReportDir("survey")

	// Outputs
	AnalysisOutputDir  = ReportDir("analysis-output")
	ReportOutputDir    = ReportDir("report-output")
	DashboardOutputDir = ReportDir("dashboard-output")
	SummaryOutputDir   = ReportDir("summary-output")
)

func (r *TaskEnv) outputDirs() []string {
	return []string{
		r.pathForDir(AnalysisOutputDir),
		r.pathForDir(ReportOutputDir),
		r.pathForDir(DashboardOutputDir),
		r.pathForDir(SummaryOutputDir),
	}
}

func (r *TaskEnv) asEnvVars() []string {
	return []string{
		"BENCHMARKS_DIR=" + r.benchmarksDir,
		"PACTA_DATA_DIR=" + r.pactaDataDir,
		"PORTFOLIO_DIR=" + r.pathForDir(PortfoliosDir),
		"REAL_ESTATE_DIR=" + r.pathForDir(RealEstateDir),
		"SCORE_CARD_DIR=" + r.pathForDir(ScoreCardDir),
		"SURVEY_DIR=" + r.pathForDir(SurveyDir),
		"ANALYSIS_OUTPUT_DIR=" + r.pathForDir(AnalysisOutputDir),
		"REPORT_OUTPUT_DIR=" + r.pathForDir(ReportOutputDir),
		"DASHBOARD_OUTPUT_DIR=" + r.pathForDir(DashboardOutputDir),
		"SUMMARY_OUTPUT_DIR=" + r.pathForDir(SummaryOutputDir),
	}
}

func (r *TaskEnv) pathForDir(d ReportDir) string {
	return filepath.Join(r.rootDir, string(d))
}

func (r *TaskEnv) makeDirectories() error {
	var rErr error
	makeDir := func(reportDir ReportDir) {
		if rErr != nil {
			return
		}
		dir := r.pathForDir(reportDir)
		if err := os.Mkdir(dir, 0700); err != nil {
			rErr = fmt.Errorf("failed to create dir %q: %w", dir, err)
			return
		}
	}

	// Inputs
	makeDir(PortfoliosDir)
	makeDir(RealEstateDir) // Used as part of specific projects, empty for now.
	makeDir(ScoreCardDir)  // Used as part of specific projects, empty for now.
	makeDir(SurveyDir)     // Used as part of specific projects, empty for now.

	// Outputs
	makeDir(AnalysisOutputDir)
	makeDir(DashboardOutputDir)
	makeDir(ReportOutputDir)
	makeDir(SummaryOutputDir)

	if rErr != nil {
		return rErr
	}
	return nil
}

func (h *Handler) CreateDashboard(ctx context.Context, taskID task.ID, req *task.CreateDashboardRequest, dashboardContainer string) error {
	if n := len(req.BlobURIs); n != 1 {
		return fmt.Errorf("expected exactly one blob URI as input, got %d", n)
	}
	blobURI := req.BlobURIs[0]

	//  We use this instead of /mnt/... because the base image (quite
	// reasonably) uses a non-root user, so we can't be creating directories in the
	// root filesystem all willy nilly.
	baseDir := filepath.Join("/", "home", "workflow-pacta-webapp")

	dashEnv, err := initEnv(h.benchmarkDir, h.pactaDataDir, baseDir, "create-dashboard")
	if err != nil {
		return fmt.Errorf("failed to init report env: %w", err)
	}

	// Load the parsed portfolio from blob storage, place it in our PORFOLIO_DIR,
	// where the `prepare_dashboard_data.R` script expects it to be.
	fileNameWithExt := filepath.Base(string(blobURI))
	if !strings.HasSuffix(fileNameWithExt, ".csv") {
		return fmt.Errorf("given blob wasn't a CSV-formatted portfolio, %q", fileNameWithExt)
	}
	destPath := filepath.Join(dashEnv.pathForDir(PortfoliosDir), fileNameWithExt)
	if err := h.downloadBlob(ctx, string(blobURI), destPath); err != nil {
		return fmt.Errorf("failed to download processed portfolio blob: %w", err)
	}

	inp := DashboardInput{
		Portfolio: DashboardInputPortfolio{
			Files:        []string{fileNameWithExt},
			HoldingsDate: "2023-12-31",   // TODO(#206)
			Name:         "FooPortfolio", // TODO(#206)
		},
		Inherit: "GENERAL_2023Q4", // TODO(#206): Should this be configurable
	}

	var inpJSON bytes.Buffer
	if err := json.NewEncoder(&inpJSON).Encode(inp); err != nil {
		return fmt.Errorf("failed to encode report input as JSON: %w", err)
	}

	cmd := exec.CommandContext(ctx,
		"/usr/local/bin/Rscript",
		"--vanilla", "/workflow.pacta.dashboard/inst/extdata/scripts/prepare_dashboard_data.R",
		inpJSON.String())

	cmd.Env = append(cmd.Env, dashEnv.asEnvVars()...)
	cmd.Env = append(cmd.Env,
		"LOG_LEVEL=DEBUG",
		"HOME=/root", /* Required by pandoc */
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run pacta dashboard script: %w", err)
	}

	var artifacts []*task.AnalysisArtifact
	uploadDir := func(dir string) error {
		aas, err := h.uploadDirectory(ctx, dir, dashboardContainer, req.AnalysisID)
		if err != nil {
			return fmt.Errorf("failed to upload report directory: %w", err)
		}
		artifacts = append(artifacts, aas...)
		return nil
	}

	for _, outDir := range dashEnv.outputDirs() {
		if err := uploadDir(outDir); err != nil {
			return fmt.Errorf("failed to upload artifacts %q: %w", outDir, err)
		}
	}

	events := []publisher.Event{
		{
			Data: task.CreateDashboardResponse{
				TaskID:    taskID,
				Request:   req,
				Artifacts: artifacts,
			},
			DataVersion: to.Ptr("1.0"),
			EventType:   to.Ptr("created-dashboard"),
			EventTime:   to.Ptr(time.Now()),
			ID:          to.Ptr(string(taskID)),
			Subject:     to.Ptr(string(taskID)),
		},
	}

	if _, err := h.pubsub.PublishEvents(ctx, events, nil); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	h.logger.Info("created report", zap.String("task_id", string(taskID)))

	return nil
}

func (h *Handler) CreateReport(ctx context.Context, taskID task.ID, req *task.CreateReportRequest, reportContainer string) error {
	if n := len(req.BlobURIs); n != 1 {
		return fmt.Errorf("expected exactly one blob URI as input, got %d", n)
	}
	blobURI := req.BlobURIs[0]

	//  We use this instead of /mnt/... because the base image (quite
	// reasonably) uses a non-root user, so we can't be creating directories in the
	// root filesystem all willy nilly.
	baseDir := filepath.Join("/", "home", "workflow-pacta-webapp")

	reportEnv, err := initEnv(h.benchmarkDir, h.pactaDataDir, baseDir, "create-report")
	if err != nil {
		return fmt.Errorf("failed to init report env: %w", err)
	}

	// Load the parsed portfolio from blob storage, place it in our PORFOLIO_DIR,
	// where the `run_pacta.R` script expects it to be.
	fileNameWithExt := filepath.Base(string(blobURI))
	if !strings.HasSuffix(fileNameWithExt, ".csv") {
		return fmt.Errorf("given blob wasn't a CSV-formatted portfolio, %q", fileNameWithExt)
	}
	destPath := filepath.Join(reportEnv.pathForDir(PortfoliosDir), fileNameWithExt)
	if err := h.downloadBlob(ctx, string(blobURI), destPath); err != nil {
		return fmt.Errorf("failed to download processed portfolio blob: %w", err)
	}

	inp := ReportInput{
		Portfolio: ReportInputPortfolio{
			Files:        []string{fileNameWithExt},
			HoldingsDate: "2023-12-31",   // TODO(#206)
			Name:         "FooPortfolio", // TODO(#206)
		},
		Inherit: "GENERAL_2023Q4", // TODO(#206): Should this be configurable
	}

	var inpJSON bytes.Buffer
	if err := json.NewEncoder(&inpJSON).Encode(inp); err != nil {
		return fmt.Errorf("failed to encode report input as JSON: %w", err)
	}

	cmd := exec.CommandContext(ctx,
		"/usr/local/bin/Rscript",
		"--vanilla", "/workflow.pacta.webapp/inst/extdata/scripts/run_pacta.R",
		inpJSON.String())

	cmd.Env = append(cmd.Env, reportEnv.asEnvVars()...)
	cmd.Env = append(cmd.Env,
		"LOG_LEVEL=DEBUG",
		"HOME=/root", /* Required by pandoc */
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run pacta test CLI: %w", err)
	}

	var artifacts []*task.AnalysisArtifact
	uploadDir := func(dir string) error {
		aas, err := h.uploadDirectory(ctx, dir, reportContainer, req.AnalysisID)
		if err != nil {
			return fmt.Errorf("failed to upload report directory: %w", err)
		}
		artifacts = append(artifacts, aas...)
		return nil
	}

	for _, outDir := range reportEnv.outputDirs() {
		if err := uploadDir(outDir); err != nil {
			return fmt.Errorf("failed to upload artifacts %q: %w", outDir, err)
		}
	}

	events := []publisher.Event{
		{
			Data: task.CreateReportResponse{
				TaskID:    taskID,
				Request:   req,
				Artifacts: artifacts,
			},
			DataVersion: to.Ptr("1.0"),
			EventType:   to.Ptr("created-report"),
			EventTime:   to.Ptr(time.Now()),
			ID:          to.Ptr(string(taskID)),
			Subject:     to.Ptr(string(taskID)),
		},
	}

	if _, err := h.pubsub.PublishEvents(ctx, events, nil); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	h.logger.Info("created report", zap.String("task_id", string(taskID)))

	return nil
}

func (h *Handler) downloadBlob(ctx context.Context, srcURI, destPath string) error {
	// Make sure the destination exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0700); err != nil {
		return fmt.Errorf("failed to create directory to download blob to: %w", err)
	}

	destF, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create dest file: %w", err)
	}
	defer destF.Close() // Best-effort in case something fails

	br, err := h.blob.ReadBlob(ctx, srcURI)
	if err != nil {
		return fmt.Errorf("failed to read raw portfolio: %w", err)
	}
	defer br.Close() // Best-effort in case something fails

	if _, err := io.Copy(destF, br); err != nil {
		return fmt.Errorf("failed to load raw portfolio: %w", err)
	}

	if err := br.Close(); err != nil {
		return fmt.Errorf("failed to close blob reader: %w", err)
	}

	if err := destF.Close(); err != nil {
		return fmt.Errorf("failed to close dest file: %w", err)
	}

	return nil
}

func (h *Handler) uploadDirectory(ctx context.Context, dirPath, container string, analysisID pacta.AnalysisID) ([]*task.AnalysisArtifact, error) {
	base := filepath.Base(dirPath)

	var artifacts []*task.AnalysisArtifact
	err := filepath.WalkDir(dirPath, func(path string, info fs.DirEntry, err error) error {
		if info.IsDir() {
			return nil
		}

		// This is a file, let's upload it to the container
		uri := blob.Join(h.blob.Scheme(), container, string(analysisID), base, strings.TrimPrefix(path, dirPath+"/"))
		if err := h.uploadBlob(ctx, path, uri); err != nil {
			return fmt.Errorf("failed to upload blob: %w", err)
		}

		fn := filepath.Base(path)
		// Returns pacta.FileType_UNKNOWN for unrecognized extensions, which we'll serve as binary blobs.
		ft := fileTypeFromFilename(fn)
		if ft == pacta.FileType_UNKNOWN {
			h.logger.Error("unhandled file extension", zap.String("dir", dirPath), zap.String("filename", fn), zap.String("file_ext", filepath.Ext(fn)))
		}
		artifacts = append(artifacts, &task.AnalysisArtifact{
			BlobURI:  pacta.BlobURI(uri),
			FileName: fn,
			FileType: ft,
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error while walking dir/uploading blobs: %w", err)
	}
	return artifacts, nil
}

func fileTypeFromFilename(fn string) pacta.FileType {
	ext := filepath.Ext(fn)

	switch ext {
	case ".csv":
		return pacta.FileType_CSV
	case ".yaml":
		return pacta.FileType_YAML
	case ".zip":
		return pacta.FileType_ZIP
	case ".html":
		return pacta.FileType_HTML
	case ".json":
		return pacta.FileType_JSON
	case ".txt":
		return pacta.FileType_TEXT
	case ".css":
		return pacta.FileType_CSS
	case ".js":
		return pacta.FileType_JS
	case ".map":
		switch ext2 := filepath.Ext(strings.TrimSuffix(fn, ext)); ext2 {
		case ".js":
			return pacta.FileType_JS_MAP
		case ".css":
			return pacta.FileType_CSS_MAP
		default:
			return pacta.FileType_UNKNOWN
		}
	case ".ttf":
		return pacta.FileType_TTF
	case ".woff":
		return pacta.FileType_WOFF
	case ".woff2":
		return pacta.FileType_WOFF2
	case ".eot":
		return pacta.FileType_EOT
	case ".svg":
		return pacta.FileType_SVG
	case ".png":
		return pacta.FileType_PNG
	case ".jpg":
		return pacta.FileType_JPG
	case ".pdf":
		return pacta.FileType_PDF
	case ".xlsx":
		return pacta.FileType_XLSX
	case ".rds":
		return pacta.FileType_RDS
	default:
		return pacta.FileType_UNKNOWN
	}
}

func (h *Handler) uploadBlob(ctx context.Context, srcPath, destURI string) error {
	h.logger.Info("uploading blob", zap.String("src", srcPath), zap.String("dest", destURI))

	srcF, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open file for upload: %w", err)
	}
	defer srcF.Close() // Best-effort in case something fails

	if err := h.blob.WriteBlob(ctx, destURI, srcF); err != nil {
		return fmt.Errorf("failed to write file to blob storage: %w", err)
	}

	if err := srcF.Close(); err != nil {
		return fmt.Errorf("failed to close source file: %w", err)
	}

	return nil
}
