<script setup lang="ts">
import { useSlots } from 'vue'
import { type EditorField, type EditorValue, isValid } from '@/lib/editor'

const { t } = useI18n()
const tt = (key: string) => t(`components/form/EditorField.${key}`)

interface Props {
  editorField: EditorField<any, keyof any>
  editorValue: EditorValue<any, keyof any>
  isLoading?: boolean
}
const props = withDefaults(defineProps<Props>(), {
  isLoading: false,
})
const slots = useSlots()

const helpTextSlotExists = computed(() => slots['help-text'] !== undefined)
const valid = computed(() => isValid(props.editorField, props.editorValue))
const hasValidation = computed(() => (props.editorField.validation ?? []).length > 0)
const loadingLabel = computed(() => props.editorField.loadingLabel ?? tt('Loading...'))
const invalidLabel = computed(() => props.editorField.invalidLabel ?? tt('Needs Attention'))
const validLabel = computed(() => props.editorField.validLabel ?? '')
const startHelpTextExpanded = computed(() => props.editorField.startHelpTextExpanded ?? false)
</script>

<template>
  <FormField
    :label="props.editorField.label"
    :help-text="props.editorField.helpText"
    :start-help-text-expanded="startHelpTextExpanded"
    :is-loading="props.isLoading"
    :loading-label="loadingLabel"
    :has-validation="hasValidation"
    :is-valid="valid"
    :invalid-label="invalidLabel"
    :valid-label="validLabel"
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
