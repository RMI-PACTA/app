interface SimpleAsyncDataReturn<T> {
  data: Ref<T>
  refresh: () => Promise<void>
}

export const useSimpleAsyncData = async <T>(key: string, fn: () => Promise<T>): Promise<SimpleAsyncDataReturn<T>> => {
  // const { loading: { withLoading } } = useModal()
  // withLoading(fn, key)
  const { data: dataRef, refresh, error } = await useAsyncData(key, fn)
  if (error.value) {
    console.log('error', error.value)
    throw createError(error.value)
  }
  const data = dataRef as unknown as Ref<T>
  return { data, refresh }
}
