<script setup lang="ts">
interface Props {
  onLabel: string
  offLabel: string
  value: boolean
}
const props = defineProps<Props>()

const emit = defineEmits<(e: 'update:value', value: boolean) => void>()

const value = computed({
  get: () => props.value,
  set: (v: boolean) => { emit('update:value', v) },
})
</script>

<template>
  <div class="flex flex-wrap gap-3 align-items-center">
    <PVInputSwitch
      v-model="value"
    />
    <span>{{ value ? props.onLabel : props.offLabel }}</span>
  </div>
</template>
