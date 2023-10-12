/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
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
     * @throws ApiError
     */
    public findPactaVersionById(
        id: string,
    ): CancelablePromise<PactaVersion> {
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
     * @param requestBody PACTA Version object properties to update
     * @returns void
     * @throws ApiError
     */
    public updatePactaVersion(
        id: string,
        requestBody: PactaVersionChanges,
    ): CancelablePromise<void> {
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
     * @returns void
     * @throws ApiError
     */
    public deletePactaVersion(
        id: string,
    ): CancelablePromise<void> {
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
     * @returns void
     * @throws ApiError
     */
    public markPactaVersionAsDefault(
        id: string,
    ): CancelablePromise<void> {
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
     * @throws ApiError
     */
    public listPactaVersions(): CancelablePromise<Array<PactaVersion>> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/pacta-versions',
        });
    }

    /**
     * Creates a PACTA version
     * Creates a PACTA version
     * @param requestBody PACTA Version object properties to update
     * @returns void
     * @throws ApiError
     */
    public createPactaVersion(
        requestBody: PactaVersionCreate,
    ): CancelablePromise<void> {
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
     * @throws ApiError
     */
    public findInitiativeById(
        id: string,
    ): CancelablePromise<Initiative> {
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
     * @param requestBody initiative object properties to update
     * @returns void
     * @throws ApiError
     */
    public updateInitiative(
        id: string,
        requestBody: InitiativeChanges,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'PATCH',
            url: '/initiative/{id}',
            path: {
                'id': id,
            },
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Deletes an initiative by id
     * deletes an initiative based on the ID supplied
     * @param id ID of initiative to delete
     * @returns void
     * @throws ApiError
     */
    public deleteInitiative(
        id: string,
    ): CancelablePromise<void> {
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
     * @throws ApiError
     */
    public listInitiatives(): CancelablePromise<Array<Initiative>> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/initiatives',
        });
    }

    /**
     * Creates a initiative
     * Creates a new initiative
     * @param requestBody Initiative object properties to update
     * @returns void
     * @throws ApiError
     */
    public createInitiative(
        requestBody: InitiativeCreate,
    ): CancelablePromise<void> {
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
     * @throws ApiError
     */
    public findUserByMe(): CancelablePromise<User> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/user/me',
        });
    }

    /**
     * Returns a user by ID
     * Returns a user based on a single ID
     * @param id ID of user to fetch
     * @returns User user response
     * @throws ApiError
     */
    public findUserById(
        id: string,
    ): CancelablePromise<User> {
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
     * @returns void
     * @throws ApiError
     */
    public updateUser(
        id: string,
        requestBody: UserChanges,
    ): CancelablePromise<void> {
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
     * @returns void
     * @throws ApiError
     */
    public deleteUser(
        id: string,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'DELETE',
            url: '/user/{id}',
            path: {
                'id': id,
            },
        });
    }

}
