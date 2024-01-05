package pactasrv

import (
	"context"
	"fmt"

	"github.com/RMI/pacta/cmd/server/pactasrv/conv"
	"github.com/RMI/pacta/db"
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
	user, err := s.DB.User(s.DB.NoTxn(ctx), pacta.UserID(request.Id))
	if err != nil {
		return nil, oapierr.Internal("failed to retrieve user", zap.Error(err))
	}
	result, err := conv.UserToOAPI(user)
	if err != nil {
		return nil, err
	}
	return api.FindUserById200JSONResponse(*result), nil
}

// Updates user properties
// (PATCH /user/{id})
func (s *Server) UpdateUser(ctx context.Context, request api.UpdateUserRequestObject) (api.UpdateUserResponseObject, error) {
	// TODO(#12) Implement Authorization
	mutations := []db.UpdateUserFn{}
	if request.Body.Admin != nil {
		mutations = append(mutations, db.SetUserAdmin(*request.Body.Admin))
	}
	if request.Body.SuperAdmin != nil {
		mutations = append(mutations, db.SetUserSuperAdmin(*request.Body.SuperAdmin))
	}
	if request.Body.Name != nil {
		mutations = append(mutations, db.SetUserName(*request.Body.Name))
	}
	if request.Body.PreferredLanguage != nil {
		lang, err := conv.LanguageFromOAPI(*request.Body.PreferredLanguage)
		if err != nil {
			return nil, oapierr.BadRequest("invalid language", zap.Error(err))
		}
		mutations = append(mutations, db.SetUserPreferredLanguage(lang))
	}
	err := s.DB.UpdateUser(s.DB.NoTxn(ctx), pacta.UserID(request.Id), mutations...)
	if err != nil {
		return nil, oapierr.Internal("failed to update user", zap.Error(err))
	}
	return api.UpdateUser204Response{}, nil
}

// Deletes a user by ID
// (DELETE /user/{id})
func (s *Server) DeleteUser(ctx context.Context, request api.DeleteUserRequestObject) (api.DeleteUserResponseObject, error) {
	// TODO(#12) Implement Authorization
	blobURIs, err := s.DB.DeleteUser(s.DB.NoTxn(ctx), pacta.UserID(request.Id))
	if err != nil {
		return nil, oapierr.Internal("failed to delete user", zap.Error(err))
	}
	if err := s.deleteBlobs(ctx, blobURIs...); err != nil {
		return nil, err
	}
	return api.DeleteUser204Response{}, nil
}

// Returns the logged in user
// (GET /user/me)
func (s *Server) FindUserByMe(ctx context.Context, request api.FindUserByMeRequestObject) (api.FindUserByMeResponseObject, error) {
	// Looking for how users get created or populated into the context?
	// It's in the HTTP handler for adding a user to the context in main.go
	meID, err := getUserID(ctx)
	if err != nil {
		return nil, oapierr.Unauthorized("no user id found in context", zap.Error(err))
	}
	user, err := s.DB.User(s.DB.NoTxn(ctx), meID)
	if err != nil {
		return nil, oapierr.Internal("failed to retrieve user", zap.Error(err))
	}
	result, err := conv.UserToOAPI(user)
	if err != nil {
		return nil, err
	}
	return api.FindUserByMe200JSONResponse(*result), nil
}

// a callback after login to create or return the user
// (POST /user/authentication-followup)
func (s *Server) UserAuthenticationFollowup(ctx context.Context, _request api.UserAuthenticationFollowupRequestObject) (api.UserAuthenticationFollowupResponseObject, error) {
	token, _, err := jwtauth.FromContext(ctx)
	if err != nil {
		return nil, oapierr.BadRequest("error getting authorization token", zap.Error(err))
	}
	if token == nil {
		return nil, oapierr.BadRequest("nil authorization token")
	}
	emailsClaim, ok := token.PrivateClaims()["emails"]
	if !ok {
		return nil, oapierr.BadRequest("no email claim in token")
	}
	emails, ok := emailsClaim.([]interface{})
	if !ok || len(emails) == 0 {
		return nil, oapierr.BadRequest("couldn't find email claim in token", zap.String("emails_claim_type", fmt.Sprintf("%T", emailsClaim)))
	}
	// TODO(#18) Handle Multiple Emails in the Token Claims gracefully
	if len(emails) > 1 {
		return nil, oapierr.BadRequest(fmt.Sprintf("multiple emails in token: %+v", emails))
	}
	email, ok := emails[0].(string)
	if !ok {
		return nil, oapierr.BadRequest("wrong type for email claim", zap.String("email_claim_type", fmt.Sprintf("%T", emails[0])))
	}
	canonical, err := pacta.CanonicalizeEmail(email)
	if err != nil {
		return nil, oapierr.BadRequest(fmt.Sprintf("invalid email: %q", email), zap.String("email", email), zap.Error(err))
	}
	authnID := token.Subject()
	if authnID == "" {
		return nil, oapierr.BadRequest("couldn't find authn id in jwt")
	}
	user, err := s.DB.GetOrCreateUserByAuthn(s.DB.NoTxn(ctx), pacta.AuthnMechanism_EmailAndPass, authnID, email, canonical)
	if err != nil {
		return nil, fmt.Errorf("failed to GetOrCreateUser by authn: %w", err)
	}
	result, err := conv.UserToOAPI(user)
	if err != nil {
		return nil, err
	}
	return api.UserAuthenticationFollowup200JSONResponse(*result), nil

}
