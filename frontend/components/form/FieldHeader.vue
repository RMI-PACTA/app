<script setup lang="ts">
interface Props {
  label: string
  helpText: string
  startHelpTextExpanded: boolean
  helpTextExists: boolean
  required: boolean
  loading: boolean
  completed: boolean
  requiredLabel: string
  loadingLabel: string
  completedLabel: string
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
        class="flex align-items-center gap-1 p-error"
      >
        <i
          class="pi pi-circle"
        />
        <span>{{ props.requiredLabel }}</span>
      </div>
      <div
        v-if="props.required && props.completed"
        class=" flex align-items-center gap-1 text-success"
      >
        <i
          class="pi pi-check-circle"
        />
        <span>{{ props.completedLabel }}</span>
      </div>
      <div
        v-if="props.loading"
        class="flex align-items-center gap-1 text-700"
      >
        <i
          class="pi pi-sync pi-spin"
        />
        <span>{{ props.loadingLabel }}</span>
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
