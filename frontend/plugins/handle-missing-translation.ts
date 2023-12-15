import { presentOrFileBug } from '@/globalimports/present'

export default defineNuxtPlugin((nuxtApp) => {
  const prefix = 'plugins/missing-translations'
  const missingTranslations = useState<Map<string, Set<string>>>(`${prefix}/missingTranslations`, () => new Map())
  const numberMissing = useState<number>(`${prefix}/numberMissing`, () => 0)
  const handleMissingTranslation = (locale: string, key: string) => {
    if (!missingTranslations.value.has(locale)) {
      missingTranslations.value.set(locale, new Set())
    }
    const size = presentOrFileBug(missingTranslations.value.get(locale)).size
    presentOrFileBug(missingTranslations.value.get(locale)).add(key)
    if (presentOrFileBug(missingTranslations.value.get(locale)).size > size) {
      numberMissing.value++
      console.log(`missing ${locale} translation for "${key}" (${numberMissing.value} total)`)
    }
  }

  nuxtApp.vueApp.provide('handleMissingTranslation', handleMissingTranslation)
  const values = computed(() => {
    return missingTranslations.value
  })
  return {
    provide: {
      missingTranslations: {
        values,
        numberMissing,
      },
    },
  }
})
