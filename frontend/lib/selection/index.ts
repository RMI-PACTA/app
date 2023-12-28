export const selectedCountSuffix = <T>(ts: T[] | undefined | null) => {
  if (!ts || ts.length < 2) {
    return ''
  }
  return ` (${ts.length})`
}
