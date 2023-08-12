<template>
  <div class="app-default-layout">
    <StandardNav />
    <div
      class="flex flex-column align-items-center relative"
    >
      <main
        class="px-3 md:px-4 w-full lg:w-10 xl:w-8 mx-auto"
        style="min-height: calc(100vh - 9rem - 4px);"
      >
        <NuxtPage />
      </main>
    </div>
    <StandardFooter />
  </div>
</template>
