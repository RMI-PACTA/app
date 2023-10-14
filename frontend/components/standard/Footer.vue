<script setup lang="ts">
const { showStandardDebug } = useLocalStorage()
const { t } = useI18n()

const prefix = 'StandardFooter'
const tt = (s: string) => t(`${prefix}.${s}`)
</script>

<template>
  <div class="border-top-2 border-primary w-full text-primary py-3 px-4 flex gap-2 flex-wrap justify-content-between align-items-center">
    <div class="flex flex-column gap-1">
      <span>{{ tt('A project of') }} <a
        href="https://rmi.org"
        class="text-primary"
      >Rocky Mountain Institute</a></span>
      <span class="text-600">© 2023 RMI</span>
    </div>
    <div class="flex flex-row gap-3 align-items-center">
      <PVButton
        icon="pi pi-bolt"
        class="p-1 w-auto"
        :class="showStandardDebug ? '' : 'p-button-text'"
        @click="() => showStandardDebug = !showStandardDebug"
      />
      <div class="flex row-gap-1 column-gap-3 flex-wrap">
        <a
          class="text-primary"
          href="#"
        >{{ tt('Terms of Use') }}</a>
        <a
          class="text-primary"
          href="#"
        >{{ tt('Privacy') }} </a>
        <a
          class="text-primary"
          href="#"
        >{{ tt('About') }}</a>
      </div>
      <LocaleSelector />
    </div>
  </div>
</template>
