/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { AuditLog } from './AuditLog';

export type AuditLogQueryResp = {
    auditLogs: Array<AuditLog>;
    /**
     * describes whether there are more records to query
     */
    hasNextPage: boolean;
    /**
     * the parameter to re-request with to continue this query on the next page of results
     */
    cursor: string;
};

