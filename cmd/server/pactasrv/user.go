package pactasrv

import (
	"context"
	"fmt"

	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"github.com/go-chi/jwtauth/v5"
)

// Returns a user by ID
// (GET /user/{id})
func (s *Server) FindUserById(ctx context.Context, request api.FindUserByIdRequestObject) (api.FindUserByIdResponseObject, error) {
	u, err := s.findUserById(ctx, request)
	if err != nil {
		return errToAPIError(err)
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
		return errToAPIError(err)
	}
	return api.UpdateUser200JSONResponse{}, nil
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
		return errToAPIError(err)
	}
	return api.DeleteUser200JSONResponse{}, nil
}

func (s *Server) deleteUser(ctx context.Context, request api.DeleteUserRequestObject) error {
	// TODO(#12) Implement Authorization
	return errorNotImplemented("deleteUser")
}

// Returns the logged in user
// (GET /user/me)
func (s *Server) FindUserByMe(ctx context.Context, request api.FindUserByMeRequestObject) (api.FindUserByMeResponseObject, error) {
	user, err := s.findUserByMe(ctx, request)
	if err != nil {
		return errToAPIError(err)
	}
	return api.FindUserByMe200JSONResponse(*user), nil
}

func (s *Server) findUserByMe(ctx context.Context, request api.FindUserByMeRequestObject) (*api.User, error) {
	token, _, err := jwtauth.FromContext(ctx)
	if err != nil {
		// TODO(grady) upgrade this to the new error handling strategy after #12
		return nil, errorUnauthorized("lookup self", "without authorization token")
	}
	if token == nil {
		return nil, errorUnauthorized("lookup self ", "without authorization token")
	}
	emailsClaim, ok := token.PrivateClaims()["emails"]
	if !ok {
		return nil, errorUnauthorized("lookup self", "without email claim in token")
	}
	emails, ok := emailsClaim.([]string)
	if !ok || len(emails) == 0 {
		return nil, errorInternal(fmt.Errorf("couldn't parse email claim in token"))
	}
	// TODO(#18) Handle Multiple Emails in the Token Claims gracefully
	if len(emails) > 1 {
		return nil, errorBadRequest("jwt token", "multiple emails in token")
	}
	email := emails[0]
	canonical, err := pacta.CanonicalizeEmail(email)
	if err != nil {
		return nil, errorBadRequest("email", "invalid email on token")
	}
	authnID := token.Subject()
	if authnID == "" {
		return nil, fmt.Errorf("couldn't find authn id in jwt")
	}
	user, err := s.DB.GetOrCreateUserByAuthn(s.DB.NoTxn(ctx), pacta.AuthnMechanism_EmailAndPass, authnID, email, canonical)
	if err != nil {
		return nil, fmt.Errorf("getting user by authn: %w", err)
	}
	return userToOAPI(user)
}
