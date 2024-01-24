/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { AccessBlobContentReq } from '../models/AccessBlobContentReq';
import type { AccessBlobContentResp } from '../models/AccessBlobContentResp';
import type { Analysis } from '../models/Analysis';
import type { AnalysisArtifactChanges } from '../models/AnalysisArtifactChanges';
import type { AnalysisChanges } from '../models/AnalysisChanges';
import type { AuditLogQueryReq } from '../models/AuditLogQueryReq';
import type { AuditLogQueryResp } from '../models/AuditLogQueryResp';
import type { CompletePortfolioUploadReq } from '../models/CompletePortfolioUploadReq';
import type { CompletePortfolioUploadResp } from '../models/CompletePortfolioUploadResp';
import type { FindUserByMeResp } from '../models/FindUserByMeResp';
import type { IncompleteUpload } from '../models/IncompleteUpload';
import type { IncompleteUploadChanges } from '../models/IncompleteUploadChanges';
import type { Initiative } from '../models/Initiative';
import type { InitiativeAllData } from '../models/InitiativeAllData';
import type { InitiativeChanges } from '../models/InitiativeChanges';
import type { InitiativeCreate } from '../models/InitiativeCreate';
import type { InitiativeInvitation } from '../models/InitiativeInvitation';
import type { InitiativeInvitationCreate } from '../models/InitiativeInvitationCreate';
import type { InitiativeUserRelationship } from '../models/InitiativeUserRelationship';
import type { InitiativeUserRelationshipChanges } from '../models/InitiativeUserRelationshipChanges';
import type { ListAnalysesResp } from '../models/ListAnalysesResp';
import type { ListIncompleteUploadsResp } from '../models/ListIncompleteUploadsResp';
import type { ListPortfolioGroupsResp } from '../models/ListPortfolioGroupsResp';
import type { ListPortfoliosResp } from '../models/ListPortfoliosResp';
import type { MergeUsersReq } from '../models/MergeUsersReq';
import type { MergeUsersResp } from '../models/MergeUsersResp';
import type { PactaVersion } from '../models/PactaVersion';
import type { PactaVersionChanges } from '../models/PactaVersionChanges';
import type { PactaVersionCreate } from '../models/PactaVersionCreate';
import type { Portfolio } from '../models/Portfolio';
import type { PortfolioChanges } from '../models/PortfolioChanges';
import type { PortfolioGroup } from '../models/PortfolioGroup';
import type { PortfolioGroupChanges } from '../models/PortfolioGroupChanges';
import type { PortfolioGroupCreate } from '../models/PortfolioGroupCreate';
import type { PortfolioGroupMembershipIds } from '../models/PortfolioGroupMembershipIds';
import type { RunAnalysisReq } from '../models/RunAnalysisReq';
import type { RunAnalysisResp } from '../models/RunAnalysisResp';
import type { StartPortfolioUploadReq } from '../models/StartPortfolioUploadReq';
import type { StartPortfolioUploadResp } from '../models/StartPortfolioUploadResp';
import type { User } from '../models/User';
import type { UserChanges } from '../models/UserChanges';
import type { UserQueryReq } from '../models/UserQueryReq';
import type { UserQueryResp } from '../models/UserQueryResp';

import type { CancelablePromise } from '../core/CancelablePromise';
import type { BaseHttpRequest } from '../core/BaseHttpRequest';

export class DefaultService {

    constructor(public readonly httpRequest: BaseHttpRequest) {}

    /**
     * Gives the caller access to the blob
     * Checks whether the user can access the blobs, and if so, returns blob download URLs for each, generating an audit log along the way
     * @param requestBody Information about the blobs that are requested
     * @returns AccessBlobContentResp the user can access the blobs, and the access URLs are returned, along with information about their expiration
     * @throws ApiError
     */
    public accessBlobContent(
        requestBody: AccessBlobContentReq,
    ): CancelablePromise<AccessBlobContentResp> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/access-blob-content',
            body: requestBody,
            mediaType: 'application/json',
        });
    }

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
     * Merges two users together
     * Merges two users together
     * @param requestBody a request describing the two users to merge
     * @returns MergeUsersResp the users were merged successfully
     * @throws ApiError
     */
    public mergeUsers(
        requestBody: MergeUsersReq,
    ): CancelablePromise<MergeUsersResp> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/admin/merge-users',
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
     * Returns all of the portfolios that are participating in the initiative
     * @param id ID of the initiative to fetch data for
     * @returns InitiativeAllData the initiative data
     * @throws ApiError
     */
    public allInitiativeData(
        id: string,
    ): CancelablePromise<InitiativeAllData> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/initiative/{id}/all-data',
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
     * Returns all initiative user relationships for this user that the caller has access to view
     * @param userId ID of the user to fetch relationships for
     * @returns InitiativeUserRelationship
     * @throws ApiError
     */
    public listInitiativeUserRelationshipsByUser(
        userId: string,
    ): CancelablePromise<Array<InitiativeUserRelationship>> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/initiative/user-relationships/{userId}',
            path: {
                'userId': userId,
            },
        });
    }

    /**
     * Returns all initiative invitations associated with the initiative
     * @param id ID of the initiative to fetch invitations for
     * @returns InitiativeInvitation
     * @throws ApiError
     */
    public listInitiativeInvitations(
        id: string,
    ): CancelablePromise<Array<InitiativeInvitation>> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/initiative/{id}/invitations',
            path: {
                'id': id,
            },
        });
    }

    /**
     * Creates an initiative invitation
     * Creates an initiative invitation
     * @param requestBody
     * @returns void
     * @throws ApiError
     */
    public createInitiativeInvitation(
        requestBody: InitiativeInvitationCreate,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/initiative-invitation',
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Returns the initiative invitation from this id, if it exists
     * @param id ID of the invitation to fetch details about
     * @returns InitiativeInvitation
     * @throws ApiError
     */
    public getInitiativeInvitation(
        id: string,
    ): CancelablePromise<InitiativeInvitation> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/initiative-invitation/{id}',
            path: {
                'id': id,
            },
        });
    }

    /**
     * Claims this initiative invitation, if it exists
     * @param id ID of the invitation to claim
     * @returns void
     * @throws ApiError
     */
    public claimInitiativeInvitation(
        id: string,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/initiative-invitation/{id}:claim',
            path: {
                'id': id,
            },
            errors: {
                409: `initiative invitation already claimed`,
            },
        });
    }

    /**
     * Deletes an initiative invitation by id
     * deletes an initiative based on the ID supplied
     * @param id ID of initiative invitation to delete
     * @returns void
     * @throws ApiError
     */
    public deleteInitiativeInvitation(
        id: string,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'DELETE',
            url: '/initiative-invitation/{id}:claim',
            path: {
                'id': id,
            },
        });
    }

    /**
     * Updates initiative user relationship properties
     * Updates a given user's relationship properties for a given initiative
     * @param initiativeId ID of the initiative
     * @param userId ID of the user
     * @param requestBody Relationship object properties to update
     * @returns void
     * @throws ApiError
     */
    public updateInitiativeUserRelationship(
        initiativeId: string,
        userId: string,
        requestBody: InitiativeUserRelationshipChanges,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'PATCH',
            url: '/initiative/{initiativeId}/user-relationship/{userId}',
            path: {
                'initiativeId': initiativeId,
                'userId': userId,
            },
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * creates an initiative portfolio relationship
     * creates a membership relationship between the portfolio and the initiative
     * @param initiativeId ID of the initiative
     * @param portfolioId ID of the portfolio
     * @returns void
     * @throws ApiError
     */
    public createInitiativePortfolioRelationship(
        initiativeId: string,
        portfolioId: string,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/initiative/{initiativeId}/portfolio-relationship/{portfolioId}',
            path: {
                'initiativeId': initiativeId,
                'portfolioId': portfolioId,
            },
        });
    }

    /**
     * Deletes an initiative:portfolio relationship
     * Deletes a given portfolio's relationship with a given initiative
     * @param initiativeId ID of the initiative
     * @param portfolioId ID of the portfolio
     * @returns void
     * @throws ApiError
     */
    public deleteInitiativePortfolioRelationship(
        initiativeId: string,
        portfolioId: string,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'DELETE',
            url: '/initiative/{initiativeId}/portfolio-relationship/{portfolioId}',
            path: {
                'initiativeId': initiativeId,
                'portfolioId': portfolioId,
            },
        });
    }

    /**
     * Returns the portfolio groups that the user has access to
     * @returns ListPortfolioGroupsResp
     * @throws ApiError
     */
    public listPortfolioGroups(): CancelablePromise<ListPortfolioGroupsResp> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/portfolio-groups',
        });
    }

    /**
     * Creates a portfolio group
     * Creates a new portfolio group
     * @param requestBody Initial portfolio group object properties
     * @returns PortfolioGroup portfolio group created successfully
     * @throws ApiError
     */
    public createPortfolioGroup(
        requestBody: PortfolioGroupCreate,
    ): CancelablePromise<PortfolioGroup> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/portfolio-groups',
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Returns a portfolio group by ID
     * Returns a portfolio group based on a single ID
     * @param id ID of portfolio group to fetch
     * @returns PortfolioGroup portfolio group response
     * @throws ApiError
     */
    public findPortfolioGroupById(
        id: string,
    ): CancelablePromise<PortfolioGroup> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/portfolio-group/{id}',
            path: {
                'id': id,
            },
        });
    }

    /**
     * Updates portfolio group properties
     * Updates a portfolio group's settable properties
     * @param id ID of the portfolio group to update
     * @param requestBody Portfolio Group object properties to update
     * @returns void
     * @throws ApiError
     */
    public updatePortfolioGroup(
        id: string,
        requestBody: PortfolioGroupChanges,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'PATCH',
            url: '/portfolio-group/{id}',
            path: {
                'id': id,
            },
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Deletes a portfolio group by ID
     * deletes a portfolio group based on the ID supplied - note this does not delete the portfolios that are members to this group
     * @param id ID of portfolio group to delete
     * @returns void
     * @throws ApiError
     */
    public deletePortfolioGroup(
        id: string,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'DELETE',
            url: '/portfolio-group/{id}',
            path: {
                'id': id,
            },
        });
    }

    /**
     * creates a portfolio group membership
     * creates a portfolio group membership
     * @param requestBody Portfolio Group membership to create
     * @returns void
     * @throws ApiError
     */
    public createPortfolioGroupMembership(
        requestBody: PortfolioGroupMembershipIds,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'PUT',
            url: '/portfolio-group-membership',
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Deletes a portfolio group membership
     * removes the membership of a portfolio in a portfolio - note this does not delete the portfolio or the portfolio group
     * @param requestBody Portfolio Group membership to delete
     * @returns void
     * @throws ApiError
     */
    public deletePortfolioGroupMembership(
        requestBody: PortfolioGroupMembershipIds,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'DELETE',
            url: '/portfolio-group-membership',
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Gets the incomplete uploads that the user is the owner of
     * @returns ListIncompleteUploadsResp
     * @throws ApiError
     */
    public listIncompleteUploads(): CancelablePromise<ListIncompleteUploadsResp> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/incomplete-uploads',
        });
    }

    /**
     * Returns an incomplete upload by ID
     * Returns an incomplete upload based on a single ID
     * @param id ID of incomplete upload to fetch
     * @returns IncompleteUpload incomplete upload response
     * @throws ApiError
     */
    public findIncompleteUploadById(
        id: string,
    ): CancelablePromise<IncompleteUpload> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/incomplete-upload/{id}',
            path: {
                'id': id,
            },
        });
    }

    /**
     * Updates incomplete upload properties
     * Updates a incomplete upload's settable properties
     * @param id ID of incomplete upload to update
     * @param requestBody Incomplete Upload object properties to update
     * @returns void
     * @throws ApiError
     */
    public updateIncompleteUpload(
        id: string,
        requestBody: IncompleteUploadChanges,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'PATCH',
            url: '/incomplete-upload/{id}',
            path: {
                'id': id,
            },
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Deletes an incomplete upload by ID
     * deletes an incomplete upload based on the ID supplied
     * @param id ID of incomplete upload to delete
     * @returns void
     * @throws ApiError
     */
    public deleteIncompleteUpload(
        id: string,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'DELETE',
            url: '/incomplete-upload/{id}',
            path: {
                'id': id,
            },
        });
    }

    /**
     * Gets the list of portfolios that the user is the owner of
     * @returns ListPortfoliosResp
     * @throws ApiError
     */
    public listPortfolios(): CancelablePromise<ListPortfoliosResp> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/portfolios',
        });
    }

    /**
     * Returns an portfolio by ID
     * Returns an portfolio based on a single ID
     * @param id ID of portfolio to fetch
     * @returns Portfolio portfolio response
     * @throws ApiError
     */
    public findPortfolioById(
        id: string,
    ): CancelablePromise<Portfolio> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/portfolio/{id}',
            path: {
                'id': id,
            },
        });
    }

    /**
     * Updates portfolio properties
     * Updates a portfolio's settable properties
     * @param id ID of portfolio to update
     * @param requestBody portfolio object properties to update
     * @returns void
     * @throws ApiError
     */
    public updatePortfolio(
        id: string,
        requestBody: PortfolioChanges,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'PATCH',
            url: '/portfolio/{id}',
            path: {
                'id': id,
            },
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Deletes an portfolio by ID
     * deletes an portfolio based on the ID supplied
     * @param id ID of portfolio to delete
     * @returns void
     * @throws ApiError
     */
    public deletePortfolio(
        id: string,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'DELETE',
            url: '/portfolio/{id}',
            path: {
                'id': id,
            },
        });
    }

    /**
     * gets info about the logged in user
     * Returns the logged in user, if the user is logged in, otherwise returns empty
     * @returns FindUserByMeResp user response
     * @throws ApiError
     */
    public findUserByMe(): CancelablePromise<FindUserByMeResp> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/user/me',
        });
    }

    /**
     * Gets the list of users that the user is able to view, currently an admin-only action
     * @param requestBody A request describing which users should be returned
     * @returns UserQueryResp
     * @throws ApiError
     */
    public userQuery(
        requestBody: UserQueryReq,
    ): CancelablePromise<UserQueryResp> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/users',
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * a callback after login to create or return the user
     * Creates a user in the database, if the user does not yet exist, or returns the existing user.
     * @returns User user response
     * @throws ApiError
     */
    public userAuthenticationFollowup(): CancelablePromise<User> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/user/authentication-followup',
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

    /**
     * queries the platform's audit logs
     * returns back audit logs that matc the user's query
     * @param requestBody A request describing which audit logs should be returned
     * @returns AuditLogQueryResp The audit logs that matched the requested query, if any
     * @throws ApiError
     */
    public listAuditLogs(
        requestBody: AuditLogQueryReq,
    ): CancelablePromise<AuditLogQueryResp> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/audit-logs',
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Starts the process of uploading one or more portfolio files
     * Creates one or more new incomplete portfolio uploads, and creates upload URLs for the user to put their blobs into.
     * @param requestBody A request describing the portfolios that the user wants to upload
     * @returns StartPortfolioUploadResp The assets can now be uploaded via the given signed URLs.
     * @throws ApiError
     */
    public startPortfolioUpload(
        requestBody: StartPortfolioUploadReq,
    ): CancelablePromise<StartPortfolioUploadResp> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/portfolio-upload',
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Called after uploads of portfolios to cloud storage are complete.
     * Signals that the upload of the portfolios are complete, and that the server should start parsing them.
     * @param requestBody A request describing the incomplete uploads that the user wants to begin processing
     * @returns CompletePortfolioUploadResp The process to initiate the parsing of the uploads has been initiated.
     * @throws ApiError
     */
    public completePortfolioUpload(
        requestBody: CompletePortfolioUploadReq,
    ): CancelablePromise<CompletePortfolioUploadResp> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/portfolio-upload:complete',
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Gets the analyses that the user is the owner of
     * @returns ListAnalysesResp
     * @throws ApiError
     */
    public listAnalyses(): CancelablePromise<ListAnalysesResp> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/analyses',
        });
    }

    /**
     * Returns an analysis by ID
     * Returns an analysis based on a single ID
     * @param id ID of analysis to fetch
     * @returns Analysis analysis response
     * @throws ApiError
     */
    public findAnalysisById(
        id: string,
    ): CancelablePromise<Analysis> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/analysis/{id}',
            path: {
                'id': id,
            },
        });
    }

    /**
     * Updates writable analysis properties
     * Updates an analysis' settable properties
     * @param id ID of analysis to update
     * @param requestBody Analayis object properties to update
     * @returns void
     * @throws ApiError
     */
    public updateAnalysis(
        id: string,
        requestBody: AnalysisChanges,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'PATCH',
            url: '/analysis/{id}',
            path: {
                'id': id,
            },
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Deletes an analysis (and its artifacts) by ID
     * deletes an analysis based on the ID supplied
     * @param id ID of analysis to delete
     * @returns void
     * @throws ApiError
     */
    public deleteAnalysis(
        id: string,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'DELETE',
            url: '/analysis/{id}',
            path: {
                'id': id,
            },
        });
    }

    /**
     * Updates writable analysis artifact properties
     * Updates an analysis artifact's settable properties
     * @param id ID of analysis artifact to update
     * @param requestBody Analysis artifact's object properties to update
     * @returns void
     * @throws ApiError
     */
    public updateAnalysisArtifact(
        id: string,
        requestBody: AnalysisArtifactChanges,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'PATCH',
            url: '/analysis-artifact/{id}',
            path: {
                'id': id,
            },
            body: requestBody,
            mediaType: 'application/json',
        });
    }

    /**
     * Deletes an analysis artifact by ID
     * deletes an analysis artifact based on the ID supplied
     * @param id ID of analysis artifact to delete
     * @returns void
     * @throws ApiError
     */
    public deleteAnalysisArtifact(
        id: string,
    ): CancelablePromise<void> {
        return this.httpRequest.request({
            method: 'DELETE',
            url: '/analysis-artifact/{id}',
            path: {
                'id': id,
            },
        });
    }

    /**
     * Requests an anslysis be run
     * Creates a snapshot of the requested entity, and starts it running
     * @param requestBody Properties of the analysis to run
     * @returns RunAnalysisResp information about the requested analysis
     * @throws ApiError
     */
    public runAnalysis(
        requestBody: RunAnalysisReq,
    ): CancelablePromise<RunAnalysisResp> {
        return this.httpRequest.request({
            method: 'POST',
            url: '/run-analysis',
            body: requestBody,
            mediaType: 'application/json',
        });
    }

}
