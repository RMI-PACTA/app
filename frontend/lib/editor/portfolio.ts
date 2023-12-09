import { type Portfolio } from '@/openapi/generated/pacta'
import { Validation, type EditorFieldsFor, type EditorComputedValues } from './common'
import { getEditorComputedValues } from './utils'

export type EditorPortfolio = EditorFieldsFor<Portfolio>

const createEditorPortfolio = (portfolio: Portfolio): EditorPortfolio => {
  return {
    id: {
      name: 'id',
      label: 'ID',
      originalValue: portfolio.id,
      currentValue: portfolio.id,
    },
    name: {
      name: 'name',
      label: 'Name',
      validation: [Validation.NotEmpty],
      originalValue: portfolio.name,
      currentValue: portfolio.name,
    },
    description: {
      name: 'description',
      label: 'Description',
      originalValue: portfolio.description,
      currentValue: portfolio.description,
    },
    adminDebugEnabled: {
      name: 'adminDebugEnabled',
      label: 'Admin Debugging Enabled',
      originalValue: portfolio.adminDebugEnabled,
      currentValue: portfolio.adminDebugEnabled,
    },
    holdingsDate: {
      name: 'holdingsDate',
      label: 'Holdings Date',
      originalValue: portfolio.holdingsDate,
      currentValue: portfolio.holdingsDate,
    },
    createdAt: {
      name: 'createdAt',
      label: 'Created At',
      originalValue: portfolio.createdAt,
      currentValue: portfolio.createdAt,
    },
    numberOfRows: {
      name: 'numberOfRows',
      label: 'Number of Rows',
      originalValue: portfolio.numberOfRows,
      currentValue: portfolio.numberOfRows,
    },
  }
}

export const portfolioEditor = (i: Portfolio): EditorComputedValues<Portfolio> => {
  return getEditorComputedValues(`lib/editor/portfolio[${i.id}]`, i, createEditorPortfolio)
}
