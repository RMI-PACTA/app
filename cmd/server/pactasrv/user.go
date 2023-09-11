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

// Returns the logged in user
// (GET /user/me)
func (s *Server) FindUserByMe(ctx context.Context, request api.FindUserByMeRequestObject) (api.FindUserByMeResponseObject, error) {
	unauth := api.FindUserByMe401JSONResponse(map[string]interface{}{})
	token, _, err := jwtauth.FromContext(ctx)
	if err != nil {
		// TODO(grady) upgrade this to the new error handling strategy after #12
		return nil, fmt.Errorf("failed to get token from context: %w", err)
	}
	if token == nil {
		return unauth, nil
	}
	email, ok := token.PrivateClaims()["email"]
	if !ok {
		return nil, fmt.Errorf("failed to load email claim")
	}
	emailStr, ok := email.(string)
	if !ok || emailStr == "" {
		return nil, fmt.Errorf("email claim was of unexpected type %T, wanted a non-empty string", email)
	}
	authnID := token.Subject()
	if authnID == "" {
		return nil, fmt.Errorf("couldn't find authn id in jwt")
	}
	user, err := s.DB.GetOrCreateUserByAuthn(s.DB.NoTxn(ctx), pacta.AuthnMechanism_EmailAndPass, authnID, emailStr)
	if err != nil {
		return nil, fmt.Errorf("getting user by authn: %w", err)
	}
	return api.FindUserByMe200JSONResponse{
		Admin:             user.Admin,
		CanonicalEmail:    &user.CanonicalEmail,
		EnteredEmail:      user.EnteredEmail,
		Id:                string(user.ID),
		Name:              user.Name,
		PreferredLanguage: api.UserPreferredLanguage(user.PreferredLanguage),
		SuperAdmin:        user.SuperAdmin,
	}, nil
}
