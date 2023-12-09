/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { HoldingsDate } from './HoldingsDate';

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
    holdingsDate?: HoldingsDate;
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
    failureCode?: string;
    /**
     * Message describing the failure, if any
     */
    failureMessage?: string;
    /**
     * Flag to indicate whether admin debug mode is enabled
     */
    adminDebugEnabled: boolean;
};

