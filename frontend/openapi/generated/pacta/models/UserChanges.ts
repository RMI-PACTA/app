/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Language } from './Language';

export type UserChanges = {
    /**
     * The new name of the user
     */
    name?: string;
    /**
     * The user's new preferred language
     */
    preferredLanguage?: Language;
    /**
     * Whether the given user is an admin
     */
    admin?: boolean;
    /**
     * Whether the given user is a super admin
     */
    superAdmin?: boolean;
};

