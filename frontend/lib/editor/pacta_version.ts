import { type PactaVersion, type PactaVersionChanges } from '@/openapi/generated/pacta'
import { type EditorField, asChange, asIncompleteField, asValue, Validation } from './common'
import { type Ref } from 'vue'

export interface EditorPactaVersion {
  id: EditorField<PactaVersion, 'id'>
  name: EditorField<PactaVersion, 'name'>
  description: EditorField<PactaVersion, 'description'>
  digest: EditorField<PactaVersion, 'digest'>
  isDefault: EditorField<PactaVersion, 'isDefault'>
  createdAt: EditorField<PactaVersion, 'createdAt'>
}

const createEditorPactaVersion = (pactaVersion: PactaVersion): EditorPactaVersion => {
  return {
    id: {
      name: 'id',
      label: 'ID',
      originalValue: pactaVersion.id,
      currentValue: pactaVersion.id,
    },
    name: {
      name: 'name',
      label: 'Name',
      isRequired: true,
      validation: Validation.NotEmpty,
      originalValue: pactaVersion.name,
      currentValue: pactaVersion.name,
    },
    description: {
      name: 'description',
      label: 'Description',
      isRequired: true,
      validation: Validation.NotEmpty,
      originalValue: pactaVersion.description,
      currentValue: pactaVersion.description,
    },
    digest: {
      name: 'digest',
      label: 'Docker Image Digest',
      isRequired: true,
      validation: Validation.NotEmpty,
      originalValue: pactaVersion.digest,
      currentValue: pactaVersion.digest,
    },
    isDefault: {
      name: 'isDefault',
      label: 'Is Default Version',
      originalValue: pactaVersion.isDefault,
      currentValue: pactaVersion.isDefault,
    },
    createdAt: {
      name: 'createdAt',
      label: 'Created At',
      originalValue: pactaVersion.createdAt,
      currentValue: pactaVersion.createdAt,
    },
  }
}

const getComputedPactaVersionChanges = (ref: Ref<EditorPactaVersion>): ComputedRef<PactaVersionChanges> => {
  return computed(() => {
    const epv = ref.value
    return {
      ...asChange(epv.id),
      ...asChange(epv.name),
      ...asChange(epv.description),
      ...asChange(epv.digest),
      ...asChange(epv.isDefault),
      ...asChange(epv.createdAt),
    }
  })
}

const getComputedPactaVersion = (ref: Ref<EditorPactaVersion>): ComputedRef<PactaVersion> => {
  return computed(() => {
    const epv = ref.value
    return {
      ...asValue(epv.id),
      ...asValue(epv.name),
      ...asValue(epv.description),
      ...asValue(epv.digest),
      ...asValue(epv.isDefault),
      ...asValue(epv.createdAt),
    }
  })
}

const getComputedIncompleteFields = (ref: Ref<EditorPactaVersion>): ComputedRef<string[]> => {
  return computed(() => {
    const epv = ref.value
    return [
      ...asIncompleteField(epv.id),
      ...asIncompleteField(epv.name),
      ...asIncompleteField(epv.description),
      ...asIncompleteField(epv.digest),
      ...asIncompleteField(epv.isDefault),
      ...asIncompleteField(epv.createdAt),
    ]
  })
}

export const pactaVersionEditor = (pv: PactaVersion) => {
  const epv = useState<EditorPactaVersion>('lib/editor/pacta-version')
  epv.value = createEditorPactaVersion(pv)
  const incompleteFields = getComputedIncompleteFields(epv)
  const changes = getComputedPactaVersionChanges(epv)
  const pactaVersion = getComputedPactaVersion(epv)
  // const pactaVersion = getComputedPactaVersion(epv)
  const setPactaVersion = (pv: PactaVersion) => { epv.value = createEditorPactaVersion(pv) }
  const hasChanges = computed(() => changes.value && Object.keys(changes.value).length > 0)
  const isIncomplete = computed(() => incompleteFields.value.length > 0)
  return {
    editorPactaVersion: epv,
    setPactaVersion,
    incompleteFields,
    changes,
    hasChanges,
    isIncomplete,
    pactaVersion,
  }
}
