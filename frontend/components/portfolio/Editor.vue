<script setup lang="ts">
import {
  type EditorPortfolioFields as EditorFields,
  type EditorPortfolioValues as EditorValues,
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
    <FormEditorField
      :editor-field="efs.propertyHoldingsDate"
      :editor-value="evs.propertyHoldingsDate"
    >
      <InputsHoldingsDate
        v-model:value="evs.propertyHoldingsDate.currentValue"
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="efs.propertyESG"
      :editor-value="evs.propertyESG"
    >
      <InputsEsg
        v-model:value="evs.propertyESG.currentValue"
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="efs.propertyExternal"
      :editor-value="evs.propertyExternal"
    >
      <InputsExternal
        v-model:value="evs.propertyExternal.currentValue"
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="efs.propertyEngagementStrategy"
      :editor-value="evs.propertyEngagementStrategy"
    >
      <InputsEngagementStrategy
        v-model:value="evs.propertyEngagementStrategy.currentValue"
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="efs.adminDebugEnabled"
      :editor-value="evs.adminDebugEnabled"
    >
      <AdminDebugEnabledToggleButton
        v-model:value="evs.adminDebugEnabled.currentValue"
      />
    </FormEditorField>
  </div>
</template>
