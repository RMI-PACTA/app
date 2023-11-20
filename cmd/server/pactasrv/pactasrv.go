package pactasrv

import (
	"context"
	"fmt"

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

	IncompleteUpload(tx db.Tx, id pacta.IncompleteUploadID) (*pacta.IncompleteUpload, error)
	IncompleteUploads(tx db.Tx, ids []pacta.IncompleteUploadID) (map[pacta.IncompleteUploadID]*pacta.IncompleteUpload, error)
	CreateIncompleteUpload(tx db.Tx, i *pacta.IncompleteUpload) (pacta.IncompleteUploadID, error)
	UpdateIncompleteUpload(tx db.Tx, id pacta.IncompleteUploadID, mutations ...db.UpdateIncompleteUploadFn) error
	DeleteIncompleteUpload(tx db.Tx, id pacta.IncompleteUploadID) (pacta.BlobURI, error)

	GetOrCreateUserByAuthn(tx db.Tx, mech pacta.AuthnMechanism, authnID, email, canonicalEmail string) (*pacta.User, error)
	User(tx db.Tx, id pacta.UserID) (*pacta.User, error)
	Users(tx db.Tx, ids []pacta.UserID) (map[pacta.UserID]*pacta.User, error)
	UpdateUser(tx db.Tx, id pacta.UserID, mutations ...db.UpdateUserFn) error
	DeleteUser(tx db.Tx, id pacta.UserID) error
}

type Blob interface {
	Scheme() blob.Scheme

	// For uploading portfolios
	SignedUploadURL(ctx context.Context, uri string) (string, error)
	// For downloading reports
	SignedDownloadURL(ctx context.Context, uri string) (string, error)
	DeleteBlobs(ctx context.Context, uris []string) error
}

type Server struct {
	DB                DB
	TaskRunner        TaskRunner
	Logger            *zap.Logger
	Blob              Blob
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

func getUserID(ctx context.Context) (pacta.UserID, error) {
	userID, err := session.UserIDFromContext(ctx)
	if err != nil {
		return "", oapierr.Unauthorized("error getting authorization token", zap.Error(err))
	}
	return pacta.UserID(userID), nil
}
