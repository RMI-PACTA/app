export const useTime = () => {
  const { locale: i18nLocale, t } = useI18n()
  const defaultLocale = 'en-US'
  const locale = computed(() => i18nLocale.value || defaultLocale)

  const prefix = 'composables/useTime'
  const tt = (key: string) => t(`${prefix}.${key}`)

  const humanReadableDateLong = (t: Date): ComputedRef<string> => {
    return computed(() => {
      if (new Date().getDate() === t.getDate()) {
        return tt('today')
      }
      return t.toLocaleString(locale.value, { dateStyle: 'long' })
    })
  }

  const humanReadableDate = (t: Date): ComputedRef<string> => {
    return computed(() => {
      if (new Date().getDate() === t.getDate()) {
        return tt('today')
      }
      return t.toLocaleString(locale.value, { dateStyle: 'short' })
    })
  }

  const humanReadableTime = (t: Date): ComputedRef<string> => {
    return computed(() => {
      if (new Date().getDate() === t.getDate()) {
        return t.toLocaleTimeString(locale.value, { timeStyle: 'short' })
      }
      return t.toLocaleString(locale.value, { dateStyle: 'short', timeStyle: 'short' })
    })
  }

  const humanReadableDateLongFromStandardString = (s: string): ComputedRef<string> => humanReadableDateLong(new Date(Date.parse(s)))
  const humanReadableDateFromStandardString = (s: string): ComputedRef<string> => humanReadableDate(new Date(Date.parse(s)))
  const humanReadableTimeFromStandardString = (asStr: string): ComputedRef<string> => humanReadableTime(new Date(Date.parse(asStr)))

  const humanReadableTimeWithSeconds = (t: Date): ComputedRef<string> => {
    return computed(() => {
      const now = new Date()
      if (now.getDate() === t.getDate()) {
        if (now.getHours() === t.getHours()) {
          return t.toLocaleTimeString(locale.value, { hour: 'numeric', minute: '2-digit', second: '2-digit' })
        }
        return t.toLocaleTimeString(locale.value, { timeStyle: 'short' })
      }
      return t.toLocaleString(locale.value, { dateStyle: 'short', timeStyle: 'short' })
    })
  }

  return {
    humanReadableDateFromStandardString,
    humanReadableDateLongFromStandardString,
    humanReadableTimeFromStandardString,
    humanReadableTimeWithSeconds,
  }
}
