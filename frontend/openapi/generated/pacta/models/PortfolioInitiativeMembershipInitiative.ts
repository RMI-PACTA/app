/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Initiative } from './Initiative';

export type PortfolioInitiativeMembershipInitiative = {
    /**
     * the user that added this portfolio to the initiative
     */
    addedByUserId?: string;
    initiative: Initiative;
    /**
     * The time at which this relationship was created
     */
    createdAt: string;
};

