const defaultLocale = 'en-US'

const getLocale = (): string => {
  const { locale } = useI18n()
  if (locale) {
    return locale.value
  }
  return defaultLocale
}

export const humanReadableTimeFromStandardString = (asStr: string): string => {
  return humanReadableTime(new Date(Date.parse(asStr)))
}

export const humanReadableDateLongFromStandardString = (s: string): string => {
  return humanReadableDateLong(new Date(Date.parse(s)))
}

export const humanReadableDateLong = (t: Date): string => {
  const locale = getLocale()
  if (new Date().getDate() === t.getDate()) {
    return 'today'
  }
  return t.toLocaleString(locale, { dateStyle: 'long' })
}

export const humanReadableDateFromStandardString = (s: string): string => {
  return humanReadableDate(new Date(Date.parse(s)))
}

export const humanReadableDate = (t: Date): string => {
  const locale = getLocale()
  if (new Date().getDate() === t.getDate()) {
    return 'today'
  }
  return t.toLocaleString(locale, { dateStyle: 'short' })
}

export const humanReadableTime = (t: Date): string => {
  const locale = getLocale()
  if (new Date().getDate() === t.getDate()) {
    return t.toLocaleTimeString(locale, { timeStyle: 'short' })
  }
  return t.toLocaleString(locale, { dateStyle: 'short', timeStyle: 'short' })
}

export const humanReadableTimeWithSeconds = (t: Date): string => {
  const locale = getLocale()
  const now = new Date()
  if (now.getDate() === t.getDate()) {
    if (now.getHours() === t.getHours()) {
      return t.toLocaleTimeString(locale, { hour: 'numeric', minute: '2-digit', second: '2-digit' })
    }
    return t.toLocaleTimeString(locale, { timeStyle: 'short' })
  }
  return t.toLocaleString(locale, { dateStyle: 'short', timeStyle: 'short' })
}
