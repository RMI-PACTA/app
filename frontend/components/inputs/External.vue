<script setup lang="ts">
import { type OptionalBoolean } from '@/openapi/generated/pacta'

const { t } = useI18n()

const tt = (s: string) => t(`components/inputs/External.${s}`)

interface Props {
  value: OptionalBoolean
  disabled?: boolean
}
const props = defineProps<Props>()

interface Emits {
  (e: 'update:value', value: OptionalBoolean): void
}
const emit = defineEmits<Emits>()

const model = computed({
  get: () => props.value,
  set: (value: OptionalBoolean) => { emit('update:value', value) },
})
</script>

<template>
  <ExplicitTriStateCheckbox
    v-model:value="model"
    :true-label="tt('True')"
    :false-label="tt('False')"
    :unset-label="tt('Unset')"
    :disabled="props.disabled"
  />
</template>
