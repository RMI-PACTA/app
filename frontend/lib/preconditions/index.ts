import type { Maybe } from 'graphql/jsutils/Maybe'

export const present = <T>(t: T | undefined | null | Maybe<T>): T => {
  if (t === undefined) {
    throw new Error(`Expected to be present but was undefined: ${typeof t}.`)
  }
  if (t === null) {
    throw new Error(`Expected to be present but was null: ${typeof t}.`)
  }
  return t
}
