<script setup lang="ts">
const { permissionDenied: { permissionDeniedVisible, permissionDeniedError }, error: { error } } = useModal()
const router = useRouter()

const navigateBack = () => {
  // If the user didn't land here directly, go back to the previous page.
  if (window.history.length > 1) {
    router.back()
    return
  }
  // Otherwise, just take them home.
  return router.push('/')
}
</script>

<template>
  <StandardModal
    v-model:visible="permissionDeniedVisible"
    header="Permission Denied"
    sub-header="You aren't authorized to view this page"
    @closed="navigateBack"
  >
    <p>
      It is likely that you have gotten here either by entering in an invalid
      URL, or by clicking on a link or button that is not set up correctly.
    </p>
    <p>
      If you think this is a bug (something that you should have access
      to), or if you want to let us know how you got here, please file a bug.
    </p>
    <StandardDebug
      label="Technical Error"
      :value="permissionDeniedError"
      always
    />
  </StandardModal>
</template>
