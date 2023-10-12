<script setup lang="ts">
import { type EditorPactaVersion, isComplete } from '@/lib/editor'

const props = defineProps<{
  editorPactaVersion: EditorPactaVersion
}>()
const emit = defineEmits<(e: 'update:editorPactaVersion', epv: EditorPactaVersion) => void>()
const epv = computed({
  get: () => props.editorPactaVersion,
  set: (epv) => { emit('update:editorPactaVersion', epv) },
})
</script>

<template>
  <div>
    <FormField
      label="Version Name"
      help-text="The name of the version of the PACTA algorithm."
      :required="epv.name.isRequired"
      :completed="isComplete(epv.name)"
    >
      <PVInputText
        v-model="epv.name.currentValue"
      />
    </FormField>
    <FormField
      label="Version Description"
      help-text="An optional description of this version of the PACTA algorithm."
      :required="epv.description.isRequired"
      :completed="isComplete(epv.description)"
    >
      <PVTextarea
        v-model="epv.description.currentValue"
        auto-resize
      />
    </FormField>
    <FormField
      label="Docker Image Digest"
      help-text="The SHA hash of the docker image that should correspond to this version of the PACTA version."
      :required="epv.digest.isRequired"
      :completed="isComplete(epv.digest)"
    >
      <PVInputText
        v-model="epv.digest.currentValue"
      />
    </FormField>
  </div>
</template>
