import { type User } from '@/openapi/generated/pacta'
import { type EditorFieldsFor, Validation, type EditorComputedValues } from './common'
import { getEditorComputedValues } from './utils'

export type EditorUser = EditorFieldsFor<User>

// TODO(grady) Add i18n for forms.
const createEditorUser = (user: User): EditorUser => {
  return {
    id: {
      name: 'id',
      label: 'ID',
      originalValue: user.id,
      currentValue: user.id,
    },
    enteredEmail: {
      name: 'enteredEmail',
      label: 'Entered Email',
      originalValue: user.enteredEmail,
      currentValue: user.enteredEmail,
    },
    canonicalEmail: {
      name: 'canonicalEmail',
      label: 'Canonical Email',
      originalValue: user.canonicalEmail,
      currentValue: user.canonicalEmail,
    },
    admin: {
      name: 'admin',
      label: 'Admin',
      originalValue: user.admin,
      currentValue: user.admin,
    },
    superAdmin: {
      name: 'superAdmin',
      label: 'Super Admin',
      originalValue: user.superAdmin,
      currentValue: user.superAdmin,
    },
    name: {
      name: 'name',
      label: 'Name',
      validation: [Validation.NotEmpty],
      originalValue: user.name,
      currentValue: user.name,
    },
    preferredLanguage: {
      name: 'preferredLanguage',
      label: 'Preferred Language',
      originalValue: user.preferredLanguage,
      currentValue: user.preferredLanguage,
    },
  }
}

export const userEditor = (u: User): EditorComputedValues<User> => {
  return getEditorComputedValues('lib/editor/user', u, createEditorUser)
}
