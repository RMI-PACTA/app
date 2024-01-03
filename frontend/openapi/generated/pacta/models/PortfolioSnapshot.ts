/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Initiative } from './Initiative';
import type { Portfolio } from './Portfolio';
import type { PortfolioGroup } from './PortfolioGroup';

/**
 * represents an immutable description of a collection of portfolios at a point in time, used to ensure reproducibility and change detection
 */
export type PortfolioSnapshot = {
    /**
     * the system assigned unique identifier of the snapshot
     */
    id: string;
    /**
     * the full set of denormalized portfolios included in this analysis
     */
    portfolioIds: Array<string>;
    /**
     * if populated, this snapshot represents a snapshot of this initiative
     */
    initiative?: Initiative;
    /**
     * if populated, this snapshot represents a snapshot of this portfolio group
     */
    portfolioGroup?: PortfolioGroup;
    /**
     * if populated, this snapshot represents a snapshot of this solitary portfolio
     */
    portfolio?: Portfolio;
};

