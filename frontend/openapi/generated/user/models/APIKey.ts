/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type APIKey = {
    /**
     * Unique identifier for the API key
     */
    id: string;
    /**
     * An opaque string used to authenticate with various RMI APIs
     */
    key: string;
    /**
     * Timestamp when the token expires, RFC3339-formatted.
     */
    expiresAt?: string;
};

