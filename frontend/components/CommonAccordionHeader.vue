<script setup lang="ts">
interface Props {
  heading: string
  subHeading?: string | undefined
  icon?: string | undefined
}
const props = defineProps<Props>()
</script>

<template>
  <div class="flex justify-content-between align-items-center w-full flex-wrap px-2">
    <div class="flex gap-1 flex-column">
      <div class="text-lg font-bold">
        {{ props.heading }}
      </div>
      <div
        v-if="props.subHeading"
        class="text-sm font-light"
      >
        {{ props.subHeading }}
      </div>
    </div>
    <i
      v-if="props.icon"
      :class="props.icon"
    />
    <slot />
  </div>
</template>
