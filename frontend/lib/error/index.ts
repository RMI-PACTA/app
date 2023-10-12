import { type NuxtError } from 'nuxt/app'

export enum Remediation {
  None = 'none',
  Reload = 'reload',
  FileBug = 'file-bug',
  CheckUrl = 'check-url',
}

interface RemediationData {
  type: 'remediation'
  remediation: Remediation
}

type PACTAErrorData = RemediationData

interface PACTAError extends NuxtError {
  data: PACTAErrorData
}

export const createErrorWithRemediation = (err: string | Partial<NuxtError>, r: Remediation): PACTAError => {
  console.log(err)
  const nuxtErr = createError(err)
  if (!nuxtErr.data || typeof (nuxtErr.data) !== 'object') {
    nuxtErr.data = {}
  }
  nuxtErr.data.type = 'remediation'
  nuxtErr.data.remediation = r

  // TypeScript doesn't automatically pick up that `data` is always set, so we
  // give it a nudge.
  return nuxtErr as PACTAError
}
