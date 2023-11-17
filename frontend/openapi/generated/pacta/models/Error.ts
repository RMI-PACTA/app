/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type Error = {
    /**
     * Human readable error message (in English)
     */
    message: string;
    /**
     * An enum-like type indicating a more specific type of error.
     *
     * An example might be getting a 401 Unauthorized because you're logged in with multiple emails and haven't selected one, the error_id could be 'multiple_emails'.
     */
    error_id: string;
};

