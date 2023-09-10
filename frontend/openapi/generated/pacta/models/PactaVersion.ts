/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type PactaVersion = {
    /**
     * Unique id of the pacta version - system assigned
     */
    id: string;
    /**
     * the human meaningful name of the version of the PACTA model
     */
    name: string;
    /**
     * Additional information about the version of the PACTA model
     */
    description: string;
    /**
     * The hash (typically SHA256) that uniquely identifies this version of the PACTA model.
     */
    digest: string;
    /**
     * The time at which this version of the PACTA model was created
     */
    createdAt: string;
    /**
     * Whether this version of the PACTA model is the default version
     */
    isDefault: boolean;
};

