package pactasrv

import (
	"context"

	api "github.com/RMI/pacta/openapi/pacta"
)

// Returns a user by ID
// (GET /user/{id})
func (s *Server) FindUserById(ctx context.Context, request api.FindUserByIdRequestObject) (api.FindUserByIdResponseObject, error) {
	u, err := s.findUserById(ctx, request)
	if err != nil {
		e := errToAPIError(err)
		return api.FindUserByIddefaultJSONResponse{
			Body:       e,
			StatusCode: int(e.Code),
		}, nil
	}
	return api.FindUserById200JSONResponse(*u), nil
}

func (s *Server) findUserById(ctx context.Context, request api.FindUserByIdRequestObject) (*api.User, error) {
	// TODO(#12) Implement Authorization
	return nil, errorNotImplemented("findUserById")
}

// Updates user properties
// (PATCH /user/{id})
func (s *Server) UpdateUser(ctx context.Context, request api.UpdateUserRequestObject) (api.UpdateUserResponseObject, error) {
	err := s.updateUser(ctx, request)
	if err != nil {
		e := errToAPIError(err)
		return api.UpdateUserdefaultJSONResponse{
			Body:       e,
			StatusCode: int(e.Code),
		}, nil
	}
	return api.UpdateUser200JSONResponse(emptySuccess), nil
}

func (s *Server) updateUser(ctx context.Context, request api.UpdateUserRequestObject) error {
	// TODO(#12) Implement Authorization
	return errorNotImplemented("updateUser")
}

// Deletes a user by ID
// (DELETE /user/{id})
func (s *Server) DeleteUser(ctx context.Context, request api.DeleteUserRequestObject) (api.DeleteUserResponseObject, error) {
	err := s.deleteUser(ctx, request)
	if err != nil {
		e := errToAPIError(err)
		return api.DeleteUserdefaultJSONResponse{
			Body:       e,
			StatusCode: int(e.Code),
		}, nil
	}
	return api.DeleteUser200JSONResponse(emptySuccess), nil
}

func (s *Server) deleteUser(ctx context.Context, request api.DeleteUserRequestObject) error {
	// TODO(#12) Implement Authorization
	return errorNotImplemented("deleteUser")
}
