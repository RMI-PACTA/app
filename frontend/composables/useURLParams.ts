import type { RouteParams, LocationQuery } from 'vue-router'
import { useRoute, stringifyQuery } from 'vue-router'
import { computed, type WritableComputedRef } from 'vue'

export const useURLParams = () => {
  const route = useRoute()
  const router = useRouter()

  // A common bug prior to the pending-values approach was that within a single tick, we'd do multiple updates like
  // (1) setVal('foo', 1)
  // (2) setVal('bar', 2)
  // **router.currentRoute** is a reactive object that **does not get updated between ticks**.
  // Thus setVal('foo', 1) would set the value in the URL, but then setVal('bar', 2) would overwrite it.
  // Said another way, `setVal('foo', 1)` immediately followed by `getVal('foo'), would return the original value.
  // The PendingValues approach solves this by tracking the values that are pending to be set in the URL,
  // **and using them as the source of truth when set for `getVal` queries**.
  const pendingValues = useState<Map<string, string | undefined>>('useURLParams.pendingValues', () => new Map<string, string | undefined>())
  const hasPendingValues = computed(() => pendingValues.value.size > 0)

  // An interaction pattern we often have is
  //   (a) change something where the source of truth state is in the URL, then
  //   (b) reload the core data of the page
  // This is a helper method to enable a setter of one or more URL params to wait for them
  // to be reflected in the URL Query before initiating
  const waitForURLToUpdate = () => {
    if (!hasPendingValues.value) {
      return Promise.resolve()
    }
    return new Promise<void>((resolve) => {
      // Note, this just lets us know that the value set has been emptied into the URL.
      // It doesn't mean that the URL has been updated (i.e. reading currentRoute would yield the incorrect result).
      // To sidestep that we watch the currentRoute and wait for it to change AFTER the pending value set is emptied.
      const unwatchPV = watch(hasPendingValues, (hpv) => {
        if (!hpv) {
          unwatchPV()
          const unwatchRoute = watch(router.currentRoute, () => {
            unwatchRoute()
            resolve()
          })
        }
      })
    })
  }

  const getVal = (src: RouteParams | LocationQuery, key: string): string | undefined => {
    // NOTE: if the value is in pendingValues, that means it's the source of truth
    const pvs = pendingValues.value
    if (pvs.has(key)) {
      return pvs.get(key)
    }

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

  const resolvePendingValues = () => {
    const pvs = pendingValues.value
    const query = new URLSearchParams(stringifyQuery(router.currentRoute.value.query))
    for (const [key, val] of pvs) {
      if (val === undefined) {
        query.delete(key)
      } else {
        query.set(key, val)
      }
    }
    let qs = query.toString()
    if (qs) {
      qs = '?' + qs
    }
    void router.replace(qs).then(() => nextTick(() => {
      pendingValues.value = new Map<string, string | undefined>()
    }))
  }

  const setVal = (key: string, val: string | undefined) => {
    pendingValues.value.set(key, val)
    // Note: we only try to resolve pending values upon next tick.
    // This isn't required, but it prevents us from doing more work on the browser,
    // since in multi-update cases, we'll generate multiple replace() calls which will
    // become redundant upon the next tick updating the URLs.
    void nextTick(resolvePendingValues)
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
    waitForURLToUpdate,
  }
}
