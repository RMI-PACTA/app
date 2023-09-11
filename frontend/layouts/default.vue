<script setup lang="ts">
import { onMounted } from 'vue'

const { loading: { onMountedWithLoading, loadingSet }, anyBlockingModalOpen, error: { setError } } = useModal()

const handleError = (event: Event & { reason: Error }) => {
  event.preventDefault()
  const { reason } = event
  setError('fallback')(reason)
  loadingSet.value.clear()
}

onMountedWithLoading(() => { /* nothing to do */ }, 'defaultLayout.onMountedWithLoading')
onMounted(() => {
  window.addEventListener('unhandledrejection', handleError)
})

</script>

<template>
  <div class="app-default-layout">
    <StandardNav />
    <div
      class="flex flex-column align-items-center relative"
      :aria-hidden="anyBlockingModalOpen"
    >
      <main
        class="px-3 md:px-6 w-full lg:w-10 xl:w-8 mx-auto"
        style="min-height: calc(100vh - 9rem - 4px);"
      >
        <NuxtErrorBoundary>
          <template #error="{ error, clearError }">
            {{ setError(error) }}
            {{ clearError() }}
          </template>
          <NuxtPage />
        </NuxtErrorBoundary>
      </main>
    </div>
    <ModalGroup />
    <StandardFooter />
  </div>
</template>
