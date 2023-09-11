<script setup lang="ts">
const { isAuthenticated, signIn, signOut } = await useMSAL()
</script>

<template>
  <StandardContent>
    <TitleBar title="RMI PACTA" />
    <p>
      This will eventually be the site for RMI's PACTA, but now it's mostly just a placeholder.
    </p>
    <p>
      This project is open source. You can view the code at <a href="https://github.com/RMI/pacta">github.com/RMI/pacta</a>.
    </p>

    <PVButton
      v-if="!isAuthenticated"
      label="Sign In"
      @click="signIn"
    />
    <PVButton
      v-else
      class="inline-block"
      label="Sign Out"
      icon="pi pi-sign-out"
      @click="signOut"
    />
  </StandardContent>
</template>
