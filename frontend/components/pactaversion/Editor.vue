<script setup lang="ts">
import {
  type EditorPactaVersionFields as EditorFields,
  type EditorPactaVersionValues as EditorValues,
} from '@/lib/editor'

interface Props {
  editorFields: EditorFields
  editorValues: EditorValues
}
interface Emits {
  (e: 'update:editorValues', evs: EditorValues): void
}
const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const efs = computed(() => props.editorFields)
const evs = computed({
  get: () => props.editorValues,
  set: (evs) => { emit('update:editorValues', evs) },
})
</script>

<template>
  <div>
    <FormEditorField
      :editor-field="efs.name"
      :editor-value="evs.name"
    >
      <PVInputText
        v-model="evs.name.currentValue"
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="efs.description"
      :editor-value="evs.description"
    >
      <PVTextarea
        v-model="evs.description.currentValue"
        auto-resize
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="efs.digest"
      :editor-value="evs.digest"
    >
      <PVInputText
        v-model="evs.digest.currentValue"
      />
    </FormEditorField>
  </div>
</template>
