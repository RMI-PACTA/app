<script setup lang="ts">
import { OptionalBoolean } from '@/openapi/generated/pacta'

interface Props {
  value: OptionalBoolean
  trueLabel: string
  falseLabel: string
  unsetLabel: string
  disabled?: boolean
}
interface Emits {
  (e: 'update:value', value: OptionalBoolean): void
}
const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const model = computed({
  get: () => {
    if (props.value === OptionalBoolean.OPTIONAL_BOOLEAN_TRUE) {
      return true
    } else if (props.value === OptionalBoolean.OPTIONAL_BOOLEAN_FALSE) {
      return false
    } else {
      return undefined
    }
  },
  set: (v: boolean | undefined) => {
    if (v === true) {
      emit('update:value', OptionalBoolean.OPTIONAL_BOOLEAN_TRUE)
    } else if (v === false) {
      emit('update:value', OptionalBoolean.OPTIONAL_BOOLEAN_FALSE)
    } else {
      emit('update:value', OptionalBoolean.OPTIONAL_BOOLEAN_UNSET)
    }
  },
})
const label = computed(() => {
  if (props.value === OptionalBoolean.OPTIONAL_BOOLEAN_TRUE) {
    return props.trueLabel
  } else if (props.value === OptionalBoolean.OPTIONAL_BOOLEAN_FALSE) {
    return props.falseLabel
  } else {
    return props.unsetLabel
  }
})
</script>

<template>
  <div class="flex flex-wrap gap-3 align-items-center">
    <PVTriStateCheckbox
      v-model="model"
      :disabled="props.disabled"
    />
    <span>{{ label }}</span>
  </div>
</template>
