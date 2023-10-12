<script setup lang="ts">
import { type LanguageCode, LanguageOptions } from '@/lib/language'

interface Props {
  value: LanguageCode
}
interface Emits {
  (e: 'update:value', value: LanguageCode): void
}
const props = defineProps<Props>()
const emits = defineEmits<Emits>()
const value = computed<LanguageCode>({
  get: () => props.value,
  set: (v: LanguageCode) => { emits('update:value', v) },
})
</script>

<template>
  <PVDropdown
    v-model="value"
    option-label="label"
    option-value="code"
    :options="LanguageOptions"
  >
    <template #value="slotProps">
      <LanguageRepresentation :code="slotProps.value" />
    </template>
    <template #option="slotProps">
      <LanguageRepresentation :code="slotProps.option.code" />
    </template>
  </PVDropdown>
</template>
