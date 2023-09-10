import { ErrorWithRemediation, Remediation } from '@/lib/error'

export const present = <T>(t: T | undefined | null, r: Remediation): T => {
  if (t === undefined) {
    throw new ErrorWithRemediation(`expected to be present but was undefined: ${typeof t}.`, r)
  }
  if (t === null) {
    throw new ErrorWithRemediation(`expected to be present but was null: ${typeof t}.`, r)
  }
  return t
}

export const presentOrSuggestReload = <T>(t: T | undefined | null): T => present(t, Remediation.Reload)
export const presentOrFileBug = <T>(t: T | undefined | null): T => present(t, Remediation.FileBug)
export const presentOrCheckURL = <T>(t: T | undefined | null): T => present(t, Remediation.CheckUrl)
