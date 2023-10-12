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
        class="px-3 md:px-6 w-full lg:w-10 xl:w-8 mx-auto full-height-minus-header-footer"
      >
        <NuxtPage />
      </main>
    </div>
    <ModalGroup />
    <StandardFooter />
  </div>
</template>
