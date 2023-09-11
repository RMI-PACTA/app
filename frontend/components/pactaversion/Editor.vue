<script setup lang="ts">
import { type PactaVersion } from '@/openapi/generated/pacta'
import FormField from '@/components/form/Field.vue'

const props = defineProps<{
  pactaVersion: PactaVersion
}>()
const emit = defineEmits<(e: 'update:pactaVersion', pactaVersion: PactaVersion) => void>()
const model = computed({
  get: () => props.pactaVersion,
  set: (pactaVersion: PactaVersion) => { emit('update:pactaVersion', pactaVersion) }
})

const nameCompleted = computed(() => model.value.name.length > 0)
const digestCompleted = computed(() => model.value.digest.length > 0)
const incompleteFields = computed<string[]>(() => {
  const result: string[] = []
  if (!nameCompleted.value) { result.push('Version Name') }
  if (!digestCompleted.value) { result.push('Docker Image Digest') }
  return result
})

defineExpose({ incompleteFields })
</script>

<template>
  <div>
    <FormField
      label="Version Name"
      help-text="The name of the version of the PACTA algorithm."
      required
      :completed="nameCompleted"
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
