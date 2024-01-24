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
	"github.com/RMI/pacta/session"
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
	NoTxn(context.Context) db.Tx

	Analysis(tx db.Tx, id pacta.AnalysisID) (*pacta.Analysis, error)
	AnalysisArtifactsForAnalysis(tx db.Tx, id pacta.AnalysisID) ([]*pacta.AnalysisArtifact, error)
	Blobs(tx db.Tx, ids []pacta.BlobID) (map[pacta.BlobID]*pacta.Blob, error)
	GetOwnerForUser(tx db.Tx, userID pacta.UserID) (pacta.OwnerID, error)
	CreateAuditLog(tx db.Tx, log *pacta.AuditLog) (pacta.AuditLogID, error)
	User(tx db.Tx, id pacta.UserID) (*pacta.User, error)
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
	r.Get("/report/{analysis_id}", func(w http.ResponseWriter, r *http.Request) {
		analysisID := chi.URLParam(r, "analysis_id")
		newPath := "/report/" + analysisID + "/"
		http.Redirect(w, r, newPath, http.StatusTemporaryRedirect)
	})
	r.Get("/report/{analysis_id}/*", s.serveReport)
}

func (s *Server) serveReport(w http.ResponseWriter, r *http.Request) {
	aID := pacta.AnalysisID(chi.URLParam(r, "analysis_id"))
	if aID == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	artifacts, err := s.db.AnalysisArtifactsForAnalysis(s.db.NoTxn(ctx), aID)
	if err != nil {
		s.logger.Error("failed to load artifacts for analysis", zap.String("analysis_id", string(aID)), zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var blobIDs []pacta.BlobID
	for _, a := range artifacts {
		blobIDs = append(blobIDs, a.Blob.ID)
	}

	blobs, err := s.db.Blobs(s.db.NoTxn(ctx), blobIDs)
	if err != nil {
		s.logger.Error("failed to load blobs", zap.Strings("blob_ids", asStrs(blobIDs)), zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	a, err := s.db.Analysis(s.db.NoTxn(ctx), aID)
	if err != nil {
		if strings.HasPrefix(string(aID), "analysis") {
			s.logger.Error("failed to load analysis", zap.String("analysis_id", string(aID)), zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		} else {
			s.logger.Info("poorly constructed analysis id", zap.String("path", string(aID)), zap.Error(err))
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
	}

	subPath := strings.TrimPrefix(r.URL.Path, "/report/"+string(aID))
	if strings.HasPrefix(subPath, "/") {
		subPath = subPath[1:]
	}
	if subPath == "" {
		subPath = "index.html"
	}

	for _, aa := range artifacts {
		// Container is just 'reports', we can ignore that.
		b, ok := blobs[aa.Blob.ID]
		if !ok {
			s.logger.Error("no blob loaded for blob ID", zap.String("analysis_artifact_id", string(aa.ID)), zap.String("blob_id", string(aa.Blob.ID)))
			continue
		}
		uri := string(b.BlobURI)
		_, path, ok := blob.SplitURI(s.blob.Scheme(), uri)
		if !ok {
			s.logger.Error("blob had invalid URI", zap.String("analysis_artifact_id", string(aa.ID)), zap.String("blob_uri", uri))
			continue
		}
		_, uriPath, ok := strings.Cut(path, "/")
		if !ok {
			s.logger.Error("path had no UUID prefix", zap.String("analysis_artifact_id", string(aa.ID)), zap.String("blob_uri", uri), zap.String("blob_path", path))
			continue
		}

		// This isn't the asset we're looking for.
		if uriPath != subPath {
			continue
		}

		if ok := s.doAuthzAndAuditLog(a, aa, w, r); !ok {
			// Note that doAuthzAndAuditLog will have already written the response.
			return
		}

		// Load the blob from blob storage and copy it byte-for-byte to the HTTP response.
		r, err := s.blob.ReadBlob(ctx, uri)
		if err != nil {
			http.Error(w, "failed to read blob: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Close()

		w.Header().Set("Content-Type", fileTypeToMIME(b.FileType))
		if _, err := io.Copy(w, r); err != nil {
			s.logger.Error("failed to read/write blob", zap.String("blob_uri", uri), zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}

	s.logger.Info("unknown report asset", zap.String("analysis_id", string(aID)), zap.String("req_path", r.URL.Path), zap.String("asset_path", subPath))
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	return
}

func (s *Server) doAuthzAndAuditLog(a *pacta.Analysis, aa *pacta.AnalysisArtifact, w http.ResponseWriter, r *http.Request) bool {
	ctx := r.Context()
	const unauthenticatedUserID = "unauthenticated user"
	var actorOwner *pacta.Owner
	actorID, _ := session.UserIDFromContext(ctx)
	if actorID == "" {
		actorID = unauthenticatedUserID
		actorOwner = &pacta.Owner{ID: unauthenticatedUserID}
	} else {
		ownerID, err := s.db.GetOwnerForUser(s.db.NoTxn(ctx), actorID)
		if err != nil {
			s.logger.Error("failed to get owner for user", zap.String("user_id", string(actorID)), zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return false
		}
		actorOwner = &pacta.Owner{ID: ownerID}
	}
	auditLog := &pacta.AuditLog{
		Action:               pacta.AuditLogAction_Download,
		ActorID:              string(actorID),
		ActorOwner:           actorOwner,
		PrimaryTargetType:    pacta.AuditLogTargetType_AnalysisArtifact,
		PrimaryTargetID:      string(aa.ID),
		PrimaryTargetOwner:   a.Owner,
		SecondaryTargetType:  pacta.AuditLogTargetType_Analysis,
		SecondaryTargetID:    string(a.ID),
		SecondaryTargetOwner: a.Owner,
	}
	allowIfAuditLogSaves := func() bool {
		if _, err := s.db.CreateAuditLog(s.db.NoTxn(ctx), auditLog); err != nil {
			s.logger.Error("failed to create audit log", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return false
		}
		return true
	}
	if aa.SharedToPublic {
		auditLog.ActorType = pacta.AuditLogActorType_Public
		return allowIfAuditLogSaves()
	}
	if actorID == unauthenticatedUserID {
		s.logger.Info("unauthenticated user attempted to read asset", zap.String("user_id", string(actorID)))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return false
	}
	if a.Owner.ID == actorOwner.ID {
		auditLog.ActorType = pacta.AuditLogActorType_Owner
		return allowIfAuditLogSaves()
	}
	if aa.AdminDebugEnabled {
		user, err := s.db.User(s.db.NoTxn(ctx), actorID)
		if err != nil {
			s.logger.Error("failed to get user", zap.String("user_id", string(actorID)), zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return false
		}
		if user.Admin || user.SuperAdmin {
			auditLog.ActorType = pacta.AuditLogActorType_Admin
			return allowIfAuditLogSaves()
		}
	}
	s.logger.Info("unauthorized user attempted to read asset", zap.String("user_id", string(actorID)), zap.String("analysis_artifact_id", string(aa.ID)), zap.String("analysis_id", string(a.ID)))
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	return false
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

func asStrs[T ~string](ts []T) []string {
	ss := make([]string, len(ts))
	for i, t := range ts {
		ss[i] = string(t)
	}
	return ss
}
