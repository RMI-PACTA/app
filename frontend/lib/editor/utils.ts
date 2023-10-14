import { type EditorFieldsFor, isValid, type EditorComputedValues } from './common'

const getCurrentValue = <R> (record: EditorFieldsFor<R>): R => {
  const result: Partial<R> = {}
  for (const key in record) {
    result[key] = record[key].currentValue
  }
  return result as R
}

const getInvalidFields = <R> (record: EditorFieldsFor<R>): string[] => {
  const result: string[] = []
  for (const key in record) {
    if (!isValid(record[key])) {
      result.push(key)
    }
  }
  return result
}

const getChanges = <R> (record: EditorFieldsFor<R>): Partial<R> => {
  const result: Partial<R> = {}
  for (const key in record) {
    if (record[key].currentValue !== record[key].originalValue) {
      result[key] = record[key].currentValue
    }
  }
  return result
}

type ToEFF <R> = (r: R) => EditorFieldsFor<R>

export const getEditorComputedValues = <R> (key: string, r: R, fn: ToEFF<R>): EditorComputedValues<R> => {
  const eff = useState<EditorFieldsFor<R>>(key)
  eff.value = fn(r)
  const invalidFields = computed(() => getInvalidFields(eff.value))
  const changes = computed(() => getChanges(eff.value))
  const currentValue = computed(() => getCurrentValue(eff.value))
  const setEditorValue = (r: R) => { eff.value = fn(r) }
  const hasChanges = computed(() => changes.value && Object.keys(changes.value).length > 0)
  const isInvalid = computed(() => invalidFields.value.length > 0)
  const canSave = computed(() => hasChanges.value && !isInvalid.value)
  const saveTooltip = computed<string | undefined>(() => {
    if (!hasChanges.value) { return 'All changes saved' }
    if (isInvalid.value) { return `Cannot save with invalid fields: ${invalidFields.value.join(', ')}` }
    return undefined
  })

  return {
    setEditorValue,
    editorObject: eff,
    invalidFields,
    changes,
    currentValue,
    hasChanges,
    isInvalid,
    canSave,
    saveTooltip,
  }
}
