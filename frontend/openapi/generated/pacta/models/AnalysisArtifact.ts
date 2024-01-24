/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Blob } from './Blob';

export type AnalysisArtifact = {
    /**
     * the system assigned unique identifier of the artifact
     */
    id: string;
    /**
     * Whether the admin debug mode is enabled for this artifact
     */
    adminDebugEnabled: boolean;
    /**
     * Whether this artifact is publicly accessible
     */
    sharedToPublic: boolean;
    /**
     * Information about the file/artifact itself
     */
    blob: Blob;
};

