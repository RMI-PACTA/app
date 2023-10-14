export type LanguageCode = 'en' | 'de' | 'fr' | 'es'

// TODO(grady) Add support for icon-based flags
export interface LanguageOption { label: string, code: LanguageCode }

export const LanguageOptions: LanguageOption[] = [
  { label: 'English', code: 'en' },
  { label: 'Deutsch', code: 'de' },
  { label: 'Français', code: 'fr' },
  { label: 'Español', code: 'es' },
]
