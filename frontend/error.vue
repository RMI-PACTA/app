<script setup lang="ts">
// We get the equivalent of props.error from useError() in the ModalError component.
// const props = defineProps({
//   error: Object
// })

const { error: { errorModalVisible } } = useModal()
errorModalVisible.value = true
</script>

<template>
  <ModalError :is-full-page="true" />
</template>
