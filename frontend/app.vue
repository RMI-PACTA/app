<script setup lang="ts">
const { loading: { clearLoading }, error: { errorModalVisible, error } } = useModal()

const handleError = (err: Error) => {
  error.value = err
  errorModalVisible.value = true
  clearLoading()
}

onErrorCaptured((err: unknown, _instance: ComponentPublicInstance | null, _info: string) => {
  let error: Error | undefined
  if (err instanceof Error) {
    error = err
  } else if (typeof (err) === 'string') {
    error = new Error(err)
  } else {
    error = new Error('unknown error', { cause: err })
  }
  handleError(error)
  return false // Don't propagate
})
</script>

<template>
  <NuxtLayout />
</template>
