/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type UserChanges = {
    /**
     * The new name of the user
     */
    name?: string;
    /**
     * The user's new preferred language
     */
    preferredLanguage?: UserChanges.preferredLanguage;
    /**
     * Whether the given user is an admin
     */
    admin?: boolean;
    /**
     * Whether the given user is a super admin
     */
    superAdmin?: boolean;
};

export namespace UserChanges {

    /**
     * The user's new preferred language
     */
    export enum preferredLanguage {
        EN = 'en',
        FR = 'fr',
        ES = 'es',
        DE = 'de',
    }


}

