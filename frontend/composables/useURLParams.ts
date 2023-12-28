import type { RouteParams, LocationQuery } from 'vue-router'
import { useRoute, stringifyQuery } from 'vue-router'
import { computed, type WritableComputedRef } from 'vue'

export const useURLParams = () => {
  const route = useRoute()
  const router = useRouter()

  const getVal = (src: RouteParams | LocationQuery, key: string): string | undefined => {
    const val = src[key]
    if (!val) {
      return undefined
    }

    if (Array.isArray(val)) {
      if (val.length === 0) {
        return undefined
      }
      if (!val[0]) {
        return undefined
      }
      return val[0]
    }

    return val
  }

  const setVal = (key: string, val: string | undefined) => {
    const query = new URLSearchParams(stringifyQuery(router.currentRoute.value.query))
    if (val) {
      query.set(key, val)
    } else {
      query.delete(key)
    }
    let qs = query.toString()
    if (qs) {
      qs = '?' + qs
    }
    void router.replace(qs)
  }

  const fromQueryReactive = (key: string): WritableComputedRef<string | undefined> => {
    return computed({
      get: () => getVal(router.currentRoute.value.query, key),
      set: (val: string | undefined) => { setVal(key, val) },
    })
  }

  const fromQueryReactiveWithDefault = (key: string, def: string): WritableComputedRef<string> => {
    const fqr = fromQueryReactive(key)
    return computed({
      get: () => fqr.value ?? def,
      set: (val: string) => {
        if (val === def) {
          fqr.value = undefined
        } else {
          fqr.value = val
        }
      },
    })
  }

  return {
    fromQuery: (key: string): string | undefined => {
      return getVal(route.query, key)
    },
    fromQueryReactive,
    fromQueryReactiveWithDefault,
    fromParams: (key: string): string | undefined => {
      return getVal(route.params, key)
    },
  }
}
