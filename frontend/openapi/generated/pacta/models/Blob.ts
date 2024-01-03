/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { FileType } from './FileType';

export type Blob = {
    /**
     * the system assigned unique identifier of the blob
     */
    id: string;
    /**
     * the human meaningful name of the file
     */
    fileName: string;
    /**
     * the type (extension) of the file
     */
    fileType: FileType;
    /**
     * The time at which this blob was created within the system
     */
    createdAt: string;
};

