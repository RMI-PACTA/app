package pactasrv

import (
	"context"

	"github.com/RMI/pacta/cmd/server/pactasrv/conv"
	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"
)

// Returns a user by ID
// (GET /user/{id})
func (s *Server) FindUserById(ctx context.Context, request api.FindUserByIdRequestObject) (api.FindUserByIdResponseObject, error) {
	// TODO(#12) Implement Authorization
	return nil, oapierr.NotImplemented("findUserById")
}

// Updates user properties
// (PATCH /user/{id})
func (s *Server) UpdateUser(ctx context.Context, request api.UpdateUserRequestObject) (api.UpdateUserResponseObject, error) {
	// TODO(#12) Implement Authorization
	return nil, oapierr.NotImplemented("updateUser")
}

// Deletes a user by ID
// (DELETE /user/{id})
func (s *Server) DeleteUser(ctx context.Context, request api.DeleteUserRequestObject) (api.DeleteUserResponseObject, error) {
	// TODO(#12) Implement Authorization
	return nil, oapierr.NotImplemented("deleteUser")
}

// Returns the logged in user
// (GET /user/me)
func (s *Server) FindUserByMe(ctx context.Context, request api.FindUserByMeRequestObject) (api.FindUserByMeResponseObject, error) {
	token, _, err := jwtauth.FromContext(ctx)
	if err != nil {
		return nil, oapierr.Unauthorized("without authorization token")
	}
	if token == nil {
		return nil, oapierr.Unauthorized("without authorization token")
	}
	emailsClaim, ok := token.PrivateClaims()["emails"]
	if !ok {
		return nil, oapierr.Unauthorized("without email claim in token")
	}
	emails, ok := emailsClaim.([]string)
	if !ok || len(emails) == 0 {
		return nil, oapierr.Internal("couldn't parse email claim in token")
	}
	// TODO(#18) Handle Multiple Emails in the Token Claims gracefully
	if len(emails) > 1 {
		return nil, oapierr.BadRequest("multiple emails in token")
	}
	email := emails[0]
	canonical, err := pacta.CanonicalizeEmail(email)
	if err != nil {
		return nil, oapierr.BadRequest("invalid email on token").WithErrorID(invalidEmail)
	}
	authnID := token.Subject()
	if authnID == "" {
		return nil, oapierr.Internal("couldn't find authn id in jwt")
	}
	user, err := s.DB.GetOrCreateUserByAuthn(s.DB.NoTxn(ctx), pacta.AuthnMechanism_EmailAndPass, authnID, email, canonical)
	if err != nil {
		return nil, oapierr.Internal("failed to get user by authn", zap.String("id", authnID), zap.Error(err))
	}
	u, err := conv.UserToOAPI(user)
	if err != nil {
		return nil, err
	}
	return api.FindUserByMe200JSONResponse(*u), nil
}
