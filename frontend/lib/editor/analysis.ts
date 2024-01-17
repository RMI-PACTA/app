import { type Analysis } from '@/openapi/generated/pacta'
import { Validation, type EditorFieldsFor, type EditorValuesFor, type EditorComputedValues } from './common'
import { getEditorComputedValues, type Translation } from './utils'

export type EditorAnalysisFields = EditorFieldsFor<Analysis>
export type EditorAnalysisValues = EditorValuesFor<Analysis>

const createEditorAnalysisFields = (translation: Translation): EditorAnalysisFields => {
  const tt = (key: string) => translation.t(`lib/editor/analysis.${key}`)
  return {
    id: {
      name: 'id',
      label: tt('ID'),
    },
    analysisType: {
      name: 'analysisType',
      label: tt('Analysis Type'),
    },
    pactaVersion: {
      name: 'pactaVersion',
      label: tt('PACTA Version'),
    },
    portfolioSnapshot: {
      name: 'portfolioSnapshot',
      label: tt('Portfolio Snapshot'),
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
    failureCode: {
      name: 'failureCode',
      label: tt('Failure Code'),
    },
    failureMessage: {
      name: 'failureMessage',
      label: tt('Failure Message'),
    },
    artifacts: {
      name: 'artifacts',
      label: tt('Artifacts'),
    },
  }
}

export const analysisEditor = (i: Analysis, translation: Translation): EditorComputedValues<Analysis> => {
  return getEditorComputedValues(`lib/editor/analysis[${i.id}]`, i, createEditorAnalysisFields, translation)
}
