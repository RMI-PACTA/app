<script setup lang="ts">
import {
  type EditorUserFields as EditorFields,
  type EditorUserValues as EditorValues,
} from '@/lib/editor'

const prefix = 'components/user/Editor'

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
      :editor-field="efs.preferredLanguage"
      :editor-value="evs.preferredLanguage"
    >
      <LanguageSelector
        v-model:value="evs.preferredLanguage.currentValue"
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="efs.admin"
      :editor-value="evs.admin"
    >
      <ExplicitInputSwitch
        v-model:value="evs.admin.currentValue"
        :on-label="tt('is admin')"
        :off-label="tt('is not admin')"
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="efs.superAdmin"
      :editor-value="evs.superAdmin"
    >
      <ExplicitInputSwitch
        v-model:value="evs.superAdmin.currentValue"
        :on-label="tt('is super admin')"
        :off-label="tt('is not super admin')"
      />
    </FormEditorField>
  </div>
</template>
