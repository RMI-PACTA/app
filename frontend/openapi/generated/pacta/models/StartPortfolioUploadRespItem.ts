/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type StartPortfolioUploadRespItem = {
    /**
     * The name of the file, including its extension, used as a round-trip id.
     */
    file_name: string;
    /**
     * The signed URL where the file should be uploaded to, using PUT semantics.
     */
    upload_url: string;
    /**
     * A unique identifier for the uploaded asset
     */
    incomplete_upload_id: string;
};

