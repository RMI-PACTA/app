<template>
  <div />
</template>

<script setup lang="ts">
const { permissionDenied: { setPermissionDenied } } = useModal()

setPermissionDenied(new Error('permission denied'))
</script>
