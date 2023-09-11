import { type Error as OPAIError } from 'openapi/generated/pacta'
import { type Ref } from 'vue'

export const useModal = () => {
  const prefix = 'useModal'

  const nVisibilityStates = useState<number>(`${prefix}.nVisibilityStates`)
  const visibilityStates: Array<Ref<boolean>> = []
  const newModalVisibilityState = (suffix: string) => {
    const result = useState<boolean>(`${prefix}.${suffix}`, () => false)
    visibilityStates.push(result)
    nVisibilityStates.value++
    return result
  }
  const anyModalVisible = computed(() => {
    // This is a trick to trigger reactivity when the number of states changes
    // without trying to serialize/state-ify the list of visibilities!
    if (nVisibilityStates.value === 0) {
      return false
    }
    return visibilityStates.some((vs) => vs.value)
  })

  // error
  const errorModalVisible = newModalVisibilityState('errorModalVisible')
  const error = useState<Error | null>(`${prefix}.error`, () => null)
  const setError = (opKey: string) => {
    return (err?: Error) => {
      error.value = err ?? new Error('an unknown error occurred')
      errorModalVisible.value = true
      clearLoading()

      // We used to re-throw here, but that just breaks the page (e.g. no more
      // navigation), since it ends up propagating to the top-level. Since
      // setError is 'handling' the error, we don't re-throw.
    }
  }
  const setAndClearError = (err?: Error, fn: () => void) => {
    setError('setAndClearError')(err)
    fn()
  }
  const withErrorHandling = async (fn: () => Promise<unknown>, opKey: string): Promise<unknown> => {
    return await fn().catch(setError(opKey))
  }

  // loading
  const loadingSet = useState<Set<string>>(`${prefix}.loadingSet`, () => new Set<string>())
  const loading = computed(() => loadingSet.value.size > 0)
  const startLoading = (loadKey: string) => {
    loadingSet.value.add(loadKey)
  }
  const stopLoading = (loadKey: string) => {
    return () => loadingSet.value.delete(loadKey)
  }
  const clearLoading = () => { loadingSet.value.clear() }
  const withLoadingAndErrorHandling = async <T> (fn: () => Promise<T>, opKey: string): (Promise<T>) => {
    startLoading(opKey)
    const p = fn()
    p.catch(setError(opKey)).finally(stopLoading(opKey))
    return await p
  }
  const onMountedWithLoading = (fn: () => void, opKey: string) => {
    startLoading(opKey)
    onMounted(() => {
      fn()
      stopLoading(opKey)()
    })
  }

  // permissionDenied
  const permissionDeniedVisible = newModalVisibilityState('permissionDeniedVisibile')
  const permissionDeniedError = useState<Error | null>(`${prefix}.permissionDeniedError`, () => null)
  const setPermissionDenied = (e: Error) => {
    permissionDeniedError.value = e
    permissionDeniedVisible.value = true
  }

  const anyBlockingModalOpen = computed(() => anyModalVisible.value || loading.value)

  const handleOAPIError = async <T>(t: OPAIError | T): Promise<T> => {
    return await new Promise<T>((resolve, reject) => {
      // TODO(#10) Rephrase this once we use 300+ for all errors
      if (t instanceof Object && Object.prototype.hasOwnProperty.call(t, 'message')) {
        reject(new Error(JSON.stringify(t)))
      } else {
        resolve(t as T)
      }
    })
  }

  return {
    anyBlockingModalOpen,
    newModalVisibilityState,
    loading: {
      withLoadingAndErrorHandling,
      onMountedWithLoading,
      startLoading,
      stopLoading,
      clearLoading,
      loading,
      loadingSet
    },
    error: {
      setError,
      setAndClearError,
      error,
      withErrorHandling,
      errorModalVisible,
      withLoadingAndErrorHandling,
      handleOAPIError
    },
    permissionDenied: {
      permissionDeniedVisible,
      permissionDeniedError,
      setPermissionDenied
    }
  }
}
