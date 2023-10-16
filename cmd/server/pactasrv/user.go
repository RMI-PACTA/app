package pactasrv

import (
	"context"

	"github.com/RMI/pacta/cmd/server/pactasrv/conv"
	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
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
		lang, err := pacta.ParseLanguage(string(*request.Body.PreferredLanguage))
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
	err := s.DB.DeleteUser(s.DB.NoTxn(ctx), pacta.UserID(request.Id))
	if err != nil {
		return nil, oapierr.Internal("failed to delete user", zap.Error(err))
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
