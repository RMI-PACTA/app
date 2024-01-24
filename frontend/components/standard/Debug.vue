<script setup lang="ts">
const { showStandardDebug } = useLocalStorage()
const { t } = useI18n()

const prefix = 'components/standard/Debug'
const tt = (s: string) => t(`${prefix}.${s}`)

interface Props {
  label?: string
  always?: boolean
  noAccordion?: boolean
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
    v-if="(showStandardDebug || props.always) && !props.noAccordion"
    class="standard-debug"
  >
    <PVAccordionTab
      :header="label"
      content-class="surface-100"
      header-class="surface-800"
    >
      <div
        class="code-block surface-50 relative min-w-full"
      >
        <div class="absolute top-0 right-0 flex gap-0">
          <CopyToClipboardButton
            :value="valueAsStr"
            class="p-button-text p-button-secondary"
          />
          <DownloadButton
            :value="valueAsStr"
            :file-name="`${props.label ?? 'pacta-metadata'}.json`"
            class="p-button-text p-button-secondary"
          />
        </div>
        {{ valueAsStr }}
      </div>
    </PVAccordionTab>
  </PVAccordion>
</template>

<style lang="scss">
  .standard-debug.p-accordion {
    width: fit-content;
    display: inline-block;

    .p-accordion-header .p-accordion-header-link {
      gap: 1rem;
      padding: 0.5rem 0.75rem;
    }
  }
  // Because accordions within accordions have poor styling, we hardcode these properties
  .standard-debug.p-accordion {
    .p-accordion-header-text {
      font-size: .9rem;
    }
    .p-accordion-header .p-accordion-header-link {
      border: 1px solid #a7a9ac !important; 
      border-radius: 2px !important;
    }
    .p-accordion-header.p-highlight .p-accordion-header-link {
      border-radius: 2px 2px 0 0 !important;
    }
  }
</style>
