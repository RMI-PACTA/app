<script setup lang="ts">
const { showStandardDebug } = useLocalStorage()

interface Props {
  label?: string
  always?: boolean
  value: unknown
}
const props = withDefaults(defineProps<Props>(), { always: false, label: 'Technical Debugging' })
</script>

<template>
  <ClientOnly>
    <PVAccordion
      v-if="showStandardDebug || props.always"
      class="standard-debug"
    >
      <PVAccordionTab
        :header="props.label || 'Debug'"
        content-class="surface-100"
        header-class="surface-800"
      >
        <div
          class="code surface-50"
        >
          {{ JSON.stringify(props.value, null, 2) }}
        </div>
      </PVAccordionTab>
    </PVAccordion>
  </ClientOnly>
</template>

<style lang="scss">
  .standard-debug.p-accordion {
    width: fit-content;
    display: inline-block;

    .code {
      display: inline-block;
      font-size: 0.75rem;
      line-height: 0.75rem;
      white-space: pre-wrap;
      font-family: monospace;
    }

    .p-accordion-header .p-accordion-header-link {
      gap: 1rem;
      padding: 0.5rem 0.75rem;
    }
  }
</style>
