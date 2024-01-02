package pactasrv

import (
	"context"

	"github.com/RMI/pacta/cmd/server/pactasrv/conv"
	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"go.uber.org/zap"
)

// queries the platform's audit logs
// (POST /audit-logs)
func (s *Server) ListAuditLogs(ctx context.Context, request api.ListAuditLogsRequestObject) (api.ListAuditLogsResponseObject, error) {
	// TODO(#12) implement authorization
	query, err := conv.AuditLogQueryFromOAPI(request.Body)
	if err != nil {
		return nil, err
	}
	// TODO(#12) implement additional authorizations, ensuring for example that:
	// - every generated query has reasonable limits + only filters by allowed search terms
	// - the actor is allowed to see the audit logs of the actor_owner, but not of other actor_owners
	// - initiative admins should be able to see audit logs of the initiative, but not initiative members
	// - admins should be able to see all
	// This is probably our most important piece of authz-ery, so it should be thoroughly tested.
	als, pi, err := s.DB.AuditLogs(s.DB.NoTxn(ctx), query)
	if err != nil {
		return nil, oapierr.Internal("querying audit logs failed", zap.Error(err))
	}
	results, err := dereference(conv.AuditLogsToOAPI(als))
	if err != nil {
		return nil, err
	}
	return api.ListAuditLogs200JSONResponse{
		AuditLogs:   results,
		Cursor:      string(pi.Cursor),
		HasNextPage: pi.HasNextPage,
	}, nil
}
