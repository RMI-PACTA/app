<script setup lang="ts">
import { onMounted } from 'vue'

interface Props {
  title: string
}
const props = withDefaults(defineProps<Props>(), {})
onMounted(() => {
  document.title = 'RMI | PACTA | ' + props.title
})
</script>

<template>
  <h1 class="inline-block text-primary title-bar">
    {{ title }}
  </h1>
</template>

<style scoped lang="scss">
.title-bar {
  margin: 0;
}
</style>
