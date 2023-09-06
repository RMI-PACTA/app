<script setup lang="ts">
import { computed } from 'vue'
interface Props {
  label: string
  helpText: string
  startHelpTextExpanded: boolean
  helpTextExists: boolean
  required: boolean
  completed: boolean
}
const props = defineProps<Props>()
const { helpTextExpanded: computedHTE } = useLocalStorage()

const id = `FormField[${useStateIDGenerator().id()}]`
const helpTextExpanded = computedHTE(props.label)
const helpTextIconClass = computed(() => helpTextExpanded.value ? 'pi pi-info-circle' : 'pi pi-info-circle text-600')
const helpTextTextClass = computed(() => helpTextExpanded.value ? 'mb-2' : 'h-0')
</script>

<template>
  <div class="flex flex-column">
    <div class="flex align-items-center mb-1 gap-2">
      <label
        class="inline-block text-lg ml-1"
        :for="id"
      >
        {{ props.label }}
      </label>
      <i
        v-if="helpTextExists"
        :class="helpTextIconClass"
        class="cursor-pointer p-1"
        @click="() => helpTextExpanded = !helpTextExpanded"
      />
      <div
        v-if="props.required && !props.completed"
        class="required-warning flex align-items-center gap-1"
      >
        <i
          class="pi pi-exclamation-triangle"
        />
        <span>Required</span>
      </div>
      <div
        v-if="props.required && props.completed"
        class="completed-warning flex align-items-center gap-1"
      >
        <i
          class="pi pi-check-circle"
        />
        <span>Completed</span>
      </div>
    </div>
    <div
      v-if="helpTextExists"
      :class="helpTextTextClass"
      class="overflow-hidden ml-1 text-sm help-text-animate"
    >
      <slot name="help-text" />
      {{ props.helpText }}
    </div>
  </div>
</template>

<style scoped lang="scss">
.required-warning {
  background: $warningMessageBg;
  border-radius: $borderRadius;
  padding: 0.25rem 0.5rem;
}
</style>
