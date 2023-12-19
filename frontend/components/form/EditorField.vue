<script setup lang="ts">
import { useSlots } from 'vue'
import { type EditorField, type EditorValue, isValid } from '@/lib/editor'

const { t } = useI18n()
const tt = (key: string) => t(`components/form/EditorField.${key}`)

// Why this convoluted type structure?
// In order for typchecking down the line, the EditorField and EditorValue need to have the SAME key.
// Enforcing that requires a check at a higher level, otherwise we'd have to directly parameterize Props
// with a non-any parameter value, which is a no-no. This leads us to Indirect1.
// Then, we have a second problem: we want to specify the `any`
// on props, but we don't want to use `keyof any`, which doesn't guarantee that the key will correspond to the parameterized type.
// The Indirect2 allows us to condense these two constraints down to one, which allows us to use the single `any` in props.
interface Indirect1<T, K extends keyof T> {
  editorField: EditorField<T, K>
  editorValue: EditorValue<T, K>
}

interface Indirect2<T> extends Indirect1<T, keyof T> {}

interface Props extends Indirect2<any> {
  isLoading?: boolean
}
const props = withDefaults(defineProps<Props>(), {
  isLoading: false,
})
const slots = useSlots()

const helpTextSlotExists = computed(() => slots['help-text'] !== undefined)
const valid = computed(() => isValid(props.editorField, props.editorValue))
const hasValidation = computed(() => (props.editorField.validation ?? []).length > 0)
const loadingLabel = computed(() => props.editorField.loadingLabel ?? tt('Loading'))
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
