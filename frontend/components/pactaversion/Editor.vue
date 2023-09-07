<script setup lang="ts">
import { type PactaVersion } from '@/openapi/generated/pacta'

const props = defineProps<{
  pactaVersion: PactaVersion
}>()

const emit = defineEmits<(e: 'update:pactaVersion', pactaVersion: PactaVersion) => void>()

const model = computed({
  get: () => props.pactaVersion,
  set: (pactaVersion: PactaVersion) => { emit('update:pactaVersion', pactaVersion) }
})
</script>

<template>
  <div>
    <FormField
      label="Version Name"
      help-text="The name of the version of the PACTA algorithm."
      required
      :completed="model.name.length > 0"
    >
      <PVInputText
        v-model="model.name"
      />
    </FormField>
    <FormField
      label="Version Description"
      help-text="An optional description of this version of the PACTA algorithm."
    >
      <PVTextarea
        v-model="model.description"
        auto-resize
      />
    </FormField>
    <FormField
      label="Docker Image Digest"
      help-text="The SHA hash of the docker image that should correspond to this version of the PACTA version."
      required
      :completed="model.digest.length > 0"
    >
      <PVInputText
        v-model="model.digest"
      />
    </FormField>
  </div>
</template>
