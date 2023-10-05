<script setup lang="ts">
import { useSlots } from 'vue'
interface Props {
  label: string
  helpText?: string
  startHelpTextExpanded?: boolean
  required?: boolean
  loading?: boolean
  completed?: boolean
  requiredLabel?: string
  loadingLabel?: string
  completedLabel?: string
}
const props = withDefaults(defineProps<Props>(), {
  helpText: '',
  startHelpTextExpanded: true,
  required: false,
  loading: false,
  completed: false,
  requiredLabel: 'Required',
  loadingLabel: 'Loading...',
  completedLabel: '',
})
const slots = useSlots()

const helpTextExists = computed(() => props.helpText !== '' || slots['help-text'] !== undefined)
</script>

<template>
  <div class="flex flex-column form-field">
    <FormFieldHeader
      :label="props.label"
      :help-text="props.helpText"
      :help-text-exists="helpTextExists"
      :start-help-text-expanded="props.startHelpTextExpanded"
      :required="props.required"
      :loading="props.loading"
      :completed="props.completed"
      :required-label="props.requiredLabel"
      :loading-label="props.loadingLabel"
      :completed-label="props.completedLabel"
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
