/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type InitiativeUserRelationship = {
    /**
     * the inititative that this relationship describes
     */
    initiativeId: string;
    /**
     * the user that this relationship describes
     */
    userId: string;
    /**
     * whether this user is a manager of the initiative
     */
    manager: boolean;
    /**
     * whether this user is a member of the initiative
     */
    member: boolean;
    /**
     * the time at which this relationship was last updated
     */
    updatedAt: string;
};

