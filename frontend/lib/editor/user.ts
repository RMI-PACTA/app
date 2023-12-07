import { type User } from '@/openapi/generated/pacta'
import { type EditorFieldsFor, type EditorValuesFor, Validation, type EditorComputedValues } from './common'
import { getEditorComputedValues, type Translation } from './utils'

export type EditorUserFields = EditorFieldsFor<User>
export type EditorUserValues = EditorValuesFor<User>

const createEditorUserValues = (user: User): EditorUserValues => {
  return {
    id: {
      originalValue: user.id,
      currentValue: user.id,
    },
    enteredEmail: {
      originalValue: user.enteredEmail,
      currentValue: user.enteredEmail,
    },
    canonicalEmail: {
      originalValue: user.canonicalEmail,
      currentValue: user.canonicalEmail,
    },
    admin: {
      originalValue: user.admin,
      currentValue: user.admin,
    },
    superAdmin: {
      originalValue: user.superAdmin,
      currentValue: user.superAdmin,
    },
    name: {
      originalValue: user.name,
      currentValue: user.name,
    },
    preferredLanguage: {
      originalValue: user.preferredLanguage,
      currentValue: user.preferredLanguage,
    },
  }
}

const createEditorUserFields = (user: User, translation: Translation): EditorUserFields => {
  const tt = (key: string) => translation.t(`lib/editor/user/${key}`)
  return {
    id: {
      name: 'id',
      label: tt('ID'),
    },
    enteredEmail: {
      name: 'enteredEmail',
      label: tt('Entered Email'),
    },
    canonicalEmail: {
      name: 'canonicalEmail',
      label: tt('Canonical Email'),
    },
    admin: {
      name: 'admin',
      label: tt('Admin'),
    },
    superAdmin: {
      name: 'superAdmin',
      label: tt('Super Admin'),
    },
    name: {
      validation: [Validation.NotEmpty],
      name: 'name',
      label: tt('Name'),
    },
    preferredLanguage: {
      name: 'preferredLanguage',
      label: tt('Preferred Language'),
    },
  }
}

export const userEditor = (u: User, translation: Translation): EditorComputedValues<User> => {
  return getEditorComputedValues('lib/editor/user', u, createEditorUserFields, createEditorUserValues, translation)
}
