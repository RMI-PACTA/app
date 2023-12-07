import { type IncompleteUpload } from '@/openapi/generated/pacta'
import { Validation, type EditorFieldsFor, type EditorComputedValues } from './common'
import { getEditorComputedValues } from './utils'

export type EditorIncompleteUpload = EditorFieldsFor<IncompleteUpload>

const createEditorIncompleteUpload = (incompleteUpload: IncompleteUpload): EditorIncompleteUpload => {
  return {
    id: {
      name: 'id',
      label: 'ID',
      originalValue: incompleteUpload.id,
      currentValue: incompleteUpload.id,
    },
    name: {
      name: 'name',
      label: 'Name',
      validation: [Validation.NotEmpty],
      originalValue: incompleteUpload.name,
      currentValue: incompleteUpload.name,
    },
    description: {
      name: 'description',
      label: 'Description',
      originalValue: incompleteUpload.description,
      currentValue: incompleteUpload.description,
    },
    adminDebugEnabled: {
      name: 'adminDebugEnabled',
      label: 'Admin Debugging Enabled',
      originalValue: incompleteUpload.adminDebugEnabled,
      currentValue: incompleteUpload.adminDebugEnabled,
    },
    holdingsDate: {
      name: 'holdingsDate',
      label: 'Holdings Date',
      originalValue: incompleteUpload.holdingsDate,
      currentValue: incompleteUpload.holdingsDate,
    },
    createdAt: {
      name: 'createdAt',
      label: 'Created At',
      originalValue: incompleteUpload.createdAt,
      currentValue: incompleteUpload.createdAt,
    },
    ranAt: {
      name: 'ranAt',
      label: 'Ran At',
      originalValue: incompleteUpload.ranAt,
      currentValue: incompleteUpload.ranAt,
    },
    completedAt: {
      name: 'completedAt',
      label: 'Completed At',
      originalValue: incompleteUpload.completedAt,
      currentValue: incompleteUpload.completedAt,
    },
    failureMessage: {
      name: 'failureMessage',
      label: 'Failure Message',
      originalValue: incompleteUpload.failureMessage,
      currentValue: incompleteUpload.failureMessage,
    },
    failureCode: {
      name: 'failureCode',
      label: 'Failure Code',
      originalValue: incompleteUpload.failureCode,
      currentValue: incompleteUpload.failureCode,
    },
  }
}

export const incompleteUploadEditor = (i: IncompleteUpload): EditorComputedValues<IncompleteUpload> => {
  return getEditorComputedValues(`lib/editor/incompleteUpload[${i.id}]`, i, createEditorIncompleteUpload)
}
