import { type Initiative } from '@/openapi/generated/pacta'
import { Validation, type EditorFieldsFor, type EditorComputedValues } from './common'
import { getEditorComputedValues } from './utils'

export type EditorInitiative = EditorFieldsFor<Initiative>

const createEditorInitiative = (initiative: Initiative): EditorInitiative => {
  return {
    id: {
      name: 'id',
      label: 'ID',
      validation: [Validation.AlphanumericAndDashesAndUnderscores],
      originalValue: initiative.id,
      currentValue: initiative.id,
    },
    name: {
      name: 'name',
      label: 'Name',
      validation: [Validation.NotEmpty],
      originalValue: initiative.name,
      currentValue: initiative.name,
    },
    affiliation: {
      name: 'affiliation',
      label: 'Affiliation',
      originalValue: initiative.affiliation,
      currentValue: initiative.affiliation,
    },
    publicDescription: {
      name: 'publicDescription',
      label: 'Public Description',
      validation: [Validation.NotEmpty],
      originalValue: initiative.publicDescription,
      currentValue: initiative.publicDescription,
    },
    internalDescription: {
      name: 'internalDescription',
      label: 'Internal Description',
      originalValue: initiative.internalDescription,
      currentValue: initiative.internalDescription,
    },
    requiresInvitationToJoin: {
      name: 'requiresInvitationToJoin',
      label: 'Requires Invitation to Join',
      originalValue: initiative.requiresInvitationToJoin,
      currentValue: initiative.requiresInvitationToJoin,
    },
    isAcceptingNewMembers: {
      name: 'isAcceptingNewMembers',
      label: 'Accepting New Members',
      originalValue: initiative.isAcceptingNewMembers,
      currentValue: initiative.isAcceptingNewMembers,
    },
    isAcceptingNewPortfolios: {
      name: 'isAcceptingNewPortfolios',
      label: 'Accepting New Portfolios',
      originalValue: initiative.isAcceptingNewPortfolios,
      currentValue: initiative.isAcceptingNewPortfolios,
    },
    language: {
      name: 'language',
      label: 'Language',
      validation: [Validation.NotEmpty],
      originalValue: initiative.language,
      currentValue: initiative.language,
    },
    pactaVersion: {
      name: 'pactaVersion',
      label: 'PACTA Version',
      validation: [Validation.NotEmpty],
      originalValue: initiative.pactaVersion,
      currentValue: initiative.pactaVersion,
    },
    createdAt: {
      name: 'createdAt',
      label: 'Created At',
      originalValue: initiative.createdAt,
      currentValue: initiative.createdAt,
    },
  }
}

export const initiativeEditor = (i: Initiative): EditorComputedValues<Initiative> => {
  return getEditorComputedValues(
    'lib/editor/initiative', i, createEditorInitiative)
}

interface Example {
  a: string
  b: string | undefined
  c?: string
  d?: string | undefined
}

export type ExampleEditor = EditorFieldsFor<Example>
