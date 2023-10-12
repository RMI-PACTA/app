import { type WritableComputedRef } from 'vue'

export interface ForUse<T> {
  isRequired: boolean
  isCompleted: ComputedRef<boolean>
  value: WritableComputedRef<T>
}

export enum Validation {
  NotEmpty = 'NotEmpty',
  AlphanumericAndDashesAndUnderscores = 'AlphanumericAndDashesAndUnderscores',
}

export interface EditorField<R, Name extends keyof R> {
  name: Name
  label: string
  isRequired?: boolean
  validation?: Validation | undefined
  originalValue: R[Name]
  currentValue: R[Name]
}

export const asChange = <R, N extends keyof R>(field: EditorField<R, N>): Partial<R> => {
  const result: Partial<R> = {}
  if (field.originalValue === field.currentValue) {
    return result
  }
  result[field.name] = field.currentValue
  return result
}

export const asValue = <R, N extends keyof R>(field: EditorField<R, N>): Pick<R, N> => {
  return { [field.name]: field.currentValue } as unknown as Pick<R, N>
}

export const asIncompleteField = <R, N extends keyof R>(field: EditorField<R, N>): string[] => {
  if (field.isRequired && !isComplete(field)) {
    return [field.label]
  }
  return []
}

export const asComputed = <R, N extends keyof R>(field: WritableComputedRef<EditorField<R, N>>): WritableComputedRef<R[N]> => {
  return computed<R[N]>({
    get: () => field.value.currentValue,
    set: (v: R[N]) => { field.value.currentValue = v },
  })
}

const alphanumericAndDashesAndUnderscores = /^[a-zA-Z0-9-_]+$/
export const isComplete = (editorField: EditorField<any, any>): boolean => {
  if (!editorField.isRequired || editorField.validation === undefined) {
    return true
  }
  switch (editorField.validation) {
    case Validation.NotEmpty:
      return !!editorField.currentValue
    case Validation.AlphanumericAndDashesAndUnderscores:
      return alphanumericAndDashesAndUnderscores.test(editorField.currentValue)
  }
}
