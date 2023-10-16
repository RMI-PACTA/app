<script setup lang="ts">
const { fakeUsers: { fakeUsersVisible } } = useModal()
const { t } = useI18n()
const { getMaybeMe, refreshMaybeMe } = useSession()
const { isAuthenticated, signIn, signOut } = await useMSAL()
const localePath = useLocalePath()

const prefix = 'ModalFakeUsers'
const tt = (s: string) => t(`${prefix}.${s}`)

// TODO(grady) move this to an env-specific location, omit the component if empty.
const hardcodedLocalUsers = [{
  id: 'user.c9a959a1c4d6b8be1fc1',
  email: 'rmi-localtest-a@siliconally.org',
  password: 'CQdCZd9lq2HOnFB',
}, {
  id: 'user.4aed45b19c7ef24e105c',
  email: 'rmi-localtest-b@siliconally.org',
  password: 'Dq7CLBX90nFtNk8',
}]

const { maybeMe } = await getMaybeMe()

const users = computed(() => {
  const me = maybeMe.value
  return hardcodedLocalUsers.map((u) => ({
    itMe: me && u.email.toLowerCase() === me.enteredEmail.toLowerCase(),
    ...u,
  }))
})
const close = () => { fakeUsersVisible.value = false }
const doSignOut = () => { void signOut().then(close) }
</script>

<template>
  <StandardModal
    v-model:visible="fakeUsersVisible"
    :header="tt('Heading')"
    :sub-header="tt('Subheading')"
    @show="refreshMaybeMe"
  >
    <p>
      {{ tt('Description') }}
    </p>
    <div class="flex flex-column gap-2 align-content-stretch">
      <div
        v-for="user in users"
        :key="user.email"
        class="shadow-3 p-2 border-round"
        :class="user.itMe ? 'border-primary border-2' : ''"
      >
        <div class="flex flex-column gap-2">
          <div class="flex gap-2 justify-content-between align-items-center">
            <span>
              {{ tt('Email') }}:
              <b class="code">{{ user.email }}</b>
              <LinkButton
                :to="localePath(`/user/${user.id}`)"
                class="p-button-xs p-button-secondary p-button-text"
                icon="pi pi-external-link"
                @click="close"
              />
            </span>
            <CopyToClipboardButton
              v-if="!user.itMe"
              :value="user.email"
              icon="pi pi-user"
              :cta="tt('Copy')"
              class="p-button-xs"
            />
          </div>
          <div class="flex gap-2 justify-content-between align-items-center">
            <span>{{ tt('Password') }}: <b class="code">{{ user.password }}</b></span>
            <CopyToClipboardButton
              v-if="!user.itMe"
              :value="user.password"
              :cta="tt('Copy')"
              class="p-button-xs"
            />
            <PVButton
              v-else
              :label="tt('Sign Out')"
              icon="pi pi-sign-out"
              class="p-button-danger p-button-xs"
              @click="doSignOut"
            />
          </div>
        </div>
      </div>
    </div>
    <PVButton
      v-if="!isAuthenticated"
      :label="tt('Sign In')"
      icon="pi pi-sign-in"
      class="align-self-start"
      @click="signIn"
    />
  </StandardModal>
</template>
