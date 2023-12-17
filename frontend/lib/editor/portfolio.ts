import { type Portfolio } from '@/openapi/generated/pacta'
import { Validation, type EditorFieldsFor, type EditorValuesFor, type EditorComputedValues } from './common'
import { getEditorComputedValues, type Translation } from './utils'

export type EditorPortfolioFields = EditorFieldsFor<Portfolio>
export type EditorPortfolioValues = EditorValuesFor<Portfolio>

const createEditorPortfolioFields = (translation: Translation): EditorPortfolioFields => {
  const tt = (key: string) => translation.t(`lib/editor/portfolio/${key}`)
  return {
    id: {
      name: 'id',
      label: tt('ID'),
    },
    name: {
      name: 'name',
      label: tt('Name'),
      validation: [Validation.NotEmpty],
      helpText: tt('The name of this portfolio.'),
    },
    description: {
      name: 'description',
      label: tt('Description'),
      helpText: tt('The description of this portfolio - helpful for record keeping, not used for anything besides organization.'),
    },
    adminDebugEnabled: {
      name: 'adminDebugEnabled',
      label: tt('Admin Debugging Enabled'),
      helpText: tt('When enabled, this portfolio can be accessed by administrators to help with debugging. Only turn this on if you\'re comfortable with system administrators accessing this data.'),
    },
    holdingsDate: {
      name: 'holdingsDate',
      label: tt('Holdings Date'),
    },
    createdAt: {
      name: 'createdAt',
      label: tt('Created At'),
    },
    numberOfRows: {
      name: 'numberOfRows',
      label: tt('Number of Rows'),
    },
  }
}

export const portfolioEditor = (i: Portfolio, translation: Translation): EditorComputedValues<Portfolio> => {
  return getEditorComputedValues(`lib/editor/portfolio[${i.id}]`, i, createEditorPortfolioFields, translation)
}
