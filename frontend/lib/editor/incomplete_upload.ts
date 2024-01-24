import { type IncompleteUpload } from '@/openapi/generated/pacta'
import { Validation, type EditorFieldsFor, type EditorValuesFor, type EditorComputedValues } from './common'
import { getEditorComputedValues, type Translation } from './utils'

export type EditorIncompleteUploadFields = EditorFieldsFor<IncompleteUpload>
export type EditorIncompleteUploadValues = EditorValuesFor<IncompleteUpload>

const createEditorIncompleteUploadFields = (translation: Translation): EditorIncompleteUploadFields => {
  const tt = (key: string) => translation.t(`lib/editor/incomplete_upload.${key}`)
  return {
    id: {
      name: 'id',
      label: tt('ID'),
    },
    name: {
      name: 'name',
      label: tt('Name'),
      validation: [Validation.NotEmpty],
      helpText: tt('NameHelpText'),
    },
    description: {
      name: 'description',
      label: tt('Description'),
      helpText: tt('DescriptionHelpText'),
    },
    adminDebugEnabled: {
      name: 'adminDebugEnabled',
      label: tt('AdminDebuggingEnabled'),
      helpText: tt('ADEHelpText'),
    },
    propertyHoldingsDate: {
      name: 'propertyHoldingsDate',
      label: tt('Holdings Date'),
      helpText: tt('HoldingsDateHelpText'),
    },
    propertyESG: {
      name: 'propertyESG',
      label: tt('ESG'),
      helpText: tt('ESGHelpText'),
    },
    propertyExternal: {
      name: 'propertyExternal',
      label: tt('External'),
      helpText: tt('ExternalHelpText'),
    },
    propertyEngagementStrategy: {
      name: 'propertyEngagementStrategy',
      label: tt('Engagement Strategy'),
      helpText: tt('EngagementStrategyHelpText'),
    },
    createdAt: {
      name: 'createdAt',
      label: tt('Created At'),
    },
    ranAt: {
      name: 'ranAt',
      label: tt('Ran At'),
    },
    completedAt: {
      name: 'completedAt',
      label: tt('Completed At'),
    },
    failureMessage: {
      name: 'failureMessage',
      label: tt('Failure Message'),
    },
    failureCode: {
      name: 'failureCode',
      label: tt('Failure Code'),
    },
  }
}

export const incompleteUploadEditor = (i: IncompleteUpload, translation: Translation): EditorComputedValues<IncompleteUpload> => {
  return getEditorComputedValues(`lib/editor/incomplete_upload[${i.id}]`, i, createEditorIncompleteUploadFields, translation)
}
