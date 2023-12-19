import { type PortfolioGroup } from '@/openapi/generated/pacta'
import { Validation, type EditorFieldsFor, type EditorValuesFor, type EditorComputedValues } from './common'
import { getEditorComputedValues, type Translation } from './utils'

export type EditorPortfolioGroupFields = EditorFieldsFor<PortfolioGroup>
export type EditorPortfolioGroupValues = EditorValuesFor<PortfolioGroup>

const createEditorPortfolioGroupFields = (translation: Translation): EditorPortfolioGroupFields => {
  const tt = (key: string) => translation.t(`lib/editor/portfolio_group.${key}`)
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
    createdAt: {
      name: 'createdAt',
      label: tt('Created At'),
    },
    members: {
      name: 'members',
      label: tt('Members'),
    },
  }
}

export const portfolioGroupEditor = (i: PortfolioGroup, translation: Translation): EditorComputedValues<PortfolioGroup> => {
  return getEditorComputedValues(`lib/editor/portfolio_group[${i.id}]`, i, createEditorPortfolioGroupFields, translation)
}
