/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type MergeUsersReq = {
    /**
     * the user id of the user to merge records from, and to be deleted after the merge
     */
    fromUserId: string;
    /**
     * the user id of the user to recieve merged records and to exist after the merge
     */
    toUserId: string;
};

