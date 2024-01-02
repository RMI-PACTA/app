/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { AuditLogAction } from './AuditLogAction';
import type { AuditLogActorType } from './AuditLogActorType';
import type { AuditLogTargetType } from './AuditLogTargetType';

export type AuditLog = {
    /**
     * the unique identifier of a given audit log
     */
    id: string;
    /**
     * the time that this audit log was created/the action was undertaken
     */
    createdAt: string;
    /**
     * the authority that this actor was acting as when performing this action
     */
    actorType: AuditLogActorType;
    /**
     * the user id of the actor that initiated this action, not populated if the system initiated the action
     */
    actorId?: string;
    /**
     * the owner id of the actor that initiated this action, not populated if the system initiated the action
     */
    actorOwnerId?: string;
    /**
     * the action that generated this audit log
     */
    action: AuditLogAction;
    /**
     * the object category that this action was performed on
     */
    primaryTargetType: AuditLogTargetType;
    /**
     * the id of the object that this action was performed on
     */
    primaryTargetId: string;
    /**
     * the id of the owner of the primary object this action was performed on
     */
    primaryTargetOwner: string;
    /**
     * the object category of the secondary object (membership partner, typically) that this action was performed on
     */
    secondaryTargetType?: AuditLogTargetType;
    /**
     * the id of the secondary object that this action was performed on
     */
    secondaryTargetId?: string;
    /**
     * the id of the owner of the secondary object this action was performed on
     */
    secondaryTargetOwner?: string;
};

