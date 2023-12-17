<script setup lang="ts">
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
  offlineSafe: undefined,
})
const tooltip = computed(() => {
  const result = { value: props.alt, class: props.tooltipClass }
  return result
})
</script>

<template>
  <PVButton
    v-tooltip="tooltip"
    :icon="props.icon"
    :aria-label="props.alt"
    :class="props.buttonClass"
    :style="props.buttonStyle"
    :offline-safe="props.offlineSafe"
  />
</template>
