package pactasrv

import (
	"context"
	"fmt"

	"github.com/RMI/pacta/oapierr"
	"github.com/RMI/pacta/pacta"
	"go.uber.org/zap"
)

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
		return fmt.Errorf("looking up populated: %w", err)
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

func populatePortfoliosInPortfolioGroups(
	s *Server,
	ctx context.Context,
	ts []*pacta.PortfolioGroup,
) error {
	getFn := func(pg *pacta.PortfolioGroup) ([]*pacta.Portfolio, error) {
		result := []*pacta.Portfolio{}
		for _, member := range pg.Members {
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

func populatePortfolioGroupsInPortfolios(
	s *Server,
	ctx context.Context,
	ts []*pacta.Portfolio,
) error {
	getFn := func(pg *pacta.Portfolio) ([]*pacta.PortfolioGroup, error) {
		result := []*pacta.PortfolioGroup{}
		for _, member := range pg.MemberOf {
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
