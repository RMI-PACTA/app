/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type InitiativeAllDataPortfolioItem = {
    /**
     * the name of the portfolio
     */
    name: string;
    /**
     * the id of the blob of the portfolio, which can be used to start a new partial download if the first download times out
     */
    blobId: string;
    /**
     * the url to download the portfolio
     */
    downloadUrl: string;
    /**
     * the time at which the download url will expire
     */
    expirationTime: string;
};

