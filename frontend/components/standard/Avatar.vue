<script setup lang="ts">
const props = defineProps<{
  name?: string | undefined
}>()

const singleCharacter = computed(() => {
  return props.name?.[0] || 'U'
})
</script>

<template>
  <div
    class="flex align-items-center justify-content-center border-circle bg-teal-400"
    style="height: 2.5rem; width: 2.5rem;"
  >
    <span class="text-2xl text-white font-bold">
      {{ singleCharacter }}
    </span>
  </div>
</template>
