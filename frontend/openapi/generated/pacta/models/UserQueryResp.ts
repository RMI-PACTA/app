/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { User } from './User';

export type UserQueryResp = {
    users: Array<User>;
    /**
     * describes whether there are more records to query
     */
    hasNextPage: boolean;
    /**
     * the parameter to re-request with to continue this query on the next page of results
     */
    cursor: string;
};

