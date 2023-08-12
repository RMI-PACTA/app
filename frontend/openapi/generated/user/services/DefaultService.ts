/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIKey } from '../models/APIKey';
import type { Error } from '../models/Error';

import type { CancelablePromise } from '../core/CancelablePromise';
import type { BaseHttpRequest } from '../core/BaseHttpRequest';

export class DefaultService {

    constructor(public readonly httpRequest: BaseHttpRequest) {}

    /**
     * Exchange a user JWT token for an auth cookie that can be used with other RMI APIs
     * Takes an auth system-issued JWT and returns a auth cookie.
     *
     * @returns string Cookie response
     * @returns Error unexpected error
     * @throws ApiError
     */
    public login(): CancelablePromise<string | Error> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/login/cookie',
            responseHeader: 'Set-Cookie',
            errors: {
                403: `User is not allowed to log in`,
            },
        });
    }

    /**
     * Log out a user from RMI APIs
     * Clears an existing API JWT
     *
     * @returns string Cookie response
     * @throws ApiError
     */
    public logout(): CancelablePromise<string> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/logout/cookie',
            responseHeader: 'Set-Cookie',
        });
    }

    /**
     * Exchange a user JWT token for an API key that can be used with other RMI APIs
     * Takes an auth system-issued JWT and returns a new API key.
     *
     * @returns APIKey API key response
     * @returns Error unexpected error
     * @throws ApiError
     */
    public createApiKey(): CancelablePromise<APIKey | Error> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/login/apikey',
            errors: {
                403: `User is not allowed to create an API key`,
            },
        });
    }

}
