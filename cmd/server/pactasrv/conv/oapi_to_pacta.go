package conv

import (
	"regexp"

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

func ifNil[T any](t *T, fallback T) T {
	if t == nil {
		return fallback
	}
	return *t
}
