import { type NuxtError } from 'nuxt/app'

export enum Remediation {
  None = 'none',
  Reload = 'reload',
  FileBug = 'file-bug',
  CheckUrl = 'check-url',
  Silent = 'silent',
}

interface RemediationData {
  type: 'remediation'
  remediation: Remediation
}

type PACTAErrorData = RemediationData

interface PACTAError extends NuxtError {
  data: PACTAErrorData
}

export const isSilent = (err: any): boolean => {
  if (!('data' in err)) {
    return false
  }

  if (typeof (err.data) !== 'object' || err.data === null) {
    return false
  }

  if (!('type' in err.data) || typeof (err.data.type) !== 'string') {
    return false
  }

  if (err.data.type !== 'remediation') {
    return false
  }

  if (!('remediation' in err.data) || typeof (err.data.remediation) !== 'string') {
    return false
  }

  return err.data.remediation === Remediation.Silent
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
