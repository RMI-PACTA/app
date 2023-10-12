<script setup lang="ts">
type GenericLang = 'en' | 'de' | 'fr' | 'es'
interface Option { label: string, value: GenericLang }

interface Props {
  value: GenericLang
}
const props = defineProps<Props>()
const emits = defineEmits<(e: 'update:value', value: GenericLang) => void>()
const value = computed<GenericLang>({
  get: () => props.value,
  set: (v: GenericLang) => { emits('update:value', v) },
})

const options: Option[] = [
  { value: 'en', label: '🇬🇧 English' },
  { value: 'fr', label: '🇫🇷 Français' },
  { value: 'de', label: '🇩🇪 Deutsch' },
  { value: 'es', label: '🇪🇸 Español' },
]
</script>

<template>
  <PVDropdown
    v-model="value"
    option-label="label"
    option-value="value"
    :options="options"
  />
</template>
