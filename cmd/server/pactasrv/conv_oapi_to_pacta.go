package pactasrv

import (
	"fmt"
	"regexp"

	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
)

var initiativeIDRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func initiativeCreateToPACTA(i *api.InitiativeCreate) (*pacta.Initiative, error) {
	if i == nil {
		return nil, errorBadRequest("InitiativeCreate", "cannot be nil")
	}
	if !initiativeIDRegex.MatchString(i.Id) {
		return nil, errorBadRequest("id", "must contain only alphanumeric characters, underscores, and dashes")
	}
	lang, err := pacta.ParseLanguage(string(i.Language))
	if err != nil {
		return nil, errorBadRequest("language", err.Error())
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

func pactaVersionCreateToPACTA(p *api.PactaVersionCreate) (*pacta.PACTAVersion, error) {
	if p == nil {
		return nil, fmt.Errorf("pactaVersionCreateToPACTA: nil pointer")
	}
	return &pacta.PACTAVersion{
		Name:        p.Name,
		Digest:      p.Digest,
		Description: p.Description,
	}, nil
}

func ifNil[T any](t *T, or T) T {
	if t == nil {
		return or
	}
	return *t
}
