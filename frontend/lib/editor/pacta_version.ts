import { type PactaVersion } from '@/openapi/generated/pacta'
import { type EditorFieldsFor, Validation, type EditorComputedValues } from './common'
import { getEditorComputedValues } from './utils'

export type EditorPactaVersion = EditorFieldsFor<PactaVersion>

const createEditorPactaVersion = (pactaVersion: PactaVersion): EditorPactaVersion => {
  return {
    id: {
      name: 'id',
      label: 'ID',
      originalValue: pactaVersion.id,
      currentValue: pactaVersion.id,
    },
    name: {
      name: 'name',
      label: 'Name',
      validation: [Validation.NotEmpty],
      originalValue: pactaVersion.name,
      currentValue: pactaVersion.name,
    },
    description: {
      name: 'description',
      label: 'Description',
      validation: [Validation.NotEmpty],
      originalValue: pactaVersion.description,
      currentValue: pactaVersion.description,
    },
    digest: {
      name: 'digest',
      label: 'Docker Image Digest',
      validation: [Validation.NotEmpty],
      originalValue: pactaVersion.digest,
      currentValue: pactaVersion.digest,
    },
    isDefault: {
      name: 'isDefault',
      label: 'Is Default Version',
      originalValue: pactaVersion.isDefault,
      currentValue: pactaVersion.isDefault,
    },
    createdAt: {
      name: 'createdAt',
      label: 'Created At',
      originalValue: pactaVersion.createdAt,
      currentValue: pactaVersion.createdAt,
    },
  }
}

export const pactaVersionEditor = (pv: PactaVersion): EditorComputedValues<PactaVersion> => {
  return getEditorComputedValues(
    'lib/editor/pacta-version', pv, createEditorPactaVersion)
}
