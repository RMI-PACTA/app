<script setup lang="ts">
import {
  type EditorInitiativeFields as EditorFields,
  type EditorInitiativeValues as EditorValues,
} from '@/lib/editor'

const prefix = 'components/initiative/Editor'
const { t } = useI18n()
const tt = (key: string) => t(`${prefix}.${key}`)

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
      :editor-field="efs.id"
      :editor-value="evs.id"
    >
      <PVInputText
        v-model="evs.id.currentValue"
        :disabled="!!evs.id.originalValue"
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="efs.affiliation"
      :editor-value="evs.affiliation"
    >
      <PVInputText
        v-model="evs.affiliation.currentValue"
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="efs.publicDescription"
      :editor-value="evs.publicDescription"
    >
      <PVTextarea
        v-model="evs.publicDescription.currentValue"
        auto-resize
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="efs.internalDescription"
      :editor-value="evs.internalDescription"
    >
      <PVTextarea
        v-model="evs.internalDescription.currentValue"
        auto-resize
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="efs.requiresInvitationToJoin"
      :editor-value="evs.requiresInvitationToJoin"
    >
      <ExplicitInputSwitch
        v-model:value="evs.requiresInvitationToJoin.currentValue"
        :on-label="tt('Requires Invitation To Join')"
        :off-label="tt('Anyone Can Join')"
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="efs.isAcceptingNewMembers"
      :editor-value="evs.isAcceptingNewMembers"
    >
      <ExplicitInputSwitch
        v-model:value="evs.isAcceptingNewMembers.currentValue"
        :on-label="tt('Accepting New Members')"
        :off-label="tt('Closed To New Members')"
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="efs.isAcceptingNewPortfolios"
      :editor-value="evs.isAcceptingNewPortfolios"
    >
      <ExplicitInputSwitch
        v-model:value="evs.isAcceptingNewPortfolios.currentValue"
        :on-label="tt('Accepting New Portfolios')"
        :off-label="tt('Closed To New Portfolios')"
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="efs.language"
      :editor-value="evs.language"
    >
      <LanguageSelector
        v-model:value="evs.language.currentValue"
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="efs.pactaVersion"
      :editor-value="evs.pactaVersion"
    >
      <PactaversionSelector
        v-model:value="evs.pactaVersion.currentValue"
      />
    </FormEditorField>
  </div>
</template>
