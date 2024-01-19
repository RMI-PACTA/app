/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { User } from './User';

export type FindUserByMeResp = {
    user?: User;
    /**
     * the id of the owner of the user, if logged in
     */
    ownerId?: string;
};

