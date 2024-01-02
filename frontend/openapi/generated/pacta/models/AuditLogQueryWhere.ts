/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { AuditLogAction } from './AuditLogAction';
import type { AuditLogActorType } from './AuditLogActorType';
import type { AuditLogTargetType } from './AuditLogTargetType';

export type AuditLogQueryWhere = {
    /**
     * a list of audit log ids to filter by
     */
    inId?: Array<string>;
    /**
     * a minimum time for audit logs to filter audit logs by
     */
    minCreatedAt?: string;
    /**
     * a maximum time for audit logs to filter audit logs by
     */
    maxCreatedAt?: string;
    /**
     * a list of audit log action types to filter audit logs by
     */
    inAction?: Array<AuditLogAction>;
    /**
     * a list of audit log actor types to filter audit logs by
     */
    inActorType?: Array<AuditLogActorType>;
    /**
     * a list of actor user ids to filter audit logs by
     */
    inActorId?: Array<string>;
    /**
     * a list of actor owner ids to filter audit logs by
     */
    inActorOwnerId?: Array<string>;
    /**
     * a list of audit log target types to filter audit logs by
     */
    inTargetType?: Array<AuditLogTargetType>;
    /**
     * a list of target ids to filter audit logs by
     */
    inTargetId?: Array<string>;
    /**
     * a list of target owner ids to filter audit logs by
     */
    inTargetOwnerId?: Array<string>;
};

