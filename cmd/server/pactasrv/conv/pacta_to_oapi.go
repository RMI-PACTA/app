package conv

import (
	"fmt"

	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"go.uber.org/zap"
)

func InitiativeToOAPI(i *pacta.Initiative) (*api.Initiative, error) {
	if i == nil {
		return nil, oapierr.Internal("initiativeToOAPI: can't convert nil pointer")
	}
	pims, err := convAll(i.PortfolioInitiativeMemberships, portfolioInitiativeMembershipToOAPIPortfolio)
	if err != nil {
		return nil, oapierr.Internal("initiativeToOAPI: portfolioInitiativeMembershipToOAPIInitiative failed", zap.Error(err))
	}
	return &api.Initiative{
		Affiliation:                    i.Affiliation,
		CreatedAt:                      i.CreatedAt,
		Id:                             string(i.ID),
		InternalDescription:            i.InternalDescription,
		IsAcceptingNewMembers:          i.IsAcceptingNewMembers,
		IsAcceptingNewPortfolios:       i.IsAcceptingNewPortfolios,
		Language:                       api.InitiativeLanguage(i.Language),
		Name:                           i.Name,
		PactaVersion:                   strPtr(i.PACTAVersion.ID),
		PublicDescription:              i.PublicDescription,
		RequiresInvitationToJoin:       i.RequiresInvitationToJoin,
		PortfolioInitiativeMemberships: pims,
	}, nil
}

func portfolioInitiativeMembershipToOAPIPortfolio(in *pacta.PortfolioInitiativeMembership) (api.PortfolioInitiativeMembershipPortfolio, error) {
	var zero api.PortfolioInitiativeMembershipPortfolio
	out := &api.PortfolioInitiativeMembershipPortfolio{
		CreatedAt: in.CreatedAt,
	}
	if in.AddedBy != nil && in.AddedBy.ID == "" {
		out.AddedByUserId = strPtr(in.AddedBy.ID)
	}
	p, err := PortfolioToOAPI(in.Portfolio)
	if err != nil {
		return zero, oapierr.Internal("portfolioInitiativeMembershipToOAPI: portfolioToOAPI failed", zap.Error(err))
	}
	out.Portfolio = *p
	return zero, nil
}

func portfolioInitiativeMembershipToOAPIInitiative(in *pacta.PortfolioInitiativeMembership) (api.PortfolioInitiativeMembershipInitiative, error) {
	var zero api.PortfolioInitiativeMembershipInitiative
	out := api.PortfolioInitiativeMembershipInitiative{
		CreatedAt: in.CreatedAt,
	}
	if in.AddedBy != nil && in.AddedBy.ID == "" {
		out.AddedByUserId = strPtr(in.AddedBy.ID)
	}
	if in.Initiative != nil && in.Initiative.PACTAVersion != nil {
		i, err := InitiativeToOAPI(in.Initiative)
		if err != nil {
			return zero, oapierr.Internal("portfolioInitiativeMembershipToOAPI: initiativeToOAPI failed", zap.Error(err))
		}
		out.Initiative = *i
	}
	return out, nil
}

func UserToOAPI(user *pacta.User) (*api.User, error) {
	if user == nil {
		return nil, oapierr.Internal("userToOAPI: can't convert nil pointer")
	}
	return &api.User{
		Admin:             user.Admin,
		CanonicalEmail:    &user.CanonicalEmail,
		EnteredEmail:      user.EnteredEmail,
		Id:                string(user.ID),
		Name:              user.Name,
		PreferredLanguage: api.UserPreferredLanguage(user.PreferredLanguage),
		SuperAdmin:        user.SuperAdmin,
	}, nil
}

func PactaVersionToOAPI(pv *pacta.PACTAVersion) (*api.PactaVersion, error) {
	if pv == nil {
		return nil, oapierr.Internal("pactaVersionToOAPI: can't convert nil pointer")
	}
	return &api.PactaVersion{
		CreatedAt:   pv.CreatedAt,
		Description: pv.Description,
		Digest:      pv.Digest,
		Id:          string(pv.ID),
		IsDefault:   pv.IsDefault,
		Name:        pv.Name,
	}, nil
}

func InitiativeInvitationToOAPI(i *pacta.InitiativeInvitation) (*api.InitiativeInvitation, error) {
	if i == nil {
		return nil, oapierr.Internal("initiativeToOAPI: can't convert nil pointer")
	}
	var usedAt *string
	if !i.UsedAt.IsZero() {
		usedAt = ptr(i.UsedAt.String())
	}
	var usedBy *string
	if i.UsedBy != nil {
		usedBy = strPtr(i.UsedBy.ID)
	}
	return &api.InitiativeInvitation{
		CreatedAt:    i.CreatedAt,
		Id:           string(i.ID),
		InitiativeId: string(i.Initiative.ID),
		UsedAt:       usedAt,
		UsedByUserId: usedBy,
	}, nil
}

func InitiativeUserRelationshipToOAPI(i *pacta.InitiativeUserRelationship) (*api.InitiativeUserRelationship, error) {
	if i == nil {
		return nil, oapierr.Internal("initiativeUserRelationshipToOAPI: can't convert nil pointer")
	}
	if i.User == nil {
		return nil, oapierr.Internal("initiativeUserRelationshipToOAPI: can't convert nil user")
	}
	if i.Initiative == nil {
		return nil, oapierr.Internal("initiativeUserRelationshipToOAPI: can't convert nil initiative")
	}
	return &api.InitiativeUserRelationship{
		UpdatedAt:    i.UpdatedAt,
		InitiativeId: string(i.Initiative.ID),
		UserId:       string(i.User.ID),
		Manager:      i.Manager,
		Member:       i.Member,
	}, nil
}

func HoldingsDateToOAPI(hd *pacta.HoldingsDate) (*api.HoldingsDate, error) {
	if hd == nil {
		return nil, nil
	}
	return &api.HoldingsDate{
		Time: hd.Time,
	}, nil
}

func IncompleteUploadsToOAPI(ius []*pacta.IncompleteUpload) ([]*api.IncompleteUpload, error) {
	return convAll(ius, IncompleteUploadToOAPI)
}

func IncompleteUploadToOAPI(iu *pacta.IncompleteUpload) (*api.IncompleteUpload, error) {
	if iu == nil {
		return nil, oapierr.Internal("incompleteUploadToOAPI: can't convert nil pointer")
	}
	hd, err := HoldingsDateToOAPI(iu.HoldingsDate)
	if err != nil {
		return nil, oapierr.Internal("incompleteUploadToOAPI: holdingsDateToOAPI failed", zap.Error(err))
	}
	return &api.IncompleteUpload{
		Id:                string(iu.ID),
		Name:              iu.Name,
		Description:       iu.Description,
		HoldingsDate:      hd,
		CreatedAt:         iu.CreatedAt,
		RanAt:             timeToNilable(iu.RanAt),
		CompletedAt:       timeToNilable(iu.CompletedAt),
		FailureCode:       stringToNilable(iu.FailureCode),
		FailureMessage:    stringToNilable(iu.FailureMessage),
		AdminDebugEnabled: iu.AdminDebugEnabled,
	}, nil
}

func PortfoliosToOAPI(ius []*pacta.Portfolio) ([]*api.Portfolio, error) {
	return convAll(ius, PortfolioToOAPI)
}

func PortfolioToOAPI(p *pacta.Portfolio) (*api.Portfolio, error) {
	if p == nil {
		return nil, oapierr.Internal("portfolioToOAPI: can't convert nil pointer")
	}
	hd, err := HoldingsDateToOAPI(p.HoldingsDate)
	if err != nil {
		return nil, oapierr.Internal("portfolioToOAPI: holdingsDateToOAPI failed", zap.Error(err))
	}
	portfolioGroupMemberships := []api.PortfolioGroupMembershipPortfolioGroup{}
	for _, m := range p.PortfolioGroupMemberships {
		pg, err := PortfolioGroupToOAPI(m.PortfolioGroup)
		if err != nil {
			return nil, oapierr.Internal("portfolioToOAPI: portfolioGroupToOAPI failed", zap.Error(err))
		}
		portfolioGroupMemberships = append(portfolioGroupMemberships, api.PortfolioGroupMembershipPortfolioGroup{
			CreatedAt:      m.CreatedAt,
			PortfolioGroup: *pg,
		})
	}
	pims, err := convAll(p.PortfolioInitiativeMemberships, portfolioInitiativeMembershipToOAPIInitiative)
	if err != nil {
		return nil, oapierr.Internal("initiativeToOAPI: portfolioInitiativeMembershipToOAPIInitiative failed", zap.Error(err))
	}
	return &api.Portfolio{
		Id:                string(p.ID),
		Name:              p.Name,
		Description:       p.Description,
		HoldingsDate:      hd,
		CreatedAt:         p.CreatedAt,
		NumberOfRows:      p.NumberOfRows,
		AdminDebugEnabled: p.AdminDebugEnabled,
		Groups:            &portfolioGroupMemberships,
		Initiatives:       &pims,
	}, nil
}

func PortfolioGroupToOAPI(pg *pacta.PortfolioGroup) (*api.PortfolioGroup, error) {
	if pg == nil {
		return nil, oapierr.Internal("portfolioGroupToOAPI: can't convert nil pointer")
	}
	members := []api.PortfolioGroupMembershipPortfolio{}
	for _, m := range pg.PortfolioGroupMemberships {
		portfolio, err := PortfolioToOAPI(m.Portfolio)
		if err != nil {
			return nil, oapierr.Internal("portfolioGroupToOAPI: portfolioToOAPI failed", zap.Error(err))
		}
		members = append(members, api.PortfolioGroupMembershipPortfolio{
			CreatedAt: m.CreatedAt,
			Portfolio: *portfolio,
		})
	}
	return &api.PortfolioGroup{
		Id:          string(pg.ID),
		Name:        pg.Name,
		Description: pg.Description,
		CreatedAt:   pg.CreatedAt,
		Members:     &members,
	}, nil
}

func PortfolioGroupsToOAPI(pgs []*pacta.PortfolioGroup) ([]*api.PortfolioGroup, error) {
	return convAll(pgs, PortfolioGroupToOAPI)
}

func auditLogActorTypeToOAPI(i pacta.AuditLogActorType) (api.AuditLogActorType, error) {
	switch i {
	case pacta.AuditLogActorType_Public:
		return api.AuditLogActorTypePUBLIC, nil
	case pacta.AuditLogActorType_Owner:
		return api.AuditLogActorTypeOWNER, nil
	case pacta.AuditLogActorType_Admin:
		return api.AuditLogActorTypeADMIN, nil
	case pacta.AuditLogActorType_SuperAdmin:
		return api.AuditLogActorTypeSUPERADMIN, nil
	case pacta.AuditLogActorType_System:
		return api.AuditLogActorTypeSYSTEM, nil
	}
	return "", oapierr.Internal(fmt.Sprintf("auditLogActorTypeToOAPI: unknown actor type: %q", i))
}

func auditLogActionToOAPI(i pacta.AuditLogAction) (api.AuditLogAction, error) {
	switch i {
	case pacta.AuditLogAction_Create:
		return api.AuditLogActionCREATE, nil
	case pacta.AuditLogAction_Update:
		return api.AuditLogActionUPDATE, nil
	case pacta.AuditLogAction_Delete:
		return api.AuditLogActionDELETE, nil
	case pacta.AuditLogAction_AddTo:
		return api.AuditLogActionADDTO, nil
	case pacta.AuditLogAction_RemoveFrom:
		return api.AuditLogActionREMOVEFROM, nil
	case pacta.AuditLogAction_EnableAdminDebug:
		return api.AuditLogActionENABLEADMINDEBUG, nil
	case pacta.AuditLogAction_DisableAdminDebug:
		return api.AuditLogActionDISABLEADMINDEBUG, nil
	case pacta.AuditLogAction_Download:
		return api.AuditLogActionDOWNLOAD, nil
	case pacta.AuditLogAction_EnableSharing:
		return api.AuditLogActionENABLESHARING, nil
	case pacta.AuditLogAction_DisableSharing:
		return api.AuditLogActionDISABLESHARING, nil
	}
	return "", oapierr.Internal(fmt.Sprintf("auditLogActionToOAPI: unknown action: %q", i))
}

func auditLogTargetTypeToOAPI(i pacta.AuditLogTargetType) (api.AuditLogTargetType, error) {
	switch i {
	case pacta.AuditLogTargetType_User:
		return api.AuditLogTargetTypeUSER, nil
	case pacta.AuditLogTargetType_Portfolio:
		return api.AuditLogTargetTypePORTFOLIO, nil
	case pacta.AuditLogTargetType_IncompleteUpload:
		return api.AuditLogTargetTypeINCOMPLETEUPLOAD, nil
	case pacta.AuditLogTargetType_PortfolioGroup:
		return api.AuditLogTargetTypePORTFOLIOGROUP, nil
	case pacta.AuditLogTargetType_Initiative:
		return api.AuditLogTargetTypeINITIATIVE, nil
	case pacta.AuditLogTargetType_PACTAVersion:
		return api.AuditLogTargetTypePACTAVERSION, nil
	case pacta.AuditLogTargetType_Analysis:
		return api.AuditLogTargetTypeANALYSIS, nil
	}
	return "", oapierr.Internal(fmt.Sprintf("auditLogTargetTypeToOAPI: unknown target type: %q", i))
}

func AuditLogToOAPI(al *pacta.AuditLog) (*api.AuditLog, error) {
	if al == nil {
		return nil, oapierr.Internal("auditLogToOAPI: can't convert nil pointer")
	}
	at, err := auditLogActorTypeToOAPI(al.ActorType)
	if err != nil {
		return nil, oapierr.Internal("auditLogToOAPI: auditLogActorTypeToOAPI failed", zap.Error(err))
	}
	act, err := auditLogActionToOAPI(al.Action)
	if err != nil {
		return nil, oapierr.Internal("auditLogToOAPI: auditLogActionToOAPI failed", zap.Error(err))
	}
	ptt, err := auditLogTargetTypeToOAPI(al.PrimaryTargetType)
	if err != nil {
		return nil, oapierr.Internal("auditLogToOAPI: auditLogTargetTypeToOAPI failed", zap.Error(err))
	}
	var aoi *string
	if al.ActorOwner != nil {
		aoi = stringToNilable(al.ActorOwner.ID)
	}
	var stt *api.AuditLogTargetType
	if al.SecondaryTargetType != "" {
		s, err := auditLogTargetTypeToOAPI(al.SecondaryTargetType)
		if err != nil {
			return nil, oapierr.Internal("auditLogToOAPI: auditLogTargetTypeToOAPI failed", zap.Error(err))
		}
		stt = &s
	}
	var sto *string
	if al.SecondaryTargetOwner != nil {
		s := string(al.SecondaryTargetOwner.ID)
		sto = &s
	}
	var sid *string
	if al.SecondaryTargetID != "" {
		s := string(al.SecondaryTargetID)
		sid = &s
	}
	return &api.AuditLog{
		Id:                   string(al.ID),
		CreatedAt:            al.CreatedAt,
		ActorType:            at,
		ActorId:              stringToNilable(al.ActorID),
		ActorOwnerId:         aoi,
		Action:               act,
		PrimaryTargetType:    ptt,
		PrimaryTargetId:      al.PrimaryTargetID,
		PrimaryTargetOwner:   string(al.PrimaryTargetOwner.ID),
		SecondaryTargetType:  stt,
		SecondaryTargetId:    sid,
		SecondaryTargetOwner: sto,
	}, nil
}

func AuditLogsToOAPI(als []*pacta.AuditLog) ([]*api.AuditLog, error) {
	return convAll(als, AuditLogToOAPI)
}
