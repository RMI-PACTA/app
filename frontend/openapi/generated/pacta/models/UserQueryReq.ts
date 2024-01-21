/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { UserQueryWhere } from './UserQueryWhere';

export type UserQueryReq = {
    /**
     * if provided, continues an existing query at the given point
     */
    cursor?: string;
    /**
     * the constraints to place on the returned records
     */
    wheres?: Array<UserQueryWhere>;
};

