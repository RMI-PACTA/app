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
      helpText: tt('IDHelpText'),
    },
    name: {
      name: 'name',
      label: tt('Name'),
      validation: [Validation.NotEmpty],
      helpText: tt('NameHelpText'),
    },
    affiliation: {
      name: 'affiliation',
      label: tt('Affiliation'),
      helpText: tt('AffiliationHelpText'),
    },
    publicDescription: {
      name: 'publicDescription',
      label: tt('Public Description'),
      validation: [Validation.NotEmpty],
      helpText: tt('PublicDescriptionHelpText'),
    },
    internalDescription: {
      name: 'internalDescription',
      label: tt('Internal Description'),
      helpText: tt('InternalDescriptionHelpText'),
    },
    requiresInvitationToJoin: {
      name: 'requiresInvitationToJoin',
      label: tt('Requires Invitation to Join'),
      helpText: tt('RequiresInvitationToJoinHelpText'),
    },
    isAcceptingNewMembers: {
      name: 'isAcceptingNewMembers',
      label: tt('Accepting New Members'),
      helpText: tt('AcceptingNewMembersHelpText'),
    },
    isAcceptingNewPortfolios: {
      name: 'isAcceptingNewPortfolios',
      label: tt('Accepting New Portfolios'),
      helpText: tt('AcceptingNewPortfoliosHelpText'),
    },
    language: {
      name: 'language',
      label: tt('Language'),
      validation: [Validation.NotEmpty],
      helpText: tt('LanguageHelpText'),
    },
    pactaVersion: {
      name: 'pactaVersion',
      label: tt('PACTA Version'),
      validation: [Validation.NotEmpty],
      helpText: tt('PactaVersionHelpText'),
    },
    createdAt: {
      name: 'createdAt',
      label: tt('Created At'),
    },
    portfolioInitiativeMemberships: {
      name: 'portfolioInitiativeMemberships',
      label: tt('Portfolio Initiative Memberships'),
    },
    initiativeUserRelationships: {
      name: 'initiativeUserRelationships',
      label: tt('Initiative User Relationships'),
    },
  }
}

export const initiativeEditor = (i: Initiative, translation: Translation): EditorComputedValues<Initiative> => {
  return getEditorComputedValues('lib/editor/initiative', i, createEditorInitiativeFields, translation)
}
