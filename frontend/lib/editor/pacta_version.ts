import { type PactaVersion } from '@/openapi/generated/pacta'
import { type EditorFieldsFor, type EditorValuesFor, Validation, type EditorComputedValues } from './common'
import { getEditorComputedValues, type Translation } from './utils'

export type EditorPactaVersionFields = EditorFieldsFor<PactaVersion>
export type EditorPactaVersionValues = EditorValuesFor<PactaVersion>

const createEditorPactaVersionFields = (translation: Translation): EditorPactaVersionFields => {
  const tt = (key: string) => translation.t(`lib/editor/pacta-version.${key}`)
  return {
    id: {
      name: 'id',
      label: tt('ID'),
    },
    name: {
      name: 'name',
      label: tt('Name'),
      validation: [Validation.NotEmpty],
      helpText: tt('The name of the version of the PACTA algorithm.'),
    },
    description: {
      name: 'description',
      label: tt('Description'),
      validation: [Validation.NotEmpty],
      helpText: tt('An optional description of this version of the PACTA algorithm.'),
    },
    digest: {
      name: 'digest',
      label: tt('Docker Image Digest'),
      validation: [Validation.NotEmpty],
      helpText: tt('The SHA hash of the docker image that should correspond to this version of the PACTA version.'),
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

export const pactaVersionEditor = (pv: PactaVersion, translation: Translation): EditorComputedValues<PactaVersion> => {
  return getEditorComputedValues(
    'lib/editor/pacta_version',
    pv,
    createEditorPactaVersionFields,
    translation)
}
