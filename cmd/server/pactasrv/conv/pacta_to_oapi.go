package conv

import (
	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"go.uber.org/zap"
)

func InitiativeToOAPI(i *pacta.Initiative) (*api.Initiative, error) {
	if i == nil {
		return nil, oapierr.Internal("initiativeToOAPI: can't convert nil pointer")
	}
	return &api.Initiative{
		Affiliation:              i.Affiliation,
		CreatedAt:                i.CreatedAt,
		Id:                       string(i.ID),
		InternalDescription:      i.InternalDescription,
		IsAcceptingNewMembers:    i.IsAcceptingNewMembers,
		IsAcceptingNewPortfolios: i.IsAcceptingNewPortfolios,
		Language:                 api.InitiativeLanguage(i.Language),
		Name:                     i.Name,
		PactaVersion:             ptr(string(i.PACTAVersion.ID)),
		PublicDescription:        i.PublicDescription,
		RequiresInvitationToJoin: i.RequiresInvitationToJoin,
	}, nil
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
		usedBy = ptr(string(i.UsedBy.ID))
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
	memberOfs := []api.PortfolioGroupMembershipPortfolioGroup{}
	for _, m := range p.MemberOf {
		pg, err := PortfolioGroupToOAPI(m.PortfolioGroup)
		if err != nil {
			return nil, oapierr.Internal("portfolioToOAPI: portfolioGroupToOAPI failed", zap.Error(err))
		}
		memberOfs = append(memberOfs, api.PortfolioGroupMembershipPortfolioGroup{
			CreatedAt:      m.CreatedAt,
			PortfolioGroup: *pg,
		})
	}
	return &api.Portfolio{
		Id:                string(p.ID),
		Name:              p.Name,
		Description:       p.Description,
		HoldingsDate:      hd,
		CreatedAt:         p.CreatedAt,
		NumberOfRows:      p.NumberOfRows,
		AdminDebugEnabled: p.AdminDebugEnabled,
		Groups:            &memberOfs,
	}, nil
}

func PortfolioGroupToOAPI(pg *pacta.PortfolioGroup) (*api.PortfolioGroup, error) {
	if pg == nil {
		return nil, oapierr.Internal("portfolioGroupToOAPI: can't convert nil pointer")
	}
	members := []api.PortfolioGroupMembershipPortfolio{}
	for _, m := range pg.Members {
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
