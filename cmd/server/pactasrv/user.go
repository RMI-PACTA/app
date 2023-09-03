package pactasrv

import (
	"context"
	"fmt"

	api "github.com/RMI/pacta/openapi/pacta"
)

// Returns a user by ID
// (GET /user/{id})
func (s *Server) FindUserById(ctx context.Context, request api.FindUserByIdRequestObject) (api.FindUserByIdResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// Updates user properties
// (PATCH /user/{id})
func (s *Server) UpdateUser(ctx context.Context, request api.UpdateUserRequestObject) (api.UpdateUserResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// Deletes a user by ID
// (DELETE /user/{id})
func (s *Server) DeleteUser(ctx context.Context, request api.DeleteUserRequestObject) (api.DeleteUserResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
