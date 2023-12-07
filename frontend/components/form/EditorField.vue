<script setup lang="ts">
import { useSlots } from 'vue'
import { type EditorField, type EditorValue, isValid } from '@/lib/editor'

const { t } = useI18n()
const tt = (key: string) => t(`FormEditorField.${key}`)

interface Props {
  editorField: EditorField<any, keyof any>
  editorValue: EditorValue<any, keyof any>
  helpText?: string
  startHelpTextExpanded?: boolean
  isLoading?: boolean
  loadingLabel?: string
  invalidLabel?: string
  validLabel?: string
}
const props = withDefaults(defineProps<Props>(), {
  helpText: undefined,
  startHelpTextExpanded: false,
  loading: false,
  loadingLabel: tt('Loading...'),
  invalidLabel: tt('Needs Attention'),
  validLabel: '',
})
const slots = useSlots()

const helpTextSlotExists = computed(() => slots['help-text'] !== undefined)
const valid = computed(() => isValid(props.editorValue))
const hasValidation = computed(() => (props.editorValue.validation ?? []).length > 0)
</script>

<template>
  <FormField
    :label="props.editorField.label"
    :help-text="props.helpText"
    :start-help-text-expanded="props.startHelpTextExpanded"
    :is-loading="props.loading"
    :loading-label="props.loadingLabel"
    :has-validation="hasValidation"
    :is-valid="valid"
    :invalid-label="props.invalidLabel"
    :valid-label="props.validLabel"
  >
    <template
      v-if="helpTextSlotExists"
      #help-text
    >
      <slot name="help-text" />
    </template>
    <slot />
  </FormField>
</template>
