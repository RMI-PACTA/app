package conv

import (
	"fmt"
	"regexp"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"go.uber.org/zap"
)

var initiativeIDRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func LanguageFromOAPI(l api.Language) (pacta.Language, error) {
	switch l {
	case api.LanguageEN:
		return pacta.Language_EN, nil
	case api.LanguageES:
		return pacta.Language_ES, nil
	case api.LanguageFR:
		return pacta.Language_FR, nil
	case api.LanguageDE:
		return pacta.Language_DE, nil
	}
	return "", oapierr.BadRequest("unknown language", zap.String("language", string(l)))
}

func OptionalBoolFromOAPI(b api.OptionalBoolean) *bool {
	switch b {
	case api.OptionalBooleanFALSE:
		return ptr(false)
	case api.OptionalBooleanTRUE:
		return ptr(true)
	}
	return nil
}

func InitiativeCreateFromOAPI(i *api.InitiativeCreate) (*pacta.Initiative, error) {
	if i == nil {
		return nil, oapierr.BadRequest("InitiativeCreate cannot be nil")
	}
	if !initiativeIDRegex.MatchString(i.Id) {
		return nil, oapierr.BadRequest("id must contain only alphanumeric characters, underscores, and dashes")
	}
	lang, err := LanguageFromOAPI(i.Language)
	if err != nil {
		return nil, err
	}
	var pv *pacta.PACTAVersion
	if i.PactaVersion != nil {
		pv = &pacta.PACTAVersion{ID: pacta.PACTAVersionID(*i.PactaVersion)}
	}
	return &pacta.Initiative{
		Affiliation:              ifNil(i.Affiliation, ""),
		ID:                       pacta.InitiativeID(i.Id),
		InternalDescription:      ifNil(i.InternalDescription, ""),
		IsAcceptingNewMembers:    ifNil(i.IsAcceptingNewMembers, false),
		IsAcceptingNewPortfolios: ifNil(i.IsAcceptingNewPortfolios, false),
		Language:                 lang,
		Name:                     i.Name,
		PACTAVersion:             pv,
		PublicDescription:        ifNil(i.PublicDescription, ""),
		RequiresInvitationToJoin: ifNil(i.RequiresInvitationToJoin, false),
	}, nil
}

func PactaVersionCreateFromOAPI(p *api.PactaVersionCreate) (*pacta.PACTAVersion, error) {
	if p == nil {
		return nil, oapierr.BadRequest("PactaVersionCreate cannot be nil")
	}
	return &pacta.PACTAVersion{
		Name:        p.Name,
		Digest:      p.Digest,
		Description: p.Description,
	}, nil
}

func InitiativeInvitationFromOAPI(i *api.InitiativeInvitationCreate) (*pacta.InitiativeInvitation, error) {
	if i == nil {
		return nil, oapierr.Internal("initiativeInvitationToOAPI: can't convert nil pointer")
	}
	if !initiativeIDRegex.MatchString(i.Id) {
		return nil, oapierr.BadRequest("id must contain only alphanumeric characters, underscores, and dashes")
	}
	if i.InitiativeId == "" {
		return nil, oapierr.BadRequest("initiative_id must not be empty")
	}
	return &pacta.InitiativeInvitation{
		ID:         pacta.InitiativeInvitationID(i.Id),
		Initiative: &pacta.Initiative{ID: pacta.InitiativeID(i.InitiativeId)},
	}, nil
}

func AnalysisTypeFromOAPI(at api.AnalysisType) (pacta.AnalysisType, error) {
	switch at {
	case api.AnalysisTypeAUDIT:
		return pacta.AnalysisType_Audit, nil
	case api.AnalysisTypeREPORT:
		return pacta.AnalysisType_Report, nil
	case api.AnalysisTypeDASHBOARD:
		return pacta.AnalysisType_Dashboard, nil
	}
	return "", oapierr.BadRequest("analysisTypeFromOAPI: unknown analysis type", zap.String("analysis_type", string(at)))
}

func HoldingsDateFromOAPI(hd api.HoldingsDate) (*pacta.HoldingsDate, error) {
	if hd.Time == nil {
		return nil, nil
	}
	return &pacta.HoldingsDate{
		Time: *hd.Time,
	}, nil
}

func PortfolioGroupCreateFromOAPI(pg *api.PortfolioGroupCreate, ownerID pacta.OwnerID) (*pacta.PortfolioGroup, error) {
	if pg == nil {
		return nil, oapierr.Internal("portfolioGroupCreateFromOAPI: can't convert nil pointer")
	}
	if pg.Name == "" {
		return nil, oapierr.BadRequest("name must not be empty")
	}
	if ownerID == "" {
		return nil, oapierr.Internal("portfolioGroupCreateFromOAPI: ownerID must not be empty")
	}
	return &pacta.PortfolioGroup{
		Name:        pg.Name,
		Description: pg.Description,
		Owner:       &pacta.Owner{ID: ownerID},
	}, nil
}

func auditLogActionFromOAPI(i api.AuditLogAction) (pacta.AuditLogAction, error) {
	switch i {
	case api.AuditLogActionCreate:
		return pacta.AuditLogAction_Create, nil
	case api.AuditLogActionUpdate:
		return pacta.AuditLogAction_Update, nil
	case api.AuditLogActionDelete:
		return pacta.AuditLogAction_Delete, nil
	case api.AuditLogActionAddTo:
		return pacta.AuditLogAction_AddTo, nil
	case api.AuditLogActionRemoveFrom:
		return pacta.AuditLogAction_RemoveFrom, nil
	case api.AuditLogActionEnableAdminDebug:
		return pacta.AuditLogAction_EnableAdminDebug, nil
	case api.AuditLogActionDisableAdminDebug:
		return pacta.AuditLogAction_DisableAdminDebug, nil
	case api.AuditLogActionDownload:
		return pacta.AuditLogAction_Download, nil
	case api.AuditLogActionEnableSharing:
		return pacta.AuditLogAction_EnableSharing, nil
	case api.AuditLogActionDisableSharing:
		return pacta.AuditLogAction_DisableSharing, nil
	case api.AuditLogActionReadMetadata:
		return pacta.AuditLogAction_ReadMetadata, nil
	case api.AuditLogActionTransferOwnership:
		return pacta.AuditLogAction_TransferOwnership, nil
	}
	return "", oapierr.BadRequest("unknown audit log action", zap.String("audit_log_action", string(i)))
}

func auditLogActorTypeFromOAPI(i api.AuditLogActorType) (pacta.AuditLogActorType, error) {
	switch i {
	case api.AuditLogActorTypePublic:
		return pacta.AuditLogActorType_Public, nil
	case api.AuditLogActorTypeOwner:
		return pacta.AuditLogActorType_Owner, nil
	case api.AuditLogActorTypeAdmin:
		return pacta.AuditLogActorType_Admin, nil
	case api.AuditLogActorTypeSuperAdmin:
		return pacta.AuditLogActorType_SuperAdmin, nil
	case api.AuditLogActorTypeSystem:
		return pacta.AuditLogActorType_System, nil
	}
	return "", oapierr.BadRequest("unknown audit log actor type", zap.String("audit_log_actor_type", string(i)))
}

func auditLogTargetTypeFromOAPI(i api.AuditLogTargetType) (pacta.AuditLogTargetType, error) {
	switch i {
	case api.AuditLogTargetTypeUser:
		return pacta.AuditLogTargetType_User, nil
	case api.AuditLogTargetTypePortfolio:
		return pacta.AuditLogTargetType_Portfolio, nil
	case api.AuditLogTargetTypeIncompleteUpload:
		return pacta.AuditLogTargetType_IncompleteUpload, nil
	case api.AuditLogTargetTypePortfolioGroup:
		return pacta.AuditLogTargetType_PortfolioGroup, nil
	case api.AuditLogTargetTypeInitiative:
		return pacta.AuditLogTargetType_Initiative, nil
	case api.AuditLogTargetTypeInitiativeInvitation:
		return pacta.AuditLogTargetType_InitiativeInvitation, nil
	case api.AuditLogTargetTypePactaVersion:
		return pacta.AuditLogTargetType_PACTAVersion, nil
	case api.AuditLogTargetTypeAnalysis:
		return pacta.AuditLogTargetType_Analysis, nil
	case api.AuditLogTargetTypeAnalysisArtifact:
		return pacta.AuditLogTargetType_AnalysisArtifact, nil
	}
	return "", oapierr.BadRequest("unknown audit log target type", zap.String("audit_log_target_type", string(i)))
}

func auditLogQueryWhereFromOAPI(i api.AuditLogQueryWhere) (*db.AuditLogQueryWhere, error) {
	result := &db.AuditLogQueryWhere{}
	if i.InId != nil {
		result.InID = fromStrs[pacta.AuditLogID](*i.InId)
	}
	if i.MinCreatedAt != nil {
		result.MinCreatedAt = *i.MinCreatedAt
	}
	if i.MaxCreatedAt != nil {
		result.MaxCreatedAt = *i.MaxCreatedAt
	}
	if i.InAction != nil {
		as, err := convAll(*i.InAction, auditLogActionFromOAPI)
		if err != nil {
			return nil, fmt.Errorf("converting audit log query where in action: %w", err)
		}
		result.InAction = as
	}
	if i.InActorType != nil {
		at, err := convAll(*i.InActorType, auditLogActorTypeFromOAPI)
		if err != nil {
			return nil, fmt.Errorf("converting audit log query where in actor type: %w", err)
		}
		result.InActorType = at
	}
	if i.InActorId != nil {
		result.InActorID = *i.InActorId
	}
	if i.InActorOwnerId != nil {
		result.InActorOwnerID = fromStrs[pacta.OwnerID](*i.InActorOwnerId)
	}
	if i.InTargetType != nil {
		tt, err := convAll(*i.InTargetType, auditLogTargetTypeFromOAPI)
		if err != nil {
			return nil, fmt.Errorf("converting audit log query where in target type: %w", err)
		}
		result.InTargetType = tt
	}
	if i.InTargetId != nil {
		result.InTargetID = *i.InTargetId
	}
	if i.InTargetOwnerId != nil {
		result.InTargetOwnerID = fromStrs[pacta.OwnerID](*i.InTargetOwnerId)
	}
	return result, nil
}

func auditLogQuerySortByFromOAPI(i api.AuditLogQuerySortBy) (db.AuditLogQuerySortBy, error) {
	switch i {
	case api.AuditLogQuerySortByCreatedAt:
		return db.AuditLogQuerySortBy_CreatedAt, nil
	case api.AuditLogQuerySortByActorType:
		return db.AuditLogQuerySortBy_ActorType, nil
	case api.AuditLogQuerySortByActorId:
		return db.AuditLogQuerySortBy_ActorID, nil
	case api.AuditLogQuerySortByActorOwnerId:
		return db.AuditLogQuerySortBy_ActorOwnerID, nil
	case api.AuditLogQuerySortByPrimaryTargetId:
		return db.AuditLogQuerySortBy_PrimaryTargetID, nil
	case api.AuditLogQuerySortByPrimaryTargetType:
		return db.AuditLogQuerySortBy_PrimaryTargetType, nil
	case api.AuditLogQuerySortByPrimaryTargetOwnerId:
		return db.AuditLogQuerySortBy_PrimaryTargetOwnerID, nil
	case api.AuditLogQuerySortBySecondaryTargetId:
		return db.AuditLogQuerySortBy_SecondaryTargetID, nil
	case api.AuditLogQuerySortBySecondaryTargetType:
		return db.AuditLogQuerySortBy_SecondaryTargetType, nil
	case api.AuditLogQuerySortBySecondaryTargetOwnerId:
		return db.AuditLogQuerySortBy_SecondaryTargetOwnerID, nil
	}
	return "", oapierr.BadRequest("unknown audit log query sort by", zap.String("audit_log_query_sort_by", string(i)))
}

func auditLogQuerySortFromOAPI(i api.AuditLogQuerySort) (*db.AuditLogQuerySort, error) {
	by, err := auditLogQuerySortByFromOAPI(i.By)
	if err != nil {
		return nil, fmt.Errorf("converting audit log query sort by: %w", err)
	}
	return &db.AuditLogQuerySort{
		By:        by,
		Ascending: i.Ascending,
	}, nil
}

func AuditLogQueryFromOAPI(q *api.AuditLogQueryReq) (*db.AuditLogQuery, error) {
	limit := 25
	if q.Limit != nil {
		limit = *q.Limit
	}
	if limit > 100 {
		limit = 100
	}
	cursor := ""
	if q.Cursor != nil {
		cursor = *q.Cursor
	}
	sorts := []*db.AuditLogQuerySort{}
	if q.Sorts != nil {
		ss, err := convAll(*q.Sorts, auditLogQuerySortFromOAPI)
		if err != nil {
			return nil, oapierr.BadRequest("error converting audit log query sorts", zap.Error(err))
		}
		sorts = ss
	}
	wheres, err := convAll(q.Wheres, auditLogQueryWhereFromOAPI)
	if err != nil {
		return nil, oapierr.BadRequest("error converting audit log query wheres", zap.Error(err))
	}
	return &db.AuditLogQuery{
		Cursor: db.Cursor(cursor),
		Limit:  limit,
		Wheres: wheres,
		Sorts:  sorts,
	}, nil
}

func userQueryWhereFromOAPI(i api.UserQueryWhere) (*db.UserQueryWhere, error) {
	result := &db.UserQueryWhere{}
	if i.NameOrEmailLike != nil {
		result.NameOrEmailLike = *i.NameOrEmailLike
	}
	return result, nil
}

func UserQueryFromOAPI(q *api.UserQueryReq) (*db.UserQuery, error) {
	limit := 100
	cursor := ""
	if q.Cursor != nil {
		cursor = *q.Cursor
	}
	wheres := []api.UserQueryWhere{}
	if q.Wheres != nil {
		wheres = append(wheres, *q.Wheres...)
	}
	ws, err := convAll(wheres, userQueryWhereFromOAPI)
	if err != nil {
		return nil, oapierr.BadRequest("error converting user query wheres", zap.Error(err))
	}
	return &db.UserQuery{
		Cursor: db.Cursor(cursor),
		Limit:  limit,
		Wheres: ws,
		Sorts: []*db.UserQuerySort{
			{By: db.UserQuerySortBy_CreatedAt, Ascending: false},
		},
	}, nil
}
