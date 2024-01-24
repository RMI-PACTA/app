/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type MergeUsersResp = {
    /**
     * the number of incomplete uploads that were transferred to the new user
     */
    incompleteUploadCount: number;
    /**
     * the number of analyses that were transferred to the new user
     */
    analysisCount: number;
    /**
     * the number of portfolios that were transferred to the new user
     */
    portfolioCount: number;
    /**
     * the number of portfolio groups that were transferred to the new user
     */
    portfolioGroupCount: number;
    /**
     * the number of audit logs that were created to record the merge
     */
    auditLogsCreated: number;
};

