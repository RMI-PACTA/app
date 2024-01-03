/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { AuditLogQuerySort } from './AuditLogQuerySort';
import type { AuditLogQueryWhere } from './AuditLogQueryWhere';

export type AuditLogQueryReq = {
    /**
     * if provided, continues an existing query at the given point
     */
    cursor?: string;
    /**
     * if provided, requests this number of records at maximum - default/maximum is 100
     */
    limit?: number;
    /**
     * the constraints to place on the returned records - this must be set to something which limits it to a scope the user should have access to
     */
    wheres: Array<AuditLogQueryWhere>;
    /**
     * the ordering that the results should be returned in - if empty, an ordering by created at date will be applied
     */
    sorts?: Array<AuditLogQuerySort>;
};

