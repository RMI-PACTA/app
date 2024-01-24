/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { FailureCode } from './FailureCode';
import type { HoldingsDate } from './HoldingsDate';
import type { OptionalBoolean } from './OptionalBoolean';

export type IncompleteUpload = {
    /**
     * Unique identifier for the incomplete upload
     */
    id: string;
    /**
     * Name of the upload
     */
    name: string;
    /**
     * Description of the upload
     */
    description: string;
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
    /**
     * The time when the upload was created
     */
    createdAt: string;
    /**
     * The time when the upload process was run
     */
    ranAt?: string;
    /**
     * The time when the upload was completed
     */
    completedAt?: string;
    /**
     * Code describing the failure, if any
     */
    failureCode?: FailureCode;
    /**
     * Message describing the failure, if any
     */
    failureMessage?: string;
    /**
     * Flag to indicate whether admin debug mode is enabled
     */
    adminDebugEnabled: boolean;
};

