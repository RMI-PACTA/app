// Package pactasrv currently just implements the Petstore API interface, the
// bog standard OpenAPI example. It implements pactaapi.StrictServerInterface,
// which is auto-generated from the OpenAPI 3 spec.
package pactasrv

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RMI/pacta/db"

	pactaapi "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
)

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

	User(tx db.Tx, id pacta.UserID) (*pacta.User, error)
	UserByAuthn(tx db.Tx, authnMechanism pacta.AuthnMechanism, authnID string) (*pacta.User, error)
	Users(tx db.Tx, ids []pacta.UserID) (map[pacta.UserID]*pacta.User, error)
	CreateUser(tx db.Tx, u *pacta.User) (pacta.UserID, error)
	UpdateUser(tx db.Tx, id pacta.UserID, mutations ...db.UpdateUserFn) error
	DeleteUser(tx db.Tx, id pacta.UserID) error
}

type Server struct {
	DB   DB
	pets []pactaapi.Pet
	idx  int
}

// Returns all pets
// (GET /pets)
func (s *Server) FindPets(ctx context.Context, req pactaapi.FindPetsRequestObject) (pactaapi.FindPetsResponseObject, error) {
	out := []pactaapi.Pet{}
	for _, p := range s.pets {
		if matchesTag(req.Params.Tags, p.Tag) {
			out = append(out, p)
		}
	}
	if req.Params.Limit != nil && *req.Params.Limit > 0 && int(*req.Params.Limit) < len(out) {
		out = out[:*req.Params.Limit]
	}

	return pactaapi.FindPets200JSONResponse(out), nil
}

func matchesTag(tags *[]string, tag *string) bool {
	if tags == nil || tag == nil || len(*tags) == 0 {
		return true
	}
	for _, t := range *tags {
		if t == *tag {
			return true
		}
	}
	return false
}

// Creates a new pet
// (POST /pets)
func (s *Server) AddPet(ctx context.Context, req pactaapi.AddPetRequestObject) (pactaapi.AddPetResponseObject, error) {
	s.idx += 1
	id := int64(s.idx)
	s.pets = append(s.pets, pactaapi.Pet{
		Id:   id,
		Name: req.Body.Name,
		Tag:  req.Body.Tag,
	})
	return pactaapi.AddPet200JSONResponse{
		Id:   id,
		Name: req.Body.Name,
		Tag:  req.Body.Tag,
	}, nil
}

// Deletes a pet by ID
// (DELETE /pets/{id})
func (s *Server) DeletePet(ctx context.Context, req pactaapi.DeletePetRequestObject) (pactaapi.DeletePetResponseObject, error) {
	for i, p := range s.pets {
		if p.Id == req.Id {
			s.pets = append(s.pets[:i], s.pets[i+1:]...)
			return pactaapi.DeletePet204Response{}, nil
		}
	}
	return pactaapi.DeletePetdefaultJSONResponse{
		Body: pactaapi.Error{
			Code:    1,
			Message: fmt.Sprintf("no pet with id %d", req.Id),
		},
		StatusCode: http.StatusNotFound,
	}, nil
}

// Returns a pet by ID
// (GET /pets/{id})
func (s *Server) FindPetByID(ctx context.Context, req pactaapi.FindPetByIDRequestObject) (pactaapi.FindPetByIDResponseObject, error) {
	for _, p := range s.pets {
		if p.Id == req.Id {
			return pactaapi.FindPetByID200JSONResponse{
				Id:   p.Id,
				Name: p.Name,
				Tag:  p.Tag,
			}, nil
		}
	}

	return pactaapi.FindPetByIDdefaultJSONResponse{
		Body: pactaapi.Error{
			Code:    2,
			Message: fmt.Sprintf("no pet with id %d", req.Id),
		},
		StatusCode: http.StatusNotFound,
	}, nil
}
