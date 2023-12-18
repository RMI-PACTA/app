/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { PortfolioGroupMembershipPortfolio } from './PortfolioGroupMembershipPortfolio';

export type PortfolioGroup = {
    /**
     * the system assigned id of the portfolio group
     */
    id: string;
    /**
     * the human meaningful name of the portfoio group
     */
    name: string;
    /**
     * the description of the contents or purpose of the portfolio group
     */
    description: string;
    /**
     * the list of portfolios that are members of this portfolio group
     */
    members?: Array<PortfolioGroupMembershipPortfolio>;
    /**
     * The time at which this initiative was created.
     */
    createdAt: string;
};

