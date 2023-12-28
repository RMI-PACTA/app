<script setup lang="ts">
const isAuthenticated = useIsAuthenticated()
const { $msal: { signOut } } = useNuxtApp()
const { signIn } = useSignIn()
</script>

<template>
  <StandardContent>
    <TitleBar title="RMI PACTA" />
    <p>
      TODO(#80) This will eventually be the site for RMI's PACTA, but now it's mostly just a placeholder.
    </p>
    <p>
      This project is open source. You can view the code at <a href="https://github.com/RMI-PACTA/app">github.com/RMI-PACTA/app</a>.
    </p>

    <PVButton
      v-if="!isAuthenticated"
      label="Sign In"
      icon="pi pi-sign-in"
      icon-pos="right"
      @click="signIn"
    />
    <template v-else>
      <LinkButton
        label="Get Started"
        icon="pi pi-arrow-right"
        icon-pos="right"
        to="/upload"
      />
      <PVButton
        class="p-button-outlined"
        label="Sign Out"
        icon="pi pi-sign-out"
        icon-pos="right"
        @click="signOut"
      />
    </template>
  </StandardContent>
</template>
