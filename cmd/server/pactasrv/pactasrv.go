// Package pactasrv currently just implements the Petstore API interface, the
// bog standard OpenAPI example. It implements pactaapi.StrictServerInterface,
// which is auto-generated from the OpenAPI 3 spec.
package pactasrv

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RMI/pacta"
	"github.com/RMI/pacta/db"

	pactaapi "github.com/RMI/pacta/openapi/pacta"
)

type DB interface {
	Begin(context.Context) (db.Tx, error)
	NoTxn(context.Context) db.Tx
	Transactional(context.Context, func(tx db.Tx) error) error
	RunOrContinueTransaction(db.Tx, func(tx db.Tx) error) error

	User(db.Tx, pacta.UserID) (*pacta.User, error)
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
