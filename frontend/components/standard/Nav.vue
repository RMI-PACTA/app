<script setup lang="ts">
import { type MenuItem } from 'primevue/menuitem'

const { t } = useI18n()
const localePath = useLocalePath()
const isAuthenticated = useIsAuthenticated()
const router = useRouter()

const { $msal: { signOut } } = useNuxtApp()
const { signIn } = useSignIn()
const { getMaybeMe } = useSession()

const { isAdmin, maybeMe } = await getMaybeMe()

const prefix = 'components/standard/Nav'
const tt = (s: string) => t(`${prefix}.${s}`)
const menuHidden = useState<boolean>(`${prefix}.menuHidden`, () => false)
const userMenu = useState<{ toggle: (e: Event) => void } | null>(`${prefix}.userMenu`, () => null)
const userMenuVisible = useState<boolean>(`${prefix}.userMenuVisible`, () => false)

const toggleUserMenu = (e: Event) => {
  if (userMenu.value) {
    userMenu.value.toggle(e)
  }
}

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
      to: 'https://pacta.rmi.org',
      external: true,
      label: tt('About'),
    },
    {
      to: localePath('/'),
      label: tt('Home'),
    },
  ]
  if (isAdmin.value) {
    result.push({
      label: tt('Admin'),
      to: localePath('/admin'),
    })
  }
  if (isAuthenticated.value) {
    result.push({
      label: tt('My Data'),
      to: localePath('/my-data'),
    })
  } else {
    result.push({
      label: tt('Sign In'),
      command: () => { void signIn() },
    })
  }
  return result
})
const userMenuItems = computed(() => {
  const result: MenuItem[] = [{
    label: tt('Account'),
    icon: 'pi pi-cog',
    to: localePath('/user/me'),
  }, {
    label: tt('My Data'),
    icon: 'pi pi-list',
    to: localePath('/my-data'),
  }, {
    label: tt('Audit Logs'),
    icon: 'pi pi-lock',
    to: localePath('/audit-logs'),
  }, {
    label: tt('Sign Out'),
    icon: 'pi pi-sign-out',
    command: () => { void signOut() },
  }]
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
        @click="() => router.push(localePath('/'))"
      />
      <div class="flex gap-1 align-items-center">
        <PVButton
          v-show="maybeMe !== undefined"
          v-tooltip="tt('Settings')"
          icon="pi pi-user"
          class="sm:hidden ml-2 flex-shrink-0"
          rounded
          :class="userMenuVisible ? 'p-button-primary' : 'p-button-text'"
          @click="toggleUserMenu"
        />
        <PVButton
          :icon="menuHidden ? 'pi pi-bars' : 'pi pi-times'"
          :class="menuHidden ? 'p-button-text' : 'border-bottom-noround p-button-primary'"
          class="sm:hidden p-button-lg h-3rem"
          @click="() => menuHidden = !menuHidden"
        />
      </div>
    </div>
    <div
      class="flex gap-2 sm:p-1 flex-1 flex-column sm:flex-row border-primary sm:border-none sm:max-h-full border-round justify-content-end align-items-center"
      :style="menuStyles"
    >
      <template
        v-for="(mi, index) in menuItems"
      >
        <LinkButton
          v-if="mi.to"
          :key="index"
          :class="mi.to === router.currentRoute.value.path ? 'border-noround sm:border-round' : 'p-button-text'"
          :to="mi.to"
          :external="mi.external"
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
      <PVButton
        v-if="isAuthenticated"
        v-tooltip.left="tt('Settings')"
        icon="pi pi-user"
        class="hidden sm:flex ml-2 flex-shrink-0"
        rounded
        :class="userMenuVisible ? 'p-button-primary' : 'p-button-outlined'"
        @click="toggleUserMenu"
      />
      <PVOverlayPanel
        ref="userMenu"
        :pt="{
          content: 'p-0',
        }"
        class="caret-primary"
        @hide="() => userMenuVisible = false"
        @show="() => userMenuVisible = true"
      >
        <div class="bg-primary p-3">
          <div class="flex gap-3 align-items-center">
            <StandardAvatar :name="maybeMe?.name" />
            <div class="flex flex-column gap-1">
              <span class="font-bold text-lg text-white">
                {{ maybeMe?.name }}
              </span>
              <span class="text-sm text-white">
                {{ maybeMe?.enteredEmail }}
              </span>
            </div>
          </div>
        </div>
        <div
          class="flex flex-column gap-1 p-2"
          @click="toggleUserMenu"
        >
          <template
            v-for="(mi, index) in userMenuItems"
          >
            <LinkButton
              v-if="mi.to"
              :key="index"
              :class="mi.to === router.currentRoute.value.fullPath ? 'border-noround sm:border-round' : 'p-button-text'"
              :to="mi.to"
              :external="mi.external"
              :icon="mi.icon"
              :label="`${mi.label}`"
            />
            <PVButton
              v-else
              :key="mi.label"
              :label="mi.label"
              :icon="mi.icon"
              class="p-button-text"
              @click="mi.command"
            />
          </template>
        </div>
      </PVOverlayPanel>
    </div>
  </div>
</template>

<style lang="scss">
.caret-primary.p-overlaypanel{
  &::before, &::after {
    border-bottom-color: var(--primary-color)
  }
}
</style>
