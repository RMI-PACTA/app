import { type EditorFieldsFor, type EditorValuesFor, isValid, type EditorComputedValues, createEditorValues } from './common'

const getCurrentValue = <R> (values: EditorValuesFor<R>): R => {
  const result: Partial<R> = {}
  for (const key in values) {
    result[key] = values[key].currentValue
  }
  return result as R
}

const getInvalidFields = <R> (fields: EditorFieldsFor<R>, values: EditorValuesFor<R>): string[] => {
  const result: string[] = []
  for (const key in values) {
    const v = values[key]
    const f = fields[key]
    if (!isValid(f, v)) {
      result.push(f.label)
    }
  }
  return result
}

const getChanges = <R> (values: EditorValuesFor<R>): Partial<R> => {
  const result: Partial<R> = {}
  for (const key in values) {
    if (values[key].currentValue !== values[key].originalValue) {
      result[key] = values[key].currentValue
    }
  }
  return result
}

export type Translation = ReturnType<typeof useI18n>
type ToEFF <R> = (translation: Translation) => EditorFieldsFor<R>

export const getEditorComputedValues = <R> (key: string, r: R, toEFF: ToEFF<R>, translation: Translation): EditorComputedValues<R> => {
  const editorValues = useState<EditorValuesFor<R>>(key)
  editorValues.value = createEditorValues(r)
  const editorFields = computed(() => toEFF(translation))
  const invalidFields = computed(() => getInvalidFields(editorFields.value, editorValues.value))
  const changes = computed(() => getChanges(editorValues.value))
  const currentValue = computed(() => getCurrentValue(editorValues.value))
  const setEditorValue = (r: R) => { editorValues.value = createEditorValues(r) }
  const hasChanges = computed(() => changes.value && Object.keys(changes.value).length > 0)
  const isInvalid = computed(() => invalidFields.value.length > 0)
  const canSave = computed(() => hasChanges.value && !isInvalid.value)
  const { t } = translation
  const allSaved = t('lib/editor/utils/AllChangesSaved')
  const cannotSave = t('lib/editor/utils/CannotSaveWithInvalidFields')
  const saveTooltip = computed<string | undefined>(() => {
    if (!hasChanges.value) { return allSaved }
    if (isInvalid.value) { return `${cannotSave}: ${invalidFields.value.join(', ')}` }
    return undefined
  })

  return {
    setEditorValue,
    editorValues,
    editorFields,
    invalidFields,
    changes,
    currentValue,
    hasChanges,
    isInvalid,
    canSave,
    saveTooltip,
  }
}
