<script setup lang="ts">
const { showStandardDebug } = useLocalStorage()
const { t } = useI18n()

const prefix = 'components/standard/Debug'
const tt = (s: string) => t(`${prefix}.${s}`)

interface Props {
  label?: string
  always?: boolean
  value: unknown
}
const props = withDefaults(defineProps<Props>(), { always: false, label: undefined })
const label = computed(() => props.label ?? tt('Debugging Information'))
const valueAsStr = computed(() => JSON.stringify(props.value, createCircularReplacer(), 2))

function createCircularReplacer (): (this: any, key: string, value: any) => any {
  const seen = new WeakSet()
  return function (this: any, key: string, value: any) {
    if (typeof value === 'object' && value !== null) {
      if (seen.has(value)) {
        return '#REF'
      }
      seen.add(value)
    }
    return value
  }
}

</script>

<template>
  <PVAccordion
    v-if="showStandardDebug || props.always"
    class="standard-debug"
  >
    <PVAccordionTab
      :header="label"
      content-class="surface-100"
      header-class="surface-800"
    >
      <div
        class="code surface-50"
      >
        {{ valueAsStr }}
      </div>
    </PVAccordionTab>
  </PVAccordion>
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
