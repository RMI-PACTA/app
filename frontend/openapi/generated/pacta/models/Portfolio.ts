/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { HoldingsDate } from './HoldingsDate';
import type { PortfolioGroupMembershipPortfolioGroup } from './PortfolioGroupMembershipPortfolioGroup';

export type Portfolio = {
    /**
     * the system assigned unique identifier of the portfolio
     */
    id: string;
    /**
     * the human meaningful name of the portfolio
     */
    name: string;
    /**
     * Additional information about the portfolio
     */
    description: string;
    /**
     * The time at which this portfolio was successfully parsed from a raw
     */
    createdAt: string;
    holdingsDate?: HoldingsDate;
    /**
     * Whether the admin debug mode is enabled for this portfolio
     */
    adminDebugEnabled: boolean;
    /**
     * The number of rows in the portfolio
     */
    numberOfRows: number;
    /**
     * The list of portfolio groups that this portfolio is a member of
     */
    groups?: Array<PortfolioGroupMembershipPortfolioGroup>;
};

