import { type Ref } from 'vue'
import { type ErrorObject } from 'serialize-error'
import { serializeError } from 'serialize-error'

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
  const error = useState<ErrorObject>('errorModal.error')
  const handleError = (err: Error) => {
    if (process.client) {
      console.log(err)
    }
    error.value = serializeError(err)
    errorModalVisible.value = true
    clearLoading()
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
  const withLoading = async <T> (fn: () => Promise<T>, opKey: string): (Promise<T>) => {
    startLoading(opKey)
    const p = fn()
    void p.finally(stopLoading(opKey))
    return await p
  }
  const onMountedWithLoading = (fn: () => void, opKey: string) => {
    startLoading(opKey)
    onMounted(() => {
      fn()
      stopLoading(opKey)()
    })
  }

  // fakeUsers
  const fakeUsersVisible = newModalVisibilityState('fakeUsersVisibile')

  // missingTranslations
  const missingTranslationsVisible = newModalVisibilityState('missingTranslationsVisibile')

  // permissionDenied
  const permissionDeniedVisible = newModalVisibilityState('permissionDeniedVisibile')
  const permissionDeniedError = useState<Error | null>(`${prefix}.permissionDeniedError`, () => null)
  const setPermissionDenied = (e: Error) => {
    permissionDeniedError.value = e
    permissionDeniedVisible.value = true
  }

  // newPortfolioGroup
  const newPortfolioGroupVisible = newModalVisibilityState('newPortfolioGroupVisibile')

  const anyBlockingModalOpen = computed(() => anyModalVisible.value || loading.value)

  return {
    anyBlockingModalOpen,
    newModalVisibilityState,
    fakeUsers: {
      fakeUsersVisible,
    },
    loading: {
      withLoading,
      onMountedWithLoading,
      startLoading,
      stopLoading,
      clearLoading,
      loading,
      loadingSet,
    },
    error: {
      error,
      errorModalVisible,
      handleError,
    },
    permissionDenied: {
      permissionDeniedVisible,
      permissionDeniedError,
      setPermissionDenied,
    },
    missingTranslations: {
      missingTranslationsVisible,
    },
    newPortfolioGroup: {
      newPortfolioGroupVisible,
    },
  }
}
