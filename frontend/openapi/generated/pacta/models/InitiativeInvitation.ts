/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type InitiativeInvitation = {
    /**
     * the human-readable id identifying this initiative invitation
     */
    id: string;
    /**
     * the id of the initiative that this invitation is for
     */
    initiativeId: string;
    /**
     * the time at which this initiative invitation was used, if it has been used
     */
    usedAt?: string;
    /**
     * the id of the user that used this initiative invitation, if it has been used
     */
    usedByUserId?: string;
    /**
     * the time at which this initiative invitation was created
     */
    createdAt: string;
};

