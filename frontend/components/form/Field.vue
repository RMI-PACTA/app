<script setup lang="ts">
import { useSlots } from 'vue'

const slots = useSlots()
const { t } = useI18n()

const prefix = 'FormField'
const tt = (s: string) => t(`${prefix}.${s}`)

interface Props {
  label: string
  helpText?: string
  startHelpTextExpanded?: boolean
  required?: boolean
  loading?: boolean
  completed?: boolean
  requiredLabel?: string | undefined
  loadingLabel?: string | undefined
  completedLabel?: string
}
const props = withDefaults(defineProps<Props>(), {
  helpText: '',
  startHelpTextExpanded: true,
  required: false,
  loading: false,
  completed: false,
  requiredLabel: undefined,
  loadingLabel: undefined,
  completedLabel: '',
})

const requiredLabel = computed(() => props.loadingLabel ?? tt('Required'))
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
      :required="props.required"
      :loading="props.loading"
      :completed="props.completed"
      :required-label="requiredLabel"
      :loading-label="loadingLabel"
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
