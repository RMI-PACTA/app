/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { HoldingsDate } from './HoldingsDate';
import type { OptionalBoolean } from './OptionalBoolean';
import type { PortfolioGroupMembershipPortfolioGroup } from './PortfolioGroupMembershipPortfolioGroup';
import type { PortfolioInitiativeMembershipInitiative } from './PortfolioInitiativeMembershipInitiative';

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
    /**
     * The list of initiatives that this portfolio is a member of
     */
    initiatives?: Array<PortfolioInitiativeMembershipInitiative>;
    propertyHoldingsDate: HoldingsDate;
    /**
     * If set, this portfolio represents ESG data
     */
    propertyESG: OptionalBoolean;
    /**
     * If set to false, this portfolio represents internal data, if set to false it represents external data, unset represents no user input
     */
    propertyExternal: OptionalBoolean;
    /**
     * If set, this portfolio represents engagement strategy data or not, if unset it represents no user input
     */
    propertyEngagementStrategy: OptionalBoolean;
};

