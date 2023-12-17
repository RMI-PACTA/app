import { type User } from '@/openapi/generated/pacta'
import { type EditorFieldsFor, type EditorValuesFor, Validation, type EditorComputedValues } from './common'
import { getEditorComputedValues, type Translation } from './utils'

export type EditorUserFields = EditorFieldsFor<User>
export type EditorUserValues = EditorValuesFor<User>

const createEditorUserFields = (translation: Translation): EditorUserFields => {
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
      helpText: tt('AdminHelpText'),
    },
    superAdmin: {
      name: 'superAdmin',
      label: tt('Super Admin'),
      helpText: tt('SuperAdminHelpText'),
    },
    name: {
      validation: [Validation.NotEmpty],
      name: 'name',
      label: tt('Name'),
      helpText: tt('NameHelpText'),
    },
    preferredLanguage: {
      name: 'preferredLanguage',
      label: tt('Preferred Language'),
      helpText: tt('PreferredLanguageHelpText'),
    },
  }
}

export const userEditor = (u: User, translation: Translation): EditorComputedValues<User> => {
  return getEditorComputedValues('lib/editor/user', u, createEditorUserFields, translation)
}
