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

func ptr[T any](t T) *T {
	return &t
}
