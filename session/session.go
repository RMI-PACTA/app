package session

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/oapierr"
	"github.com/RMI/pacta/pacta"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"
)

var userIDKey = struct{}{}
var allowedAnonymousKey = struct{}{}

const allowedAnonymousValue = "allowedAnonymous"

func WithUserID(c context.Context, id pacta.UserID) context.Context {
	return context.WithValue(c, userIDKey, id)
}

func WithAllowedAnonymous(c context.Context) context.Context {
	return context.WithValue(c, allowedAnonymousKey, allowedAnonymousValue)
}

type DB interface {
	NoTxn(context.Context) db.Tx
	UserByAuthn(tx db.Tx, mech pacta.AuthnMechanism, authnID string) (*pacta.User, error)
}

func UserIDFromContext(ctx context.Context) (pacta.UserID, error) {
	userID, ok := ctx.Value(userIDKey).(pacta.UserID)
	if !ok {
		return "", oapierr.Unauthorized("no user id in context")
	}
	if userID == "" {
		return "", oapierr.Unauthorized("empty user id in context")
	}
	return userID, nil
}

func IsAllowedAnonymous(ctx context.Context) bool {
	v := ctx.Value(allowedAnonymousKey)
	if v == nil {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	return s == allowedAnonymousValue
}

func WithAuthn(logger *zap.Logger, d DB) func(http.Handler) http.Handler {
	fn := func(c context.Context) (context.Context, error) {
		token, _, err := jwtauth.FromContext(c)
		if err != nil {
			return nil, fmt.Errorf("error getting authorization token: %w", err)
		}
		if token == nil {
			return nil, fmt.Errorf("nil authorization token")
		}
		authnID := token.Subject()
		if authnID == "" {
			return nil, fmt.Errorf("couldn't find authn id in jwt")
		}
		user, err := d.UserByAuthn(d.NoTxn(c), pacta.AuthnMechanism_EmailAndPass, authnID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user by authn: %w", err)
		}
		return WithUserID(c, user.ID), nil
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, err := fn(r.Context())
			if err != nil {
				// Optionally log errors here when debugging authentication access.
				// logger.Warn("couldn't authenticate", zap.Error(err))
				// http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				next.ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
