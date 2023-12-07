import { type PactaVersion } from '@/openapi/generated/pacta'
import { type EditorFieldsFor, type EditorValuesFor, Validation, type EditorComputedValues } from './common'
import { getEditorComputedValues, type Translation } from './utils'

export type EditorPactaVersionFields = EditorFieldsFor<PactaVersion>
export type EditorPactaVersionValues = EditorValuesFor<PactaVersion>

const createEditorPactaVersionFields = (translation: Translation): EditorPactaVersionFields => {
  const tt = (key: string) => translation.t(`lib/editor/pacta-version/${key}`)
  return {
    id: {
      name: 'id',
      label: tt('ID'),
    },
    name: {
      name: 'name',
      label: tt('Name'),
      validation: [Validation.NotEmpty],
    },
    description: {
      name: 'description',
      label: tt('Description'),
      validation: [Validation.NotEmpty],
    },
    digest: {
      name: 'digest',
      label: tt('Docker Image Digest'),
      validation: [Validation.NotEmpty],
    },
    isDefault: {
      name: 'isDefault',
      label: tt('Is Default Version'),
    },
    createdAt: {
      name: 'createdAt',
      label: tt('Created At'),
    },
  }
}

const createEditorPactaVersionValues = (pactaVersion: PactaVersion): EditorPactaVersionValues => {
  return {
    id: {
      originalValue: pactaVersion.id,
      currentValue: pactaVersion.id,
    },
    name: {
      originalValue: pactaVersion.name,
      currentValue: pactaVersion.name,
    },
    description: {
      originalValue: pactaVersion.description,
      currentValue: pactaVersion.description,
    },
    digest: {
      originalValue: pactaVersion.digest,
      currentValue: pactaVersion.digest,
    },
    isDefault: {
      originalValue: pactaVersion.isDefault,
      currentValue: pactaVersion.isDefault,
    },
    createdAt: {
      originalValue: pactaVersion.createdAt,
      currentValue: pactaVersion.createdAt,
    },
  }
}

export const pactaVersionEditor = (pv: PactaVersion, translation: Translation): EditorComputedValues<PactaVersion> => {
  return getEditorComputedValues(
    'lib/editor/pacta-version',
    pv,
    createEditorPactaVersionFields,
    createEditorPactaVersionValues,
    translation)
}
