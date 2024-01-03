/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Language } from './Language';

export type InitiativeCreate = {
    /**
     * the human readable identifier for the initiative, can only include alphanumeric characters, dashes and underscores
     */
    id: string;
    /**
     * the human meaningful name of the version of the initiative
     */
    name: string;
    /**
     * the group that sponsors/created/owns this initiative
     */
    affiliation?: string;
    /**
     * Additional information about the initiative
     */
    publicDescription?: string;
    /**
     * Additional information about the initiative, for participants only
     */
    internalDescription?: string;
    /**
     * If set, only users who have been invited to join this initiative can join it, otherwise, anyone can join it. Defaults to false.
     */
    requiresInvitationToJoin?: boolean;
    /**
     * If set, new users can join the initiative. Defaults to false.
     */
    isAcceptingNewMembers?: boolean;
    /**
     * If set, users that are members of this initiative can add portfolios to it.
     */
    isAcceptingNewPortfolios?: boolean;
    /**
     * The language this initiative should be conducted in.
     */
    language: Language;
    /**
     * The id of the PACTA model that this initiative should use, if not specified, the default PACTA model will be used.
     */
    pactaVersion?: string;
};

