package pactasrv

import (
	"context"
	"fmt"

	"github.com/RMI/pacta/oapierr"
	"github.com/RMI/pacta/pacta"
	"go.uber.org/zap"
)

func (s *Server) populatePortfoliosInPortfolioGroups(
	ctx context.Context,
	ts []*pacta.PortfolioGroup,
) error {
	getFn := func(pg *pacta.PortfolioGroup) ([]*pacta.Portfolio, error) {
		result := []*pacta.Portfolio{}
		for _, member := range pg.PortfolioGroupMemberships {
			result = append(result, member.Portfolio)
		}
		return result, nil
	}
	lookupFn := func(ids []pacta.PortfolioID) (map[pacta.PortfolioID]*pacta.Portfolio, error) {
		return s.DB.Portfolios(s.DB.NoTxn(ctx), ids)
	}
	getIDFn := func(p *pacta.Portfolio) pacta.PortfolioID {
		return p.ID
	}
	if err := populateAll(ts, getFn, getIDFn, lookupFn); err != nil {
		return oapierr.Internal("populating portfolios in portfolio groups failed", zap.Error(err))
	}
	return nil
}

func (s *Server) populateInitiativesInPortfolios(
	ctx context.Context,
	is []*pacta.Portfolio,
) error {
	getFn := func(pg *pacta.Portfolio) ([]*pacta.Initiative, error) {
		result := []*pacta.Initiative{}
		for _, member := range pg.PortfolioInitiativeMemberships {
			result = append(result, member.Initiative)
		}
		return result, nil
	}
	lookupFn := func(ids []pacta.InitiativeID) (map[pacta.InitiativeID]*pacta.Initiative, error) {
		return s.DB.Initiatives(s.DB.NoTxn(ctx), ids)
	}
	getIDFn := func(p *pacta.Initiative) pacta.InitiativeID {
		return p.ID
	}
	if err := populateAll(is, getFn, getIDFn, lookupFn); err != nil {
		return oapierr.Internal("populating initiatives in portfolios failed", zap.Error(err))
	}
	return nil
}

func (s *Server) populatePortfolioGroupsInPortfolios(
	ctx context.Context,
	ts []*pacta.Portfolio,
) error {
	getFn := func(pg *pacta.Portfolio) ([]*pacta.PortfolioGroup, error) {
		result := []*pacta.PortfolioGroup{}
		for _, member := range pg.PortfolioGroupMemberships {
			result = append(result, member.PortfolioGroup)
		}
		return result, nil
	}
	lookupFn := func(ids []pacta.PortfolioGroupID) (map[pacta.PortfolioGroupID]*pacta.PortfolioGroup, error) {
		return s.DB.PortfolioGroups(s.DB.NoTxn(ctx), ids)
	}
	getIDFn := func(p *pacta.PortfolioGroup) pacta.PortfolioGroupID {
		return p.ID
	}
	if err := populateAll(ts, getFn, getIDFn, lookupFn); err != nil {
		return oapierr.Internal("populating portfolio groups in portfolios failed", zap.Error(err))
	}
	return nil
}

func (s *Server) populateArtifactsInAnalyses(
	ctx context.Context,
	ts ...*pacta.Analysis,
) error {
	getFn := func(a *pacta.Analysis) ([]*pacta.AnalysisArtifact, error) {
		result := []*pacta.AnalysisArtifact{}
		for _, aa := range a.Artifacts {
			result = append(result, aa)
		}
		return result, nil
	}
	lookupFn := func(ids []pacta.AnalysisArtifactID) (map[pacta.AnalysisArtifactID]*pacta.AnalysisArtifact, error) {
		return s.DB.AnalysisArtifacts(s.DB.NoTxn(ctx), ids)
	}
	getIDFn := func(a *pacta.AnalysisArtifact) pacta.AnalysisArtifactID {
		return a.ID
	}
	if err := populateAll(ts, getFn, getIDFn, lookupFn); err != nil {
		return oapierr.Internal("populating analysis artifacts in analysis failed", zap.Error(err))
	}
	return nil
}

func (s *Server) populateBlobsInAnalysisArtifacts(
	ctx context.Context,
	ts ...*pacta.AnalysisArtifact,
) error {
	getFn := func(a *pacta.AnalysisArtifact) ([]*pacta.Blob, error) {
		result := []*pacta.Blob{}
		if a.Blob != nil {
			result = append(result, a.Blob)
		}
		return result, nil
	}
	lookupFn := func(ids []pacta.BlobID) (map[pacta.BlobID]*pacta.Blob, error) {
		return s.DB.Blobs(s.DB.NoTxn(ctx), ids)
	}
	getIDFn := func(a *pacta.Blob) pacta.BlobID {
		return a.ID
	}
	if err := populateAll(ts, getFn, getIDFn, lookupFn); err != nil {
		return oapierr.Internal("populating blobs in analysis artifacts failed", zap.Error(err))
	}
	return nil
}

// This helper function populates the given targets in the given sources,
// to allow for generic population of nested data structures.
// sources = entities that you want to populate sub-entity references in.
// the sub-entities should be pointers to structs with an ID populated.
// getTargetsFn = function that takes a source and returns zero or more sub-entities to populate.
// getTargetIDFn = function that takes a sub-entity and returns its ID.
// lookupTargetsFn = function that takes a list of sub-entity IDs and returns a map of ID -> sub-entity.
func populateAll[Source any, TargetID ~string, Target any](
	sources []*Source,
	getTargetsFn func(*Source) ([]*Target, error),
	getTargetIDFn func(*Target) TargetID,
	lookupTargetsFn func([]TargetID) (map[TargetID]*Target, error),
) error {
	allTargets := []*Target{}
	for i, source := range sources {
		targets, err := getTargetsFn(source)
		if err != nil {
			return fmt.Errorf("getting %d-th targets: %w", i, err)
		}
		allTargets = append(allTargets, targets...)
	}
	if len(allTargets) == 0 {
		return nil
	}

	seen := map[TargetID]bool{}
	uniqueIds := []TargetID{}
	for _, target := range allTargets {
		id := getTargetIDFn(target)
		if _, ok := seen[id]; !ok {
			uniqueIds = append(uniqueIds, id)
			seen[id] = true
		}
	}

	populatedTargets, err := lookupTargetsFn(uniqueIds)
	if err != nil {
		return fmt.Errorf("looking up %Ts: %w", allTargets[0], err)
	}
	for i, source := range sources {
		targets, err := getTargetsFn(source)
		if err != nil {
			return fmt.Errorf("re-getting %d-th targets: %w", i, err)
		}
		for _, target := range targets {
			id := getTargetIDFn(target)
			if populated, ok := populatedTargets[id]; ok {
				*target = *populated
			} else {
				return fmt.Errorf("can't find populated target %s", id)
			}
		}
	}
	return nil
}
