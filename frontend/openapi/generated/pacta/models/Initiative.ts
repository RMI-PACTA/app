/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { InitiativeUserRelationship } from './InitiativeUserRelationship';
import type { Language } from './Language';
import type { PortfolioInitiativeMembershipPortfolio } from './PortfolioInitiativeMembershipPortfolio';

export type Initiative = {
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
    affiliation: string;
    /**
     * Additional information about the initiative
     */
    publicDescription: string;
    /**
     * Additional information about the initiative, for participants only
     */
    internalDescription: string;
    /**
     * If set, only users who have been invited to join this initiative can join it, otherwise, anyone can join it. Defaults to false.
     */
    requiresInvitationToJoin: boolean;
    /**
     * If set, new users can join the initiative. Defaults to false.
     */
    isAcceptingNewMembers: boolean;
    /**
     * If set, users that are members of this initiative can add portfolios to it.
     */
    isAcceptingNewPortfolios: boolean;
    /**
     * The language this initiative should be conducted in.
     */
    language: Language;
    /**
     * The pacta model that this initiative should use, if not specified, the default pacta model will be used.
     */
    pactaVersion?: string;
    /**
     * the list of portfolios that are members of this initiative
     */
    portfolioInitiativeMemberships: Array<PortfolioInitiativeMembershipPortfolio>;
    /**
     * the list of users that are members of this initiative
     */
    initiativeUserRelationships: Array<InitiativeUserRelationship>;
    /**
     * The time at which this initiative was created.
     */
    createdAt: string;
};

