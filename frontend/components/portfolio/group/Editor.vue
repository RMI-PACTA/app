<script setup lang="ts">
import {
  type EditorPortfolioGroupFields as EditorFields,
  type EditorPortfolioGroupValues as EditorValues,
} from '@/lib/editor'

interface Props {
  editorValues: EditorValues
  editorFields: EditorFields
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
  </div>
</template>
