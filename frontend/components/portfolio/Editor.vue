<script setup lang="ts">
import {
  type EditorPortfolioFields as EditorFields,
  type EditorPortfolioValues as EditorValues,
} from '@/lib/editor'

const prefix = 'components/portfolio/Editor'

const { t } = useI18n()
const tt = (key: string) => t(`${prefix}.${key}`)

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
      :editor-field="efs.adminDebugEnabled"
      :editor-value="evs.adminDebugEnabled"
    >
      <ExplicitInputSwitch
        v-model:value="evs.adminDebugEnabled.currentValue"
        :on-label="tt('Administrator Debugging Access Enabled')"
        :off-label="tt('No Administrator Access Enabled')"
      />
    </FormEditorField>
  </div>
</template>
