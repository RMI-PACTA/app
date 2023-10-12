import { type Initiative, type InitiativeChanges } from '@/openapi/generated/pacta'
import { type EditorField, asChange, asValue, asIncompleteField, Validation } from './common'
import { type Ref } from 'vue'

export interface EditorInitiative {
  id: EditorField<Initiative, 'id'>
  name: EditorField<Initiative, 'name'>
  affiliation: EditorField<Initiative, 'affiliation'>
  publicDescription: EditorField<Initiative, 'publicDescription'>
  internalDescription: EditorField<Initiative, 'internalDescription'>
  requiresInvitationToJoin: EditorField<Initiative, 'requiresInvitationToJoin'>
  isAcceptingNewMembers: EditorField<Initiative, 'isAcceptingNewMembers'>
  isAcceptingNewPortfolios: EditorField<Initiative, 'isAcceptingNewPortfolios'>
  language: EditorField<Initiative, 'language'>
  pactaVersion: EditorField<Initiative, 'pactaVersion'>
  createdAt: EditorField<Initiative, 'createdAt'>
}

const createEditorInitiative = (initiative: Initiative): EditorInitiative => {
  return {
    id: {
      name: 'id',
      label: 'ID',
      isRequired: true,
      validation: Validation.AlphanumericAndDashesAndUnderscores,
      originalValue: initiative.id,
      currentValue: initiative.id,
    },
    name: {
      name: 'name',
      label: 'Name',
      isRequired: true,
      validation: Validation.NotEmpty,
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
      isRequired: true,
      validation: Validation.NotEmpty,
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
      isRequired: true,
      validation: Validation.NotEmpty,
      originalValue: initiative.language,
      currentValue: initiative.language,
    },
    pactaVersion: {
      name: 'pactaVersion',
      label: 'PACTA Version',
      isRequired: true,
      validation: Validation.NotEmpty,
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

const computedInitiative = (ref: Ref<EditorInitiative>): ComputedRef<Initiative> => {
  const v = ref.value
  return computed(() => ({
    ...asValue(v.id),
    ...asValue(v.name),
    ...asValue(v.affiliation),
    ...asValue(v.publicDescription),
    ...asValue(v.internalDescription),
    ...asValue(v.requiresInvitationToJoin),
    ...asValue(v.isAcceptingNewMembers),
    ...asValue(v.isAcceptingNewPortfolios),
    ...asValue(v.language),
    ...asValue(v.pactaVersion),
    ...asValue(v.createdAt),
  }))
}

const getComputedInitiativeChanges = (ref: Ref<EditorInitiative>): ComputedRef<InitiativeChanges> => {
  return computed(() => {
    const v = ref.value
    return {
      ...asChange(v.id),
      ...asChange(v.name),
      ...asChange(v.affiliation),
      ...asChange(v.publicDescription),
      ...asChange(v.internalDescription),
      ...asChange(v.requiresInvitationToJoin),
      ...asChange(v.isAcceptingNewMembers),
      ...asChange(v.isAcceptingNewPortfolios),
      ...asChange(v.language),
      ...asChange(v.pactaVersion),
      ...asChange(v.createdAt),
    }
  })
}

const getComputedIncompleteFields = (ref: Ref<EditorInitiative>): ComputedRef<string[]> => {
  return computed(() => {
    const v = ref.value
    return [
      ...asIncompleteField(v.id),
      ...asIncompleteField(v.name),
      ...asIncompleteField(v.affiliation),
      ...asIncompleteField(v.publicDescription),
      ...asIncompleteField(v.internalDescription),
      ...asIncompleteField(v.requiresInvitationToJoin),
      ...asIncompleteField(v.isAcceptingNewMembers),
      ...asIncompleteField(v.isAcceptingNewPortfolios),
      ...asIncompleteField(v.language),
      ...asIncompleteField(v.pactaVersion),
      ...asIncompleteField(v.createdAt),
    ]
  })
}

export const initiativeEditor = (i: Initiative) => {
  const ei = useState<EditorInitiative>('lib/editor/initiative')
  ei.value = createEditorInitiative(i)
  const incompleteFields = getComputedIncompleteFields(ei)
  const changes = getComputedInitiativeChanges(ei)
  const initiative = computedInitiative(ei)
  const setInitiative = (pv: Initiative) => { ei.value = createEditorInitiative(pv) }
  const hasChanges = computed(() => changes.value && Object.keys(changes.value).length > 0)
  const isIncomplete = computed(() => incompleteFields.value.length > 0)

  return {
    editorInitiative: ei,
    setInitiative,
    incompleteFields,
    changes,
    hasChanges,
    isIncomplete,
    initiative,
  }
}
