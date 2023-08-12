<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  icon: string
  alt: string
  buttonClass?: string | undefined
  buttonStyle?: string | undefined
  tooltipClass?: string | undefined
  offlineSafe?: boolean | undefined
}

const props = withDefaults(defineProps<Props>(), {
  buttonClass: undefined,
  buttonStyle: undefined,
  tooltipClass: undefined,
  offlineSafe: undefined
})
const tt = computed(() => {
  const result = { value: props.alt, class: props.tooltipClass }
  return result
})
</script>

<template>
  <PVButton
    v-tooltip="tt"
    :icon="props.icon"
    :aria-label="props.alt"
    :class="props.buttonClass"
    :style="props.buttonStyle"
    :offline-safe="props.offlineSafe"
  />
</template>
