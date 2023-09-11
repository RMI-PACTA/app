<script setup lang="ts">
const { showStandardDebug } = useLocalStorage()
</script>

<template>
  <div class="border-top-2 border-primary w-full text-primary py-3 px-4 flex gap-2 flex-wrap justify-content-between align-items-center">
    <span>A project of <a href="https://rmi.org">Rocky Mountain Institute</a></span>
    <div class="flex gap-2 align-items-center">
      <PVButton
        icon="pi pi-bolt"
        class="p-0 w-auto"
        :class="showStandardDebug ? '' : 'p-button-text'"
        @click="() => showStandardDebug = !showStandardDebug"
      />
      <span>© 2023 RMI</span>
    </div>
  </div>
</template>
