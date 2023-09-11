package pactasrv

import (
	"fmt"

	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
)

func initiativeToOAPI(i *pacta.Initiative) (*api.Initiative, error) {
	if i == nil {
		return nil, errorInternal(fmt.Errorf("initiativeToOAPI: can't convert nil pointer"))
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

func pactaVersionToOAPI(pv *pacta.PACTAVersion) (*api.PactaVersion, error) {
	if pv == nil {
		return nil, errorInternal(fmt.Errorf("pactaVersionToOAPI: can't convert nil pointer"))
	}
	return &api.PactaVersion{
		Id:          string(pv.ID),
		Name:        pv.Name,
		IsDefault:   pv.IsDefault,
		Digest:      pv.Digest,
		Description: pv.Description,
		CreatedAt:   pv.CreatedAt,
	}, nil
}

func ptr[T any](t T) *T {
	return &t
}
