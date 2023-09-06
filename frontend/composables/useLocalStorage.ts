import { computed, onMounted, type WritableComputedRef } from 'vue'
import { v4 as uuidv4 } from 'uuid'
import { present } from '@/lib/preconditions'

export const useLocalStorage = () => {
  // This code will be invoked on server, so we want a simple way of telling all of the computed
  // values produced by this composable to recompute themselves on the client once local storage
  // becomes available. We accomplish this with a single onMounted hook which forces recomputation.
  const isLocal = useState<boolean>('useLocalStorage.isLocal', () => process.client)
  onMounted(() => { isLocal.value = process.client })

  // This creates a local-storage backed state that doesn't cause hydration issues,
  // by using a default value for SSR, and then executing a onMounted hook to lookup
  // the value from local storage on the client.
  const computedStringLocalStorageValue = (key: string, defaultValue: string): WritableComputedRef<string> => {
    const state = useState<string>(`useLocalStorage.state[${key}]`, () => defaultValue)
    onMounted(() => {
      const stored = localStorage.getItem(key)
      if (stored === null) {
        state.value = defaultValue
      } else {
        state.value = stored
      }
    })
    return computed({
      get: () => state.value,
      set: (value: string) => {
        if (isLocal.value) {
          localStorage.setItem(key, value)
        }
        state.value = value
      }
    })
  }

  const computedBooleanLocalStorageValue = (key: string, defaultValue: boolean): WritableComputedRef<boolean> => {
    // We defer to a string state under the covers to allow for easier updates to the
    // mechanism used to achieve reactivity without hydration issues.
    const stringState = computedStringLocalStorageValue('bool:' + key, `${defaultValue}`)
    return computed({
      get: () => stringState.value === 'true',
      set: (value: boolean) => { stringState.value = `${value}` }
    })
  }

  const computedDateLocalStorageValue = (key: string, defaultValue: Date): WritableComputedRef<Date> => {
    const stringState = computedStringLocalStorageValue(`date:${key}`, defaultValue.toString())
    return computed({
      get: () => {
        const s = stringState.value
        return s ? new Date(s) : defaultValue
      },
      set: (d: Date) => { stringState.value = d.toString() }
    })
  }

  const computedStringSetLocalStorageValue = (key: string, defaultValue: Set<string>): WritableComputedRef<Set<string>> => {
    const stringState = computedStringLocalStorageValue(`strset:${key}`, JSON.stringify([...defaultValue]))
    return computed({
      get: () => {
        const s = stringState.value
        return s ? new Set<string>(JSON.parse(s) as string[]) : defaultValue
      },
      set: (s: Set<string>) => {
        stringState.value = JSON.stringify([...s])
      }
    })
  }

  const undefinedDeviceId = 'device-id-undefined'
  const deviceId = computedStringLocalStorageValue('device-Id', undefinedDeviceId)
  onMounted(() => {
    if (deviceId.value === undefinedDeviceId) {
      deviceId.value = uuidv4()
    }
  })

  const getDeviceId = (): string => {
    return present(deviceId.value)
  }

  const helpTextExpanded = (helpTextId: string) => computedBooleanLocalStorageValue(`helpTextExpanded-${helpTextId}`, true)

  return {
    computedBooleanLocalStorageValue,
    computedStringLocalStorageValue,
    computedStringSetLocalStorageValue,
    computedDateLocalStorageValue,
    helpTextExpanded,
    deviceId,
    getDeviceId
  }
}
