export default defineI18nConfig(() => ({
  locale: 'en',
  fallbackLocale: 'en',
  objectNotation: true,
  missing: (locale, key, vm) => {
    // TODO(grady) figure out how to skip this if we're in production + just log.
    // Consider using process.env.NODE_ENV == 'prod', etc.
    const fn = null // inject('handleMissingTranslation')
    if (fn) {
      const callable = fn as (locale: string, key: string) => void
      callable(locale, key)
    }
  },
}))
