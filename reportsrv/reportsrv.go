package reportsrv

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/RMI/pacta/blob"
	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Config struct {
	DB     DB
	Blob   Blob
	Logger *zap.Logger
}

func (c *Config) validate() error {
	if c.DB == nil {
		return errors.New("no DB was given")
	}
	if c.Blob == nil {
		return errors.New("no blob client was given")
	}
	if c.Logger == nil {
		return errors.New("no logger was given")
	}
	return nil
}

type Server struct {
	db     DB
	blob   Blob
	logger *zap.Logger
}

type DB interface {
	Begin(context.Context) (db.Tx, error)
	NoTxn(context.Context) db.Tx
	Transactional(context.Context, func(tx db.Tx) error) error
	RunOrContinueTransaction(db.Tx, func(tx db.Tx) error) error

	AnalysisArtifactsForAnalysis(tx db.Tx, id pacta.AnalysisID) ([]*pacta.AnalysisArtifact, error)
	Blobs(tx db.Tx, ids []pacta.BlobID) (map[pacta.BlobID]*pacta.Blob, error)
}

type Blob interface {
	Scheme() blob.Scheme

	ReadBlob(ctx context.Context, uri string) (io.ReadCloser, error)
}

func New(cfg *Config) (*Server, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &Server{
		db:     cfg.DB,
		blob:   cfg.Blob,
		logger: cfg.Logger,
	}, nil
}

func (s *Server) verifyRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(#12): Figure out how to share authz between the main pactasrv and this.
		// Should be a fast follow, but large enough to be a separate change.
		next.ServeHTTP(w, r)
	})
}

func (s *Server) RegisterHandlers(r chi.Router) {
	r.Use(s.verifyRequest)
	r.Get("/report/{analysis_id}/*", s.serveReport)
}

func (s *Server) serveReport(w http.ResponseWriter, r *http.Request) {
	aID := pacta.AnalysisID(chi.URLParam(r, "analysis_id"))
	if aID == "" {
		http.Error(w, "no id given", http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	artifacts, err := s.db.AnalysisArtifactsForAnalysis(s.db.NoTxn(ctx), aID)
	if err != nil {
		http.Error(w, "failed to get artifacts: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var blobIDs []pacta.BlobID
	for _, a := range artifacts {
		blobIDs = append(blobIDs, a.Blob.ID)
	}

	blobs, err := s.db.Blobs(s.db.NoTxn(ctx), blobIDs)
	if err != nil {
		http.Error(w, "failed to get blobs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	subPath := strings.TrimPrefix(r.URL.Path, "/report/"+string(aID))
	if strings.HasPrefix(subPath, "/") {
		subPath = subPath[1:]
	}
	if subPath == "" {
		subPath = "index.html"
	}

	for _, a := range artifacts {
		// Container is just 'reports', we can ignore that.
		b, ok := blobs[a.Blob.ID]
		if !ok {
			s.logger.Error("no blob loaded for blob ID", zap.String("analysis_artifact_id", string(a.ID)), zap.String("blob_id", string(a.Blob.ID)))
			continue
		}
		uri := string(b.BlobURI)
		_, path, ok := blob.SplitURI(s.blob.Scheme(), uri)
		if !ok {
			s.logger.Error("blob had invalid URI", zap.String("analysis_artifact_id", string(a.ID)), zap.String("blob_uri", uri))
			continue
		}
		_, uriPath, ok := strings.Cut(path, "/")
		if !ok {
			s.logger.Error("path had no UUID prefix", zap.String("analysis_artifact_id", string(a.ID)), zap.String("blob_uri", uri), zap.String("blob_path", path))
			continue
		}

		if uriPath != subPath {
			continue
		}

		r, err := s.blob.ReadBlob(ctx, uri)
		if err != nil {
			http.Error(w, "failed to read blob: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Close()
		w.Header().Set("Content-Type", fileTypeToMIME(b.FileType))
		if _, err := io.Copy(w, r); err != nil {
			http.Error(w, "failed to read/write blob: "+err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	http.Error(w, "not found", http.StatusNotFound)
	return
}

func fileTypeToMIME(ft pacta.FileType) string {
	switch ft {
	case pacta.FileType_CSV:
		return "text/csv"
	case pacta.FileType_YAML:
		// Note: This one is actually kinda contentious, but I don't think it matters
		// much. See https://stackoverflow.com/q/332129
		return "text/yaml"
	case pacta.FileType_ZIP:
		return "application/zip"
	case pacta.FileType_HTML:
		return "text/html"
	case pacta.FileType_JSON:
		return "application/json"
	case pacta.FileType_TEXT:
		return "text/plain"
	case pacta.FileType_CSS:
		return "text/css"
	case pacta.FileType_JS:
		return "text/javascript"
	case pacta.FileType_TTF:
		return "font/ttf"
	default:
		return "application/octet-stream"
	}

}
