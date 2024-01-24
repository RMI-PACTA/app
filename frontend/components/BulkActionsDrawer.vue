<script setup lang="ts">
const { t } = useI18n()

const prefix = 'components/BulkActionsDrawer'
const statePrefix = `${prefix}[${useStateIDGenerator().id()}]`
const tt = (s: string) => t(`${prefix}.${s}`)

const open = useState<boolean>(`${statePrefix}.open`, () => false)
const icon = computed(() => open.value ? 'pi pi-chevron-left' : 'pi pi-chevron-right')
const classes = computed(() => open.value ? 'border-2 border-primary' : 'p-button-outlined border-round')
</script>

<template>
  <span class="p-buttonset w-fit">
    <PVButton
      :label="tt('Bulk Actions')"
      class="p-button-sm"
      :class="classes"
      :icon="icon"
      icon-pos="right"
      @click="() => { open = !open}"
    />
    <template v-if="open">
      <slot />
    </template>
  </span>
</template>
