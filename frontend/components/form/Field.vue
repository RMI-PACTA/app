<script setup lang="ts">
import { useSlots } from 'vue'

const slots = useSlots()
const { t } = useI18n()

const prefix = 'components/form/Field'
const tt = (s: string) => t(`${prefix}.${s}`)

interface Props {
  label: string
  helpText?: string
  startHelpTextExpanded?: boolean
  isLoading?: boolean
  loadingLabel?: string
  hasValidation?: boolean
  isValid?: boolean
  invalidLabel?: string
  validLabel?: string
}
const props = withDefaults(defineProps<Props>(), {
  helpText: '',
  startHelpTextExpanded: true,
  isLoading: false,
  loadingLabel: undefined,
  hasValidation: false,
  isValid: false,
  invalidLabel: undefined,
  validLabel: '',
})

const invalidLabel = computed(() => props.invalidLabel ?? tt('Needs Attention'))
const loadingLabel = computed(() => props.loadingLabel ?? tt('Loading'))

const helpTextExists = computed(() => props.helpText !== '' || slots['help-text'] !== undefined)
</script>

<template>
  <div class="flex flex-column form-field">
    <FormFieldHeader
      :label="props.label"
      :help-text="props.helpText"
      :help-text-exists="helpTextExists"
      :start-help-text-expanded="props.startHelpTextExpanded"
      :is-loading="props.isLoading"
      :loading-label="loadingLabel"
      :has-validation="props.hasValidation"
      :is-valid="props.isValid"
      :invalid-label="invalidLabel"
      :valid-label="props.validLabel"
    >
      <template #help-text>
        <slot name="help-text" />
      </template>
    </FormFieldHeader>
    <slot />
  </div>
</template>

<style scoped lang="scss">
.form-field {
  margin-bottom: 1.5rem;
}
</style>
