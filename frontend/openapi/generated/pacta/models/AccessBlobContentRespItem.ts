/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type AccessBlobContentRespItem = {
    /**
     * The id of the blob to that the content is for.
     */
    blob_id: string;
    /**
     * The signed URL where the file can be downloaded from, using GET semantics.
     */
    download_url: string;
    /**
     * The time at which the signed URL will expire.
     */
    expiration_time: string;
};

