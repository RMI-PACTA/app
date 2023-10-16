package conv

import (
	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
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

func ptr[T any](t T) *T {
	return &t
}
