/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { EmptySuccess } from '../models/EmptySuccess';
import type { Error } from '../models/Error';
import type { PactaVersion } from '../models/PactaVersion';
import type { PactaVersionChanges } from '../models/PactaVersionChanges';
import type { User } from '../models/User';
import type { UserChanges } from '../models/UserChanges';

import type { CancelablePromise } from '../core/CancelablePromise';
import type { BaseHttpRequest } from '../core/BaseHttpRequest';

export class DefaultService {

    constructor(public readonly httpRequest: BaseHttpRequest) {}

    /**
     * Returns a version of the PACTA model by ID
     * @param id ID of pacta version to fetch
     * @returns PactaVersion pacta response
     * @returns Error unexpected error
     * @throws ApiError
     */
    public findPactaVersionById(
        id: string,
    ): CancelablePromise<PactaVersion | Error> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/pacta-version/{id}',
            path: {
                'id': id,
            },
        });
    }

    /**
     * Updates a PACTA version
     * Updates a PACTA version's settable properties
     * @param id ID of PACTA version to update
     * @param body PACTA Version object properties to update
     * @returns EmptySuccess pacta version updated successfully
     * @returns Error unexpected error
     * @throws ApiError
     */
    public updatePactaVersion(
        id: string,
        body: PactaVersionChanges,
    ): CancelablePromise<EmptySuccess | Error> {
        return this.httpRequest.request({
            method: 'PATCH',
            url: '/pacta-version/{id}',
            path: {
                'id': id,
            },
            query: {
                'body': body,
            },
            errors: {
                403: `caller does not have access or PACTA version does not exist`,
            },
        });
    }

    /**
     * Deletes a pacta version by ID
     * deletes a single pacta version based on the ID supplied
     * @param id ID of pacta version to delete
     * @returns EmptySuccess pacta version deleted successfully
     * @returns Error unexpected error
     * @throws ApiError
     */
    public deletePactaVersion(
        id: string,
    ): CancelablePromise<EmptySuccess | Error> {
        return this.httpRequest.request({
            method: 'DELETE',
            url: '/pacta-version/{id}',
            path: {
                'id': id,
            },
            errors: {
                403: `caller does not have access or pacta version does not exist`,
            },
        });
    }

    /**
     * Returns all versions of the PACTA model
     * @returns PactaVersion pacta versions
     * @returns Error unexpected error
     * @throws ApiError
     */
    public listPactaVersions(): CancelablePromise<Array<PactaVersion> | Error> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/pacta-versions',
        });
    }

    /**
     * Creates a PACTA version
     * Creates a PACTA version
     * @param body PACTA Version object properties to update
     * @returns EmptySuccess pacta version created successfully
     * @returns Error unexpected error
     * @throws ApiError
     */
    public createPactaVersion(
        body: PactaVersion,
    ): CancelablePromise<EmptySuccess | Error> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/pacta-versions',
            query: {
                'body': body,
            },
            errors: {
                403: `caller does not have access to create PACTA versions`,
            },
        });
    }

    /**
     * Returns a user by ID
     * Returns a user based on a single ID
     * @param id ID of user to fetch
     * @returns User user response
     * @returns Error unexpected error
     * @throws ApiError
     */
    public findUserById(
        id: string,
    ): CancelablePromise<User | Error> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/user/{id}',
            path: {
                'id': id,
            },
            errors: {
                403: `caller does not have access or user does not exist`,
            },
        });
    }

    /**
     * Updates user properties
     * Updates a user's settable properties
     * @param id ID of user to update
     * @param body User object properties to update
     * @returns User the new user object
     * @returns Error unexpected error
     * @throws ApiError
     */
    public updateUser(
        id: string,
        body: UserChanges,
    ): CancelablePromise<User | Error> {
        return this.httpRequest.request({
            method: 'PATCH',
            url: '/user/{id}',
            path: {
                'id': id,
            },
            query: {
                'body': body,
            },
            errors: {
                403: `caller does not have access or user does not exist`,
            },
        });
    }

    /**
     * Deletes a user by ID
     * deletes a single user based on the ID supplied
     * @param id ID of user to delete
     * @returns EmptySuccess user deleted
     * @returns Error unexpected error
     * @throws ApiError
     */
    public deleteUser(
        id: string,
    ): CancelablePromise<EmptySuccess | Error> {
        return this.httpRequest.request({
            method: 'DELETE',
            url: '/user/{id}',
            path: {
                'id': id,
            },
            errors: {
                403: `caller does not have access or user does not exist`,
            },
        });
    }

}
