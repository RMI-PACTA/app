import { Language } from '@/openapi/generated/pacta'

export type LanguageCode = 'en' | 'de' | 'fr' | 'es'

// TODO(grady) Add support for icon-based flags
export interface LanguageOption { label: string, code: LanguageCode, language: Language }

export const LanguageOptions: LanguageOption[] = [
  { label: 'English', code: 'en', language: Language.LANGUAGE_EN },
  { label: 'Deutsch', code: 'de', language: Language.LANGUAGE_DE },
  { label: 'Français', code: 'fr', language: Language.LANGUAGE_FR },
  { label: 'Español', code: 'es', language: Language.LANGUAGE_ES },
]

export const languageToOption = (language: Language): LanguageOption => {
  return presentOrFileBug(
    LanguageOptions.find(option => option.language === language),
    `languageToOption not found: '${language}'`,
  )
}
