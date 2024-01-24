<script setup lang="ts">
import { type HoldingsDate } from '@/openapi/generated/pacta'

interface Props {
  value: HoldingsDate | undefined
  disabled?: boolean
}
const props = defineProps<Props>()

interface Emits {
  (e: 'update:value', value: HoldingsDate | undefined): void
}
const emit = defineEmits<Emits>()

const model = computed<Date | undefined>({
  get: () => props.value?.time ? new Date(props.value.time) : undefined,
  set: (value: Date | undefined) => { emit('update:value', value ? ({ time: value.toISOString() }) : ({ time: undefined })) },
})
</script>

<template>
  <PVCalendar
    v-model="model"
    view="month"
    date-format="mm/yy"
    show-button-bar
    :disabled="props.disabled"
  />
</template>
