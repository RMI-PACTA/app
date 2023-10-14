<script setup lang="ts">
import { type EditorPactaVersion } from '@/lib/editor'

interface Props {
  editorPactaVersion: EditorPactaVersion
}
interface Emits {
  (e: 'update:editorPactaVersion', epv: EditorPactaVersion): void
}
const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const epv = computed({
  get: () => props.editorPactaVersion,
  set: (epv) => { emit('update:editorPactaVersion', epv) },
})
</script>

<template>
  <div>
    <FormEditorField
      :editor-field="epv.name"
      help-text="The name of the version of the PACTA algorithm."
    >
      <PVInputText
        v-model="epv.name.currentValue"
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="epv.description"
      help-text="An optional description of this version of the PACTA algorithm."
    >
      <PVTextarea
        v-model="epv.description.currentValue"
        auto-resize
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="epv.digest"
      help-text="The SHA hash of the docker image that should correspond to this version of the PACTA version."
    >
      <PVInputText
        v-model="epv.digest.currentValue"
      />
    </FormEditorField>
  </div>
</template>
