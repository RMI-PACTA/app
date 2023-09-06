<script setup lang="ts">
import { computed, useSlots } from 'vue'
interface Props {
  label: string
  helpText?: string
  startHelpTextExpanded?: boolean
  required?: boolean
  completed?: boolean
}
const props = withDefaults(defineProps<Props>(), {
  helpText: '',
  startHelpTextExpanded: true,
  required: false,
  completed: false
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
      :completed="props.completed"
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
