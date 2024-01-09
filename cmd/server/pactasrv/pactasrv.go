package pactasrv

import (
	"context"
	"fmt"
	"time"

	"github.com/RMI/pacta/blob"
	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/oapierr"
	"github.com/RMI/pacta/pacta"
	"github.com/RMI/pacta/session"
	"github.com/RMI/pacta/task"
	"go.uber.org/zap"
)

var (
	// Means we failed to canonicalize someone's email
	invalidEmail = oapierr.ErrorID("invalid_email")
)

type TaskRunner interface {
	ParsePortfolio(ctx context.Context, req *task.ParsePortfolioRequest) (task.ID, task.RunnerID, error)
	CreateReport(ctx context.Context, req *task.CreateReportRequest) (task.ID, task.RunnerID, error)
}

type DB interface {
	Begin(context.Context) (db.Tx, error)
	NoTxn(context.Context) db.Tx
	Transactional(context.Context, func(tx db.Tx) error) error
	RunOrContinueTransaction(db.Tx, func(tx db.Tx) error) error

	Blob(tx db.Tx, id pacta.BlobID) (*pacta.Blob, error)
	Blobs(tx db.Tx, ids []pacta.BlobID) (map[pacta.BlobID]*pacta.Blob, error)
	CreateBlob(tx db.Tx, b *pacta.Blob) (pacta.BlobID, error)
	UpdateBlob(tx db.Tx, id pacta.BlobID, mutations ...db.UpdateBlobFn) error
	DeleteBlob(tx db.Tx, id pacta.BlobID) (pacta.BlobURI, error)
	BlobContexts(tx db.Tx, ids []pacta.BlobID) ([]*pacta.BlobContext, error)

	InitiativeInvitation(tx db.Tx, id pacta.InitiativeInvitationID) (*pacta.InitiativeInvitation, error)
	InitiativeInvitationsByInitiative(tx db.Tx, iid pacta.InitiativeID) ([]*pacta.InitiativeInvitation, error)
	CreateInitiativeInvitation(tx db.Tx, ii *pacta.InitiativeInvitation) (pacta.InitiativeInvitationID, error)
	UpdateInitiativeInvitation(tx db.Tx, id pacta.InitiativeInvitationID, mutations ...db.UpdateInitiativeInvitationFn) error
	DeleteInitiativeInvitation(tx db.Tx, id pacta.InitiativeInvitationID) error

	InitiativeUserRelationship(tx db.Tx, iid pacta.InitiativeID, uid pacta.UserID) (*pacta.InitiativeUserRelationship, error)
	InitiativeUserRelationshipsByUser(tx db.Tx, uid pacta.UserID) ([]*pacta.InitiativeUserRelationship, error)
	InitiativeUserRelationshipsByInitiatives(tx db.Tx, iid pacta.InitiativeID) ([]*pacta.InitiativeUserRelationship, error)
	PutInitiativeUserRelationship(tx db.Tx, iur *pacta.InitiativeUserRelationship) error
	UpdateInitiativeUserRelationship(tx db.Tx, iid pacta.InitiativeID, uid pacta.UserID, mutations ...db.UpdateInitiativeUserRelationshipFn) error

	Initiative(tx db.Tx, id pacta.InitiativeID) (*pacta.Initiative, error)
	Initiatives(tx db.Tx, ids []pacta.InitiativeID) (map[pacta.InitiativeID]*pacta.Initiative, error)
	AllInitiatives(tx db.Tx) ([]*pacta.Initiative, error)
	CreateInitiative(tx db.Tx, i *pacta.Initiative) error
	UpdateInitiative(tx db.Tx, id pacta.InitiativeID, mutations ...db.UpdateInitiativeFn) error
	DeleteInitiative(tx db.Tx, id pacta.InitiativeID) error

	PACTAVersion(tx db.Tx, id pacta.PACTAVersionID) (*pacta.PACTAVersion, error)
	DefaultPACTAVersion(tx db.Tx) (*pacta.PACTAVersion, error)
	PACTAVersions(tx db.Tx) ([]*pacta.PACTAVersion, error)
	CreatePACTAVersion(tx db.Tx, pv *pacta.PACTAVersion) (pacta.PACTAVersionID, error)
	SetDefaultPACTAVersion(tx db.Tx, id pacta.PACTAVersionID) error
	UpdatePACTAVersion(tx db.Tx, id pacta.PACTAVersionID, mutations ...db.UpdatePACTAVersionFn) error
	DeletePACTAVersion(tx db.Tx, id pacta.PACTAVersionID) error

	PortfolioInitiativeMembershipsByPortfolio(tx db.Tx, pid pacta.PortfolioID) ([]*pacta.PortfolioInitiativeMembership, error)
	PortfolioInitiativeMembershipsByInitiative(tx db.Tx, iid pacta.InitiativeID) ([]*pacta.PortfolioInitiativeMembership, error)
	CreatePortfolioInitiativeMembership(tx db.Tx, pim *pacta.PortfolioInitiativeMembership) error
	DeletePortfolioInitiativeMembership(tx db.Tx, pid pacta.PortfolioID, iid pacta.InitiativeID) error

	Portfolio(tx db.Tx, id pacta.PortfolioID) (*pacta.Portfolio, error)
	PortfoliosByOwner(tx db.Tx, owner pacta.OwnerID) ([]*pacta.Portfolio, error)
	Portfolios(tx db.Tx, ids []pacta.PortfolioID) (map[pacta.PortfolioID]*pacta.Portfolio, error)
	CreatePortfolio(tx db.Tx, i *pacta.Portfolio) (pacta.PortfolioID, error)
	UpdatePortfolio(tx db.Tx, id pacta.PortfolioID, mutations ...db.UpdatePortfolioFn) error
	DeletePortfolio(tx db.Tx, id pacta.PortfolioID) ([]pacta.BlobURI, error)

	IncompleteUpload(tx db.Tx, id pacta.IncompleteUploadID) (*pacta.IncompleteUpload, error)
	IncompleteUploads(tx db.Tx, ids []pacta.IncompleteUploadID) (map[pacta.IncompleteUploadID]*pacta.IncompleteUpload, error)
	IncompleteUploadsByOwner(tx db.Tx, owner pacta.OwnerID) ([]*pacta.IncompleteUpload, error)
	CreateIncompleteUpload(tx db.Tx, i *pacta.IncompleteUpload) (pacta.IncompleteUploadID, error)
	UpdateIncompleteUpload(tx db.Tx, id pacta.IncompleteUploadID, mutations ...db.UpdateIncompleteUploadFn) error
	DeleteIncompleteUpload(tx db.Tx, id pacta.IncompleteUploadID) (pacta.BlobURI, error)

	CreateAnalysis(tx db.Tx, a *pacta.Analysis) (pacta.AnalysisID, error)
	UpdateAnalysis(tx db.Tx, id pacta.AnalysisID, mutations ...db.UpdateAnalysisFn) error
	DeleteAnalysis(tx db.Tx, id pacta.AnalysisID) ([]pacta.BlobURI, error)
	Analysis(tx db.Tx, id pacta.AnalysisID) (*pacta.Analysis, error)
	Analyses(tx db.Tx, ids []pacta.AnalysisID) (map[pacta.AnalysisID]*pacta.Analysis, error)
	AnalysesByOwner(tx db.Tx, ownerID pacta.OwnerID) ([]*pacta.Analysis, error)

	AnalysisArtifacts(tx db.Tx, ids []pacta.AnalysisArtifactID) (map[pacta.AnalysisArtifactID]*pacta.AnalysisArtifact, error)
	AnalysisArtifact(tx db.Tx, id pacta.AnalysisArtifactID) (*pacta.AnalysisArtifact, error)
	UpdateAnalysisArtifact(tx db.Tx, id pacta.AnalysisArtifactID, mutations ...db.UpdateAnalysisArtifactFn) error
	DeleteAnalysisArtifact(tx db.Tx, id pacta.AnalysisArtifactID) (pacta.BlobURI, error)

	CreateSnapshotOfPortfolio(tx db.Tx, pID pacta.PortfolioID) (pacta.PortfolioSnapshotID, error)
	CreateSnapshotOfPortfolioGroup(tx db.Tx, pgID pacta.PortfolioGroupID) (pacta.PortfolioSnapshotID, error)
	CreateSnapshotOfInitiative(tx db.Tx, iID pacta.InitiativeID) (pacta.PortfolioSnapshotID, error)

	GetOwnerForUser(tx db.Tx, uID pacta.UserID) (pacta.OwnerID, error)

	PortfolioGroup(tx db.Tx, id pacta.PortfolioGroupID) (*pacta.PortfolioGroup, error)
	PortfolioGroupsByOwner(tx db.Tx, owner pacta.OwnerID) ([]*pacta.PortfolioGroup, error)
	PortfolioGroups(tx db.Tx, ids []pacta.PortfolioGroupID) (map[pacta.PortfolioGroupID]*pacta.PortfolioGroup, error)
	CreatePortfolioGroup(tx db.Tx, p *pacta.PortfolioGroup) (pacta.PortfolioGroupID, error)
	UpdatePortfolioGroup(tx db.Tx, id pacta.PortfolioGroupID, mutations ...db.UpdatePortfolioGroupFn) error
	DeletePortfolioGroup(tx db.Tx, id pacta.PortfolioGroupID) error
	CreatePortfolioGroupMembership(tx db.Tx, pgID pacta.PortfolioGroupID, pID pacta.PortfolioID) error
	DeletePortfolioGroupMembership(tx db.Tx, pgID pacta.PortfolioGroupID, pID pacta.PortfolioID) error

	GetOrCreateUserByAuthn(tx db.Tx, mech pacta.AuthnMechanism, authnID, email, canonicalEmail string) (*pacta.User, error)
	User(tx db.Tx, id pacta.UserID) (*pacta.User, error)
	Users(tx db.Tx, ids []pacta.UserID) (map[pacta.UserID]*pacta.User, error)
	UpdateUser(tx db.Tx, id pacta.UserID, mutations ...db.UpdateUserFn) error
	DeleteUser(tx db.Tx, id pacta.UserID) ([]pacta.BlobURI, error)

	CreateAuditLog(tx db.Tx, a *pacta.AuditLog) (pacta.AuditLogID, error)
	CreateAuditLogs(tx db.Tx, as []*pacta.AuditLog) error
	AuditLogs(tx db.Tx, q *db.AuditLogQuery) ([]*pacta.AuditLog, *db.PageInfo, error)

	RecordUserMerge(tx db.Tx, fromUserID, toUserID, actorUserID pacta.UserID) error
	RecordOwnerMerge(tx db.Tx, fromUserID, toUserID pacta.OwnerID, actorUserID pacta.UserID) error
}

type Blob interface {
	Scheme() blob.Scheme

	SignedUploadURL(ctx context.Context, uri string) (string, time.Time, error)
	SignedDownloadURL(ctx context.Context, uri string) (string, time.Time, error)
	DeleteBlob(ctx context.Context, uri string) error
}

type Server struct {
	DB                DB
	TaskRunner        TaskRunner
	Logger            *zap.Logger
	Blob              Blob
	Now               func() time.Time
	PorfolioUploadURI string
}

func mapAll[I any, O any](is []I, f func(I) (O, error)) ([]O, error) {
	os := make([]O, len(is))
	for i, v := range is {
		o, err := f(v)
		if err != nil {
			return nil, err
		}
		os[i] = o
	}
	return os, nil
}

func dereference[T any](ts []*T, e error) ([]T, error) {
	if e != nil {
		return nil, e
	}
	result := make([]T, len(ts))
	for i, t := range ts {
		if t == nil {
			return nil, oapierr.Internal("nil pointer found in derference", zap.String("type", fmt.Sprintf("%T", t)), zap.Int("index", i))
		}
		result[i] = *t
	}
	return result, nil
}

func values[K comparable, V any](m map[K]*V) []*V {
	result := make([]*V, len(m))
	i := 0
	for _, v := range m {
		result[i] = v
		i++
	}
	return result
}

func getUserID(ctx context.Context) (pacta.UserID, error) {
	userID, err := session.UserIDFromContext(ctx)
	if err != nil {
		return "", oapierr.Unauthorized("error getting authorization token", zap.Error(err))
	}
	return pacta.UserID(userID), nil
}

func (s *Server) getUserOwnerID(ctx context.Context) (pacta.OwnerID, error) {
	userID, err := getUserID(ctx)
	if err != nil {
		return "", err
	}
	ownerID, err := s.DB.GetOwnerForUser(s.DB.NoTxn(ctx), userID)
	if err != nil {
		return "", oapierr.Internal("failed to find or create owner for user",
			zap.String("user_id", string(userID)), zap.Error(err))
	}
	return ownerID, nil
}

func (s *Server) isAdminOrSuperAdmin(ctx context.Context) (bool, bool, error) {
	userID, err := getUserID(ctx)
	if err != nil {
		return false, false, err
	}
	user, err := s.DB.User(s.DB.NoTxn(ctx), userID)
	if err != nil {
		return false, false, oapierr.Internal("failed to find user", zap.Error(err))
	}
	return user.Admin, user.SuperAdmin, nil
}

func asStrs[T ~string](ts []T) []string {
	result := make([]string, len(ts))
	for i, t := range ts {
		result[i] = string(t)
	}
	return result
}

type actorInfo struct {
	UserID       pacta.UserID
	OwnerID      pacta.OwnerID
	IsAdmin      bool
	IsSuperAdmin bool
}

func (s *Server) getactorInfoOrErrIfAnon(ctx context.Context) (*actorInfo, error) {
	actorUserID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}
	actorOwnerID, err := s.getUserOwnerID(ctx)
	if err != nil {
		return nil, err
	}
	actorIsAdmin, actorIsSuperAdmin, err := s.isAdminOrSuperAdmin(ctx)
	if err != nil {
		return nil, err
	}
	return &actorInfo{
		UserID:       actorUserID,
		OwnerID:      actorOwnerID,
		IsAdmin:      actorIsAdmin,
		IsSuperAdmin: actorIsSuperAdmin,
	}, nil
}
