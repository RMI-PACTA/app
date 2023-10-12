export interface EditorField<R, Name extends keyof R> {
  name: Name
  label: string
  validation?: Validation[]
  originalValue: R[Name]
  currentValue: R[Name]
}

export type EditorFieldsFor<R> = {
  [K in keyof R]-?: EditorField<R, K>
}

export const isValid = (editorField: EditorField<any, any>): boolean => {
  if (!editorField.validation === undefined) {
    return true
  }
  for (const v of (editorField.validation ?? [])) {
    if (!isValidFor(editorField, v)) {
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
const isValidFor = (editorField: EditorField<any, any>, validation: Validation): boolean => {
  switch (validation) {
    // TODO(gady PICKUP HERE)
    case Validation.NotEmpty:
      return !!editorField.currentValue
    case Validation.AlphanumericAndDashesAndUnderscores:
      return alphanumericAndDashesAndUnderscores.test(editorField.currentValue)
  }
}

export interface EditorComputedValues <R> {
  setEditorValue: (r: R) => void
  editorObject: Ref<EditorFieldsFor<R>>
  invalidFields: ComputedRef<string[]>
  changes: ComputedRef<Partial<R>>
  currentValue: ComputedRef<R>
  hasChanges: ComputedRef<boolean>
  isInvalid: ComputedRef<boolean>
  saveTooltip: ComputedRef<string | undefined>
  canSave: ComputedRef<boolean>
}
