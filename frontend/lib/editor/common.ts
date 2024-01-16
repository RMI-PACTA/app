export interface EditorField<R, Name extends keyof R> {
  name: Name
  label: string
  validation?: Validation[]
  startHelpTextExpanded?: boolean
  helpText?: string
  loadingLabel?: string
  invalidLabel?: string
  validLabel?: string
}

export interface EditorValue<R, Name extends keyof R> {
  originalValue: R[Name]
  currentValue: R[Name]
}

export type EditorFieldsFor<R> = {
  [K in keyof R]-?: EditorField<R, K>
}

export type EditorValuesFor<R> = {
  [K in keyof R]-?: EditorValue<R, K>
}

export const isValid = <R, K extends keyof R>(editorField: EditorField<R, K>, editorValue: EditorValue<R, K>): boolean => {
  if (!editorField.validation === undefined) {
    return true
  }
  for (const v of (editorField.validation ?? [])) {
    if (!isValidFor(editorValue, v)) {
      return false
    }
  }
  return true
}

export enum Validation {
  NotEmpty = 'NotEmpty',
  AlphanumericAndDashesAndUnderscores = 'AlphanumericAndDashesAndUnderscores',
}

const alphanumericAndDashesAndUnderscores = /^[a-zA-Z0-9-_]+$/
const isValidFor = (editorValue: EditorValue<any, any>, validation: Validation): boolean => {
  switch (validation) {
    case Validation.NotEmpty:
      return !!editorValue.currentValue
    case Validation.AlphanumericAndDashesAndUnderscores:
      return alphanumericAndDashesAndUnderscores.test(editorValue.currentValue)
  }
}

export interface EditorComputedValues <R> {
  setEditorValue: (r: R) => void
  editorValues: Ref<EditorValuesFor<R>>
  editorFields: ComputedRef<EditorFieldsFor<R>>
  invalidFields: ComputedRef<string[]>
  changes: ComputedRef<Partial<R>>
  currentValue: ComputedRef<R>
  hasChanges: ComputedRef<boolean>
  isInvalid: ComputedRef<boolean>
  saveTooltip: ComputedRef<string | undefined>
  canSave: ComputedRef<boolean>
  resetEditor: () => void
}

export const createEditorValues = <R>(r: R): EditorValuesFor<R> => Object.keys(r as any)
  .map(k => k as keyof R)
  .map((key: keyof R) => ({
    [key]: {
      originalValue: r[key],
      currentValue: r[key],
    },
  }))
  .reduce((acc, curr) => ({ ...acc, ...curr }), {}) as EditorValuesFor<R>
