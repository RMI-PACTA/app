import { type Initiative } from '@/openapi/generated/pacta'
import { Validation, type EditorFieldsFor, type EditorValuesFor, type EditorComputedValues } from './common'
import { getEditorComputedValues, type Translation } from './utils'

export type EditorInitiativeFields = EditorFieldsFor<Initiative>
export type EditorInitiativeValues = EditorValuesFor<Initiative>

const createEditorInitiativeFields = (translation: Translation): EditorInitiativeFields => {
  const tt = (key: string) => translation.t(`lib/editor/initiative.${key}`)
  return {
    id: {
      name: 'id',
      label: tt('ID'),
      validation: [Validation.AlphanumericAndDashesAndUnderscores],
      helpText: tt('This is the immutable unique identifier for the initiative. It can only contain alphanumeric characters, underscores, and dashes. This value will be shown in URLs, but will typically not be user visible.'),
    },
    name: {
      name: 'name',
      label: tt('Name'),
      validation: [Validation.NotEmpty],
      helpText: tt('The name of the PACTA initiative.'),
    },
    affiliation: {
      name: 'affiliation',
      label: tt('Affiliation'),
      helpText: tt('An optional description of the organization or entity that is hosting this initiative.'),
    },
    publicDescription: {
      name: 'publicDescription',
      label: tt('Public Description'),
      validation: [Validation.NotEmpty],
      helpText: tt('The description of the initiative that will be shown to the public. Newlines will be respected.'),
    },
    internalDescription: {
      name: 'internalDescription',
      label: tt('Internal Description'),
      helpText: tt('The description of the initiative that will be shown to members of the inititiative. Newlines will be respected.'),
    },
    requiresInvitationToJoin: {
      name: 'requiresInvitationToJoin',
      label: tt('Requires Invitation to Join'),
      helpText: tt('When disabled, anyone can join this initiative. When enabled, initiative administrators can mint invitation codes that they can share with folks to allow them to join the project.'),
    },
    isAcceptingNewMembers: {
      name: 'isAcceptingNewMembers',
      label: tt('Accepting New Members'),
      helpText: tt('When enabled, new members can join the project through the joining mechanism selected above.'),
    },
    isAcceptingNewPortfolios: {
      name: 'isAcceptingNewPortfolios',
      label: tt('Accepting New Portfolios'),
      helpText: tt('When enabled, initiative members can add new portfolios to the initiative.'),
    },
    language: {
      name: 'language',
      label: tt('Language'),
      validation: [Validation.NotEmpty],
      helpText: tt('What language should reports have when they are generated for this initiative?'),
    },
    pactaVersion: {
      name: 'pactaVersion',
      label: tt('PACTA Version'),
      validation: [Validation.NotEmpty],
      helpText: tt('What version of the PACTA algorithm should this initiative use to generate reports?'),
    },
    createdAt: {
      name: 'createdAt',
      label: tt('Created At'),
    },
  }
}

export const initiativeEditor = (i: Initiative, translation: Translation): EditorComputedValues<Initiative> => {
  return getEditorComputedValues('lib/editor/initiative', i, createEditorInitiativeFields, translation)
}
