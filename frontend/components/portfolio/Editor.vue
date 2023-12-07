<script setup lang="ts">
import { type EditorPortfolio } from '@/lib/editor'

interface Props {
  editorPortfolio: EditorPortfolio
}
interface Emits {
  (e: 'update:editorPortfolio', ei: EditorPortfolio): void
}
const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const model = computed({
  get: () => props.editorPortfolio,
  set: (editorPortfolio: EditorPortfolio) => { emit('update:editorPortfolio', editorPortfolio) },
})
</script>

<template>
  <div>
    <FormEditorField
      help-text="The name of this portfolio."
      :editor-field="model.name"
    >
      <PVInputText
        v-model="model.name.currentValue"
      />
    </FormEditorField>
    <FormEditorField
      help-text="The description of this portfolio - helpful for record keeping, not used for anything."
      :editor-field="model.description"
    >
      <PVTextarea
        v-model="model.description.currentValue"
        auto-resize
      />
    </FormEditorField>
    <FormEditorField
      help-text="When enabled, this portfolio can be accessed by administrators to help with debugging. Only turn this on if you're comfortable with system administrators accessing this data."
      :editor-field="model.adminDebugEnabled"
    >
      <ExplicitInputSwitch
        v-model:value="model.adminDebugEnabled.currentValue"
        on-label="Administrator Debugging Access Enabled"
        off-label="No Administrator Access Enabled"
      />
    </FormEditorField>
  </div>
</template>
