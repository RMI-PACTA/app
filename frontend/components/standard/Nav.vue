<script setup lang="ts">
import { type MenuItem } from 'primevue/menuitem'

const { t } = useI18n()
const localePath = useLocalePath()
const { showStandardDebug } = useLocalStorage()
const { isAuthenticated, signIn, signOut } = await useMSAL()
const router = useRouter()

const prefix = 'StandardNav'
const tt = (s: string) => t(`${prefix}.${s}`)
const menuHidden = useState<boolean>(`${prefix}.menuHidden`, () => false)

const menuStyles = computed(() => {
  return {
    transition: menuHidden.value ? 'max-height .1s ease' : 'max-height .5s ease',
    overflow: 'hidden',
    'max-height': menuHidden.value ? '0px' : '100vh',
    border: menuHidden.value ? undefined : '2px solid',
    'margin-top': menuHidden.value ? '0' : '-2px',
  }
})

const menuItems = computed(() => {
  const result: MenuItem[] = [
    {
      to: localePath('/'),
      label: tt('Home'),
    },
    {
      to: 'https://github.com/RMI-PACTA/app/issues/new',
      label: tt('File a Bug'),
    },
  ]
  if (showStandardDebug) {
    result.push({
      label: tt('Admin'),
      to: localePath('/admin'),
    })
  }
  if (isAuthenticated.value) {
    result.push({
      label: tt('Sign Out'),
      command: () => { void signOut() },
    })
  } else {
    result.push({
      label: tt('Sign In'),
      command: () => { void signIn() },
    })
  }
  return result
})
</script>

<template>
  <div class="w-full flex sm:justify-content-between border-bottom-2 border-primary p-4 pb-3 column-gap-5 flex-column sm:flex-row">
    <div class="flex w-full sm:w-auto align-items-center justify-content-between gap-2">
      <RMILogo
        background="white"
        layout="horizontal"
        class="h-3rem pb-2"
        @click="() => router.push('/')"
      />
      <PVButton
        :icon="menuHidden ? 'pi pi-bars' : 'pi pi-times'"
        :class="menuHidden ? 'p-button-text' : 'border-bottom-noround p-button-primary'"
        class="sm:hidden p-button-lg h-3rem"
        @click="() => menuHidden = !menuHidden"
      />
    </div>
    <div
      class="flex gap-2 sm:p-1 flex-1 flex-column sm:flex-row border-primary sm:border-none sm:max-h-full border-round justify-content-end"
      :style="menuStyles"
    >
      <template
        v-for="(mi, index) in menuItems"
      >
        <LinkButton
          v-if="mi.to"
          :key="index"
          :class="mi.to === router.currentRoute.value.fullPath ? 'border-noround sm:border-round' : 'p-button-text'"
          :to="mi.to"
          :label="`${mi.label}`"
        />
        <PVButton
          v-else
          :key="mi.label"
          :label="mi.label"
          class="p-button-text"
          @click="mi.command"
        />
      </template>
    </div>
  </div>
</template>
