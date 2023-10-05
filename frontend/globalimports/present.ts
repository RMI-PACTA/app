import { createErrorWithRemediation, Remediation } from '@/lib/error'

export const present = <T>(t: T | undefined | null, r: Remediation, cause?: string): T => {
  if (t === undefined) {
    throw createErrorWithRemediation({
      name: 'present error',
      message: 'expected to be present but was undefined',
      cause,
    }, r)
  }
  if (t === null) {
    throw createErrorWithRemediation({
      name: 'present error',
      message: 'expected to be present but was null',
      cause,
    }, r)
  }
  return t
}

export const presentOrSuggestReload = <T>(t: T | undefined | null, cause?: string): T => present(t, Remediation.Reload, cause)
export const presentOrFileBug = <T>(t: T | undefined | null, cause?: string): T => present(t, Remediation.FileBug, cause)
export const presentOrCheckURL = <T>(t: T | undefined | null, cause?: string): T => present(t, Remediation.CheckUrl, cause)
