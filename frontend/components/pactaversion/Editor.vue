<script setup lang="ts">
import { type PactaVersion } from '@/openapi/generated/pacta'

const props = defineProps<{
  pactaVersion: PactaVersion
}>()

const emit = defineEmits<(e: 'update:pactaVersion', pactaVersion: PactaVersion) => void, >()

const model = computed({
  get: () => props.pactaVersion,
  set: (pactaVersion: PactaVersion) => { emit('update:pactaVersion', pactaVersion) }
})
</script>

<template>
  <div class="flex flex-column gap-3">
    <FormField
      label="Version Name"
      help-text="The name of the version of the PACTA algorithm."
      required
      completed="!!model.name"
    >
      <PVInputText
        v-model="model.name"
      />
    </FormField>
  </div>
</template>
