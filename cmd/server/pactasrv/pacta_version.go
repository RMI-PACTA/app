package pactasrv

import (
	"context"
	"fmt"

	"github.com/RMI/pacta/cmd/server/pactasrv/conv"
	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"go.uber.org/zap"
)

// Returns a version of the PACTA model by ID
// (GET /pacta-version/{id})
func (s *Server) FindPactaVersionById(ctx context.Context, request api.FindPactaVersionByIdRequestObject) (api.FindPactaVersionByIdResponseObject, error) {
	id := pacta.PACTAVersionID(request.Id)
	if err := s.pactaVersionAuthz(ctx, id, pacta.AuditLogAction_ReadMetadata); err != nil {
		return nil, err
	}
	pv, err := s.DB.PACTAVersion(s.DB.NoTxn(ctx), pacta.PACTAVersionID(request.Id))
	if err != nil {
		if db.IsNotFound(err) {
			return nil, oapierr.NotFound("pacta version not found", zap.String("pacta_version_id", request.Id))
		}
		return nil, oapierr.Internal("failed to load pacta version", zap.String("pacta_version_id", request.Id), zap.Error(err))
	}
	resp, err := conv.PactaVersionToOAPI(pv)
	if err != nil {
		return nil, err
	}
	return api.FindPactaVersionById200JSONResponse(*resp), nil
}

// Returns all versions of the PACTA model
// (GET /pacta-versions)
func (s *Server) ListPactaVersions(ctx context.Context, request api.ListPactaVersionsRequestObject) (api.ListPactaVersionsResponseObject, error) {
	// No Authorization - Public
	pvs, err := s.DB.PACTAVersions(s.DB.NoTxn(ctx))
	if err != nil {
		return nil, oapierr.Internal("failed to list pacta versions", zap.Error(err))
	}
	resp, err := dereference(mapAll(pvs, conv.PactaVersionToOAPI))
	if err != nil {
		return nil, err
	}
	return api.ListPactaVersions200JSONResponse(resp), nil
}

// Creates a PACTA version
// (POST /pacta-versions)
func (s *Server) CreatePactaVersion(ctx context.Context, request api.CreatePactaVersionRequestObject) (api.CreatePactaVersionResponseObject, error) {
	if err := anyError(
		checkStringLimitSmall("name", request.Body.Name),
		checkStringLimitSmall("digest", request.Body.Digest),
		checkStringLimitMedium("description", request.Body.Description),
	); err != nil {
		return nil, err
	}
	actorInfo, err := s.getActorInfoOrErrIfAnon(ctx)
	if err != nil {
		return nil, err
	}
	var auditLogActorType pacta.AuditLogActorType
	if actorInfo.IsAdmin {
		auditLogActorType = pacta.AuditLogActorType_Admin
	} else if actorInfo.IsSuperAdmin {
		auditLogActorType = pacta.AuditLogActorType_SuperAdmin
	} else {
		return nil, oapierr.Forbidden("only admins and super admins can create initiatives")
	}
	pv, err := conv.PactaVersionCreateFromOAPI(request.Body)
	if err != nil {
		return nil, err
	}
	pvID, err := s.DB.CreatePACTAVersion(s.DB.NoTxn(ctx), pv)
	if err != nil {
		return nil, oapierr.Internal("failed to create pacta version", zap.Error(err))
	}
	if err := s.auditLogForCreateEvent(ctx, actorInfo, auditLogActorType, pacta.AuditLogTargetType_PACTAVersion, string(pvID)); err != nil {
		return nil, err
	}
	return api.CreatePactaVersion204Response{}, nil
}

// Updates a PACTA version
// (PATCH /pacta-version/{id})
func (s *Server) UpdatePactaVersion(ctx context.Context, request api.UpdatePactaVersionRequestObject) (api.UpdatePactaVersionResponseObject, error) {
	if err := anyError(
		checkStringLimitSmallPtr("name", request.Body.Name),
		checkStringLimitSmallPtr("digest", request.Body.Digest),
		checkStringLimitMediumPtr("description", request.Body.Description),
	); err != nil {
		return nil, err
	}
	id := pacta.PACTAVersionID(request.Id)
	if err := s.pactaVersionAuthz(ctx, id, pacta.AuditLogAction_Update); err != nil {
		return nil, err
	}
	mutations := []db.UpdatePACTAVersionFn{}
	b := request.Body
	if b.Description != nil {
		mutations = append(mutations, db.SetPACTAVersionDescription(*b.Description))
	}
	if b.Digest != nil {
		mutations = append(mutations, db.SetPACTAVersionDigest(*b.Digest))
	}
	if b.Name != nil {
		mutations = append(mutations, db.SetPACTAVersionName(*b.Name))
	}
	err := s.DB.UpdatePACTAVersion(s.DB.NoTxn(ctx), id, mutations...)
	if err != nil {
		return nil, oapierr.Internal("failed to update pacta version", zap.String("pacta_version_id", string(id)), zap.Error(err))
	}
	return api.UpdatePactaVersion204Response{}, nil

}

// Deletes a pacta version by ID
// (DELETE /pacta-version/{id})
func (s *Server) DeletePactaVersion(ctx context.Context, request api.DeletePactaVersionRequestObject) (api.DeletePactaVersionResponseObject, error) {
	id := pacta.PACTAVersionID(request.Id)
	if err := s.pactaVersionAuthz(ctx, id, pacta.AuditLogAction_Delete); err != nil {
		return nil, err
	}
	err := s.DB.DeletePACTAVersion(s.DB.NoTxn(ctx), id)
	if err != nil {
		return nil, oapierr.Internal("failed to delete pacta version", zap.String("pacta_version_id", string(id)), zap.Error(err))
	}
	return api.DeletePactaVersion204Response{}, nil
}

// Marks this version of the PACTA model as the default
// (POST /pacta-version/{id}/set-default)
func (s *Server) MarkPactaVersionAsDefault(ctx context.Context, request api.MarkPactaVersionAsDefaultRequestObject) (api.MarkPactaVersionAsDefaultResponseObject, error) {
	id := pacta.PACTAVersionID(request.Id)
	if err := s.pactaVersionAuthz(ctx, id, pacta.AuditLogAction_Update); err != nil {
		return nil, err
	}
	err := s.DB.SetDefaultPACTAVersion(s.DB.NoTxn(ctx), id)
	if err != nil {
		return nil, oapierr.Internal("failed to set default pacta version", zap.String("pacta_version_id", string(id)), zap.Error(err))
	}
	return api.MarkPactaVersionAsDefault204Response{}, nil
}

func (s *Server) pactaVersionAuthz(ctx context.Context, pvID pacta.PACTAVersionID, action pacta.AuditLogAction) error {
	actorInfo, err := s.getActorInfoOrAnon(ctx)
	if err != nil {
		return err
	}
	as := &authzStatus{
		primaryTargetID:      string(pvID),
		primaryTargetType:    pacta.AuditLogTargetType_PACTAVersion,
		primaryTargetOwnerID: systemOwnedEntityOwner,
		actorInfo:            actorInfo,
		action:               action,
	}
	switch action {
	case pacta.AuditLogAction_ReadMetadata:
		as.isAuthorized = true
		as.authorizedAsActorType = ptr(pacta.AuditLogActorType_Public)
	case pacta.AuditLogAction_Delete, pacta.AuditLogAction_Update:
		as.isAuthorized, as.authorizedAsActorType = allowIfAdmin(actorInfo)
	default:
		return fmt.Errorf("unknown action %q for pacta_version authz", action)
	}
	return s.auditLogIfAuthorizedOrFail(ctx, as)
}
