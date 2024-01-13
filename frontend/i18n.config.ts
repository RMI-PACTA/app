export default defineI18nConfig(() => ({
  locale: 'en',
  fallbackLocale: 'en',
  objectNotation: true,
  missing: (locale, key, vm) => {
    // TODO(grady) figure out how to skip this if we're in production + just log.
    // Consider using process.env.NODE_ENV == 'prod', etc.
    const fn = vm?.appContext.app.$nuxt.$missingTranslations.handleMissingTranslation
    if (fn) {
      fn(locale, key)
    } else {
      console.warn(`No handleMissingTranslation function found, can't handle ${locale} ${key}`)
    }
  },
}))
