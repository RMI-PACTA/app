<script setup lang="ts">
const { loading: { onMountedWithLoading }, anyBlockingModalOpen } = useModal()

onMountedWithLoading(() => { /* nothing to do */ }, 'defaultLayout.onMountedWithLoading')

// TODO(grady) Do i18n setup for SEO here. https://i18n.nuxtjs.org/guide/seo#setup
</script>

<template>
  <div class="app-default-layout">
    <StandardNav />
    <div
      class="flex flex-column align-items-center relative"
      :aria-hidden="anyBlockingModalOpen"
    >
      <main
        class="px-3 md:px-4 w-full lg:w-10 xl:w-8 mx-auto"
        style="min-height: calc(100vh - 10.5rem - 4px);"
      >
        <NuxtPage />
      </main>
    </div>
    <ModalGroup />
    <StandardFooter />
  </div>
</template>
