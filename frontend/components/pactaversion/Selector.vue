<script setup lang="ts">
interface Props {
  value: string | undefined
}
interface Emits {
  (e: 'update:value', value: string | undefined): void
}
const props = defineProps<Props>()
const emits = defineEmits<Emits>()
const value = computed({
  get: () => props.value,
  set: (v: string | undefined) => { emits('update:value', v) },
})

const pactaClient = await usePACTA()
const prefix = 'components/pactaversion/Selector'
const { data: pactaVersions, refresh } = await useSimpleAsyncData(
  `${prefix}.getPactaVersions`,
  () => pactaClient.listPactaVersions(),
)
const options = computed(() => pactaVersions.value.map((pv) => ({ label: pv.name, value: pv.id })))
</script>

<template>
  <div class="flex flex-wrap gap-2 align-items-center">
    <PVDropdown
      v-model="value"
      option-label="label"
      option-value="value"
      :options="options"
      class="flex-1"
    />
    <PVButton
      v-tooltip="'Refresh PACTA Version Options'"
      icon="pi pi-sync"
      class="p-button-text p-button-secondary"
      @click="refresh"
    />
  </div>
</template>
