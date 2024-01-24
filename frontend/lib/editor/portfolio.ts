import { type Portfolio } from '@/openapi/generated/pacta'
import { Validation, type EditorFieldsFor, type EditorValuesFor, type EditorComputedValues } from './common'
import { getEditorComputedValues, type Translation } from './utils'

export type EditorPortfolioFields = EditorFieldsFor<Portfolio>
export type EditorPortfolioValues = EditorValuesFor<Portfolio>

const createEditorPortfolioFields = (translation: Translation): EditorPortfolioFields => {
  const tt = (key: string) => translation.t(`lib/editor/portfolio.${key}`)
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
      label: tt('Admin Debugging Enabled'),
      helpText: tt('AdminDebuggingEnabledHelpText'),
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
    numberOfRows: {
      name: 'numberOfRows',
      label: tt('Number of Rows'),
    },
    groups: {
      name: 'groups',
      label: tt('Groups'),
    },
    initiatives: {
      name: 'initiatives',
      label: tt('Initiatives'),
    },
  }
}

export const portfolioEditor = (i: Portfolio, translation: Translation): EditorComputedValues<Portfolio> => {
  return getEditorComputedValues(`lib/editor/portfolio[${i.id}]`, i, createEditorPortfolioFields, translation)
}
