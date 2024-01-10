import { createErrorWithRemediation, Remediation } from '@/lib/error'

export const present = <T>(t: T | undefined | null, r: Remediation, cause?: string): T => {
  if (t !== undefined && t !== null) {
    return t
  }
  const desc = t === undefined ? 'undefined' : 'null'
  let msg = `expected present but was ${desc}`
  if (cause) {
    msg += `: ${cause}`
  }
  const offendingLine = (new Error().stack ?? '').split('\n').find((line: string, i: number) => i > 0 && !line.includes('present.ts'))
  if (offendingLine) {
    msg += `@ ${offendingLine}`
  }
  throw createErrorWithRemediation(new Error(msg), r)
}

export const presentOrSuggestReload = <T>(t: T | undefined | null, cause?: string): T => present(t, Remediation.Reload, cause)
export const presentOrFileBug = <T>(t: T | undefined | null, cause?: string): T => present(t, Remediation.FileBug, cause)
export const presentOrCheckURL = <T>(t: T | undefined | null, cause?: string): T => present(t, Remediation.CheckUrl, cause)
