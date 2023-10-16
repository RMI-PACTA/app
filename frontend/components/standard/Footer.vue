<script setup lang="ts">
const { showStandardDebug } = useLocalStorage()
const { fakeUsers: { fakeUsersVisible } } = useModal()
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
    <div class="flex flex-column align-items-center">
      <div class="flex gap-1">
        <PVButton
          v-tooltip.top="tt('Manage Fake Users')"
          icon="pi pi-user-plus"
          class="p-1 w-auto"
          :class="fakeUsersVisible ? '' : 'p-button-text'"
          @click="() => fakeUsersVisible = !fakeUsersVisible"
        />
        <PVButton
          v-tooltip.top="tt('Toggle Debug Info')"
          icon="pi pi-code"
          class="p-1 w-auto"
          :class="showStandardDebug ? '' : 'p-button-text'"
          @click="() => showStandardDebug = !showStandardDebug"
        />
      </div>
      <div class="text-xs">
        {{ tt('Dev Tools') }}
      </div>
    </div>
    <div class="flex flex-row gap-3 align-items-center flex-wrap">
      <div class="flex row-gap-1 column-gap-3 flex-wrap">
        <a
          href="https://github.com/RMI-PACTA/app/issues/new"
          target="_blank"
          class="text-primary"
        >
          {{ tt('File a Bug') }}
        </a>
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
