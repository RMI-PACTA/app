/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { HoldingsDate } from './HoldingsDate';
import type { OptionalBoolean } from './OptionalBoolean';
import type { StartPortfolioUploadReqItem } from './StartPortfolioUploadReqItem';

export type StartPortfolioUploadReq = {
    items: Array<StartPortfolioUploadReqItem>;
    propertyHoldingsDate?: HoldingsDate;
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

