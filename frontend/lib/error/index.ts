export enum Remediation {
  None = 'none',
  Reload = 'reload',
  FileBug = 'file-bug',
  CheckUrl = 'check-url',
}

export class ErrorWithRemediation extends Error {
  constructor (msg: string, public readonly remediation: Remediation) {
    super(msg)
  }
}
