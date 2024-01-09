<script setup lang="ts">
import { type LanguageOption, LanguageOptions, languageToOption } from '@/lib/language'
import { type Language } from '@/openapi/generated/pacta'

interface Props {
  value: Language
}
interface Emits {
  (e: 'update:value', value: Language): void
}
const props = defineProps<Props>()
const emits = defineEmits<Emits>()
const model = computed<LanguageOption>({
  get: () => languageToOption(props.value),
  set: (v: LanguageOption) => {
    emits('update:value', v.language)
  },
})
</script>

<template>
  <PVDropdown
    v-model="model"
    option-label="label"
    :options="LanguageOptions"
  >
    <template #value="slotProps">
      <LanguageRepresentation :code="slotProps.value.code" />
    </template>
    <template #option="slotProps">
      <LanguageRepresentation :code="slotProps.option.code" />
    </template>
  </PVDropdown>
</template>
