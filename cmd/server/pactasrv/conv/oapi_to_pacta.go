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

func InitiativeCreateFromOAPI(i *api.InitiativeCreate) (*pacta.Initiative, error) {
	if i == nil {
		return nil, oapierr.BadRequest("InitiativeCreate cannot be nil")
	}
	if !initiativeIDRegex.MatchString(i.Id) {
		return nil, oapierr.BadRequest("id must contain only alphanumeric characters, underscores, and dashes")
	}
	lang, err := pacta.ParseLanguage(string(i.Language))
	if err != nil {
		return nil, oapierr.BadRequest("failed to parse language", zap.Error(err))
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

func AnalysisTypeFromOAPI(at *api.AnalysisType) (*pacta.AnalysisType, error) {
	if at == nil {
		return nil, oapierr.BadRequest("analysisTypeFromOAPI: can't convert nil pointer")
	}
	switch string(*at) {
	case "audit":
		return ptr(pacta.AnalysisType_Audit), nil
	case "report":
		return ptr(pacta.AnalysisType_Report), nil
	}
	return nil, oapierr.BadRequest("analysisTypeFromOAPI: unknown analysis type", zap.String("analysis_type", string(*at)))
}

func HoldingsDateFromOAPI(hd *api.HoldingsDate) (*pacta.HoldingsDate, error) {
	if hd == nil {
		return nil, nil
	}
	return &pacta.HoldingsDate{
		Time: hd.Time,
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
	return pacta.ParseAuditLogAction(string(i))
}

func auditLogActorTypeFromOAPI(i api.AuditLogActorType) (pacta.AuditLogActorType, error) {
	return pacta.ParseAuditLogActorType(string(i))
}

func auditLogTargetTypeFromOAPI(i api.AuditLogTargetType) (pacta.AuditLogTargetType, error) {
	return pacta.ParseAuditLogTargetType(string(i))
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
	return db.ParseAuditLogQuerySortBy(string(i))
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
