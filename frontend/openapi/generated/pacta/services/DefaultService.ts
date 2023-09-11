/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { EmptySuccess } from '../models/EmptySuccess';
import type { Error } from '../models/Error';
import type { Initiative } from '../models/Initiative';
import type { InitiativeChanges } from '../models/InitiativeChanges';
import type { InitiativeCreate } from '../models/InitiativeCreate';
import type { PactaVersion } from '../models/PactaVersion';
import type { PactaVersionChanges } from '../models/PactaVersionChanges';
import type { PactaVersionCreate } from '../models/PactaVersionCreate';
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
            errors: {
                401: `the user is not authorized to access this resource - if logged out, try logging in`,
                403: `the user is not authorized to access this resource`,
            },
        });
    }

    /**
     * Updates a PACTA version
     * Updates a PACTA version's settable properties
     * @param id ID of PACTA version to update
     * @param requestBody PACTA Version object properties to update
     * @returns EmptySuccess pacta version updated successfully
     * @returns Error unexpected error
     * @throws ApiError
     */
    public updatePactaVersion(
        id: string,
        requestBody: PactaVersionChanges,
    ): CancelablePromise<EmptySuccess | Error> {
        return this.httpRequest.request({
            method: 'PATCH',
            url: '/pacta-version/{id}',
            path: {
                'id': id,
            },
            body: requestBody,
            mediaType: 'application/json',
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
        });
    }

    /**
     * Marks this version of the PACTA model as the default
     * @param id ID of pacta version to fetch
     * @returns EmptySuccess updated successfully
     * @returns Error unexpected error
     * @throws ApiError
     */
    public markPactaVersionAsDefault(
        id: string,
    ): CancelablePromise<EmptySuccess | Error> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/pacta-version/{id}/set-default',
            path: {
                'id': id,
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
     * @param requestBody PACTA Version object properties to update
     * @returns EmptySuccess pacta version created successfully
     * @returns Error unexpected error
     * @throws ApiError
     */
    public createPactaVersion(
        requestBody: PactaVersionCreate,
    ): CancelablePromise<EmptySuccess | Error> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/pacta-versions',
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Returns an initiative by ID
     * @param id ID of the initiative to fetch
     * @returns Initiative the initiative requested
     * @returns Error unexpected error
     * @throws ApiError
     */
    public findInitiativeById(
        id: string,
    ): CancelablePromise<Initiative | Error> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/initiative/{id}',
            path: {
                'id': id,
            },
        });
    }

    /**
     * Updates an initiative
     * Updates an initiative's settable properties
     * @param id ID of the initiative to update
     * @param body initiative object properties to update
     * @returns EmptySuccess initiative updated successfully
     * @returns Error unexpected error
     * @throws ApiError
     */
    public updateInitiative(
        id: string,
        body: InitiativeChanges,
    ): CancelablePromise<EmptySuccess | Error> {
        return this.httpRequest.request({
            method: 'PATCH',
            url: '/initiative/{id}',
            path: {
                'id': id,
            },
            query: {
                'body': body,
            },
        });
    }

    /**
     * Deletes an initiative by id
     * deletes an initiative based on the ID supplied
     * @param id ID of initiative to delete
     * @returns EmptySuccess initiative deleted successfully
     * @returns Error unexpected error
     * @throws ApiError
     */
    public deleteInitiative(
        id: string,
    ): CancelablePromise<EmptySuccess | Error> {
        return this.httpRequest.request({
            method: 'DELETE',
            url: '/initiative/{id}',
            path: {
                'id': id,
            },
        });
    }

    /**
     * Returns all initiatives
     * @returns Initiative gets all initiatives
     * @returns Error unexpected error
     * @throws ApiError
     */
    public listInitiatives(): CancelablePromise<Array<Initiative> | Error> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/initiatives',
        });
    }

    /**
     * Creates a initiative
     * Creates a new initiative
     * @param requestBody Initiative object properties to update
     * @returns EmptySuccess initiative created successfully
     * @returns Error unexpected error
     * @throws ApiError
     */
    public createInitiative(
        requestBody: InitiativeCreate,
    ): CancelablePromise<EmptySuccess | Error> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/initiatives',
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * gets info about the logged in user
     * Returns the logged in user, if the user is logged in, otherwise returns empty
     * @returns User user response
     * @returns Error unexpected error
     * @throws ApiError
     */
    public findUserByMe(): CancelablePromise<User | Error> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/user/me',
            errors: {
                401: `caller is not logged in - log in to continue`,
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
        });
    }

    /**
     * Updates user properties
     * Updates a user's settable properties
     * @param id ID of user to update
     * @param requestBody User object properties to update
     * @returns EmptySuccess the new user object
     * @returns Error unexpected error
     * @throws ApiError
     */
    public updateUser(
        id: string,
        requestBody: UserChanges,
    ): CancelablePromise<EmptySuccess | Error> {
        return this.httpRequest.request({
            method: 'PATCH',
            url: '/user/{id}',
            path: {
                'id': id,
            },
            body: requestBody,
            mediaType: 'application/json',
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
        });
    }

}
