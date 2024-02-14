// Package async provides the business logic for our async tasks.
package async

import (
	"bufio"
	"bytes"
	"context"
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
	"github.com/RMI/pacta/blob"
	"github.com/RMI/pacta/pacta"
	"github.com/RMI/pacta/task"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Config struct {
	Blob   Blob
	PubSub *publisher.Client
	Logger *zap.Logger
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
}

func New(cfg *Config) (*Handler, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config given: %w", err)
	}

	return &Handler{
		blob:   cfg.Blob,
		pubsub: cfg.PubSub,
		logger: cfg.Logger,
	}, nil
}

// TODO: Send a notification when parsing fails.
func (h *Handler) ParsePortfolio(ctx context.Context, taskID task.ID, req *task.ParsePortfolioRequest, destPortfolioContainer string) error {
	// Load the portfolio from blob storage, place it in /mnt/raw_portfolios, where
	// the `process_portfolios.R` script expects it to be.
	for _, srcURI := range req.BlobURIs {
		id := uuid.New().String()
		// TODO: Probably set the CSV extension in the signed upload URL instead.
		destPath := filepath.Join("/", "mnt", "raw_portfolios", fmt.Sprintf("%s.csv", id))
		if err := h.downloadBlob(ctx, string(srcURI), destPath); err != nil {
			return fmt.Errorf("failed to download raw portfolio blob: %w", err)
		}
	}

	processedDir := filepath.Join("/", "mnt", "processed_portfolios")
	if err := os.MkdirAll(processedDir, 0600); err != nil {
		return fmt.Errorf("failed to create directory to download blob to: %w", err)
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "/usr/local/bin/Rscript", "/app/process_portfolios.R")
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run process_portfolios script: %w", err)
	}

	sc := bufio.NewScanner(&stderr)

	// TODO: Load from the output database file (or similar, like reading the processed_portfolios dir) instead of parsing stderr
	var paths []string
	for sc.Scan() {
		line := sc.Text()
		idx := strings.Index(line, "writing to file: " /* 17 chars */)
		if idx == -1 {
			continue
		}
		paths = append(paths, strings.TrimSpace(line[idx+17:]))
	}

	// NOTE: This code could benefit from some concurrency, but I'm opting not to prematurely optimize.
	var out []*task.ParsePortfolioResponseItem
	for _, p := range paths {
		lineCount, err := countCSVLines(p)
		if err != nil {
			return fmt.Errorf("failed to count lines in file %q: %w", p, err)
		}
		fileName := filepath.Base(p)
		blobURI := pacta.BlobURI(blob.Join(h.blob.Scheme(), destPortfolioContainer, fileName))
		if err := h.uploadBlob(ctx, p, string(blobURI)); err != nil {
			return fmt.Errorf("failed to copy parsed portfolio from %q to %q: %w", p, blobURI, err)
		}
		extension := filepath.Ext(fileName)
		fileType, err := pacta.ParseFileType(extension)
		if err != nil {
			return fmt.Errorf("failed to parse file type from file name %q: %w", fileName, err)
		}
		out = append(out, &task.ParsePortfolioResponseItem{
			Blob: pacta.Blob{
				FileName: fileName,
				FileType: fileType,
				BlobURI:  blobURI,
			},
			LineCount: lineCount,
		})
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

// TODO(grady): Move this line counting into the image to prevent having our code do any read of the actual underlying data.
func countCSVLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("opening file failed: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scanner.error returned: %w", err)
	}
	// Subtract 1 for the header row
	return lineCount - 1, nil
}

func (h *Handler) CreateAudit(ctx context.Context, taskID task.ID, req *task.CreateAuditRequest) error {
	return errors.New("not implemented")
}

func (h *Handler) CreateReport(ctx context.Context, taskID task.ID, req *task.CreateReportRequest, reportContainer string) error {
	fileNames := []string{}
	for _, blobURI := range req.BlobURIs {
		// Load the parsed portfolio from blob storage, place it in /mnt/
		// processed_portfolios, where the `create_report.R` script expects it
		// to be.
		fileNameWithExt := filepath.Base(string(blobURI))
		if !strings.HasSuffix(fileNameWithExt, ".json") {
			return fmt.Errorf("given blob wasn't a JSON-formatted portfolio, %q", fileNameWithExt)
		}
		fileNames = append(fileNames, strings.TrimSuffix(fileNameWithExt, ".json"))
		destPath := filepath.Join("/", "mnt", "processed_portfolios", fileNameWithExt)
		if err := h.downloadBlob(ctx, string(blobURI), destPath); err != nil {
			return fmt.Errorf("failed to download processed portfolio blob: %w", err)
		}
	}

	reportDir := filepath.Join("/", "mnt", "reports")
	if err := os.MkdirAll(reportDir, 0600); err != nil {
		return fmt.Errorf("failed to create directory for reports to get copied to: %w", err)
	}

	cmd := exec.CommandContext(ctx, "/usr/local/bin/Rscript", "/app/create_report.R")
	cmd.Env = append(cmd.Env,
		"PORTFOLIO="+strings.Join(fileNames, ","),
		"HOME=/root", /* Required by pandoc */
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run pacta test CLI: %w", err)
	}

	// Download outputs from from /out and upload them to Azure
	dirEntries, err := os.ReadDir(reportDir)
	if err != nil {
		return fmt.Errorf("failed to read report directory: %w", err)
	}

	var artifacts []*task.AnalysisArtifact
	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			continue
		}
		dirPath := filepath.Join(reportDir, dirEntry.Name())
		tmp, err := h.uploadDirectory(ctx, dirPath, reportContainer)
		if err != nil {
			return fmt.Errorf("failed to upload report directory: %w", err)
		}
		artifacts = tmp
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
	if err := os.MkdirAll(filepath.Dir(destPath), 0600); err != nil {
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

func (h *Handler) uploadDirectory(ctx context.Context, dirPath, container string) ([]*task.AnalysisArtifact, error) {
	base := filepath.Base(dirPath)

	var artifacts []*task.AnalysisArtifact
	err := filepath.WalkDir(dirPath, func(path string, info fs.DirEntry, err error) error {
		if info.IsDir() {
			return nil
		}

		// This is a file, let's upload it to the container
		uri := blob.Join(h.blob.Scheme(), container, base, strings.TrimPrefix(path, dirPath+"/"))
		if err := h.uploadBlob(ctx, path, uri); err != nil {
			return fmt.Errorf("failed to upload blob: %w", err)
		}

		fn := filepath.Base(path)
		// Returns pacta.FileType_UNKNOWN for unrecognized extensions, which we'll serve as binary blobs.
		ft := fileTypeFromExt(filepath.Ext(fn))
		if ft == pacta.FileType_UNKNOWN {
			h.logger.Error("unhandled file extension", zap.String("dir", dirPath), zap.String("file_ext", filepath.Ext(fn)))
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

func fileTypeFromExt(ext string) pacta.FileType {
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
	case ".ttf":
		return pacta.FileType_TTF
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
