/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type User = {
    /**
     * Unique id of the user
     */
    id: string;
    /**
     * User's email address as entered
     */
    enteredEmail: string;
    /**
     * Stanard formatting of the email address of the user
     */
    canonicalEmail?: string;
    /**
     * Whether the user is an administrator of the PACTA platform
     */
    admin: boolean;
    /**
     * Whether the user is an administrator of the PACTA platform
     */
    superAdmin: boolean;
    /**
     * Name of the user
     */
    name: string;
    /**
     * The user's preferred language, if present
     */
    preferredLanguage: User.preferredLanguage;
};

export namespace User {

    /**
     * The user's preferred language, if present
     */
    export enum preferredLanguage {
        EN = 'en',
        FR = 'fr',
        ES = 'es',
        DE = 'de',
    }


}

