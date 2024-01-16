<script setup lang="ts">
import { type LanguageOption, LanguageOptions, languageToOption } from '@/lib/language'
import { type Language } from '@/openapi/generated/pacta'

const { t } = useI18n()
const tt = (key: string) => t(`components/language/Selector.${key}`)

interface Props {
  value: Language | undefined
}
interface Emits {
  (e: 'update:value', value: Language | undefined): void
}
const props = defineProps<Props>()
const emits = defineEmits<Emits>()
const model = computed<LanguageOption | undefined>({
  get: () => props.value ? languageToOption(props.value) : undefined,
  set: (v: LanguageOption | undefined) => {
    emits('update:value', v ? v.language : undefined)
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
      <LanguageRepresentation
        v-if="slotProps.value"
        :code="slotProps.value.code"
      />
      <span
        v-else
        class="font-italic font-light"
      >{{ tt('Unset') }}</span>
    </template>
    <template #option="slotProps">
      <LanguageRepresentation
        v-if="slotProps.option"
        :code="slotProps.option.code"
      />
      <span
        v-else
        class="font-italic font-light"
      >{{ tt('Unset') }}</span>
    </template>
  </PVDropdown>
</template>
