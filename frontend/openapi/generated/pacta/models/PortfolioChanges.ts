/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { HoldingsDate } from './HoldingsDate';
import type { OptionalBoolean } from './OptionalBoolean';

export type PortfolioChanges = {
    /**
     * the human meaningful name of the portfolio
     */
    name?: string;
    /**
     * Additional information about the portfolio
     */
    description?: string;
    /**
     * Whether the admin debug mode is enabled for this portfolio
     */
    adminDebugEnabled?: boolean;
    propertyHoldingsDate?: HoldingsDate;
    /**
     * If set, this portfolio represents ESG data
     */
    propertyESG?: OptionalBoolean;
    /**
     * If set to false, this portfolio represents internal data, if set to false it represents external data, unset represents no user input
     */
    propertyExternal?: OptionalBoolean;
    /**
     * If set, this portfolio represents engagement strategy data or not, if unset it represents no user input
     */
    propertyEngagementStrategy?: OptionalBoolean;
};

