<script setup lang="ts">
import { type EditorInitiative } from '@/lib/editor'

interface Props {
  editorInitiative: EditorInitiative
}
interface Emits {
  (e: 'update:editorInitiative', ei: EditorInitiative): void
}
const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const model = computed({
  get: () => props.editorInitiative,
  set: (editorInitiative: EditorInitiative) => { emit('update:editorInitiative', editorInitiative) },
})
</script>

<template>
  <div>
    <FormEditorField
      help-text="The name of the PACTA initiative."
      :editor-field="model.name"
    >
      <PVInputText
        v-model="model.name.currentValue"
      />
    </FormEditorField>
    <FormEditorField
      help-text="This is the immutable unique identifier for the initiative. It can only contain alphanumeric characters, underscores, and dashes. This value will be shown in URLs, but will typically not be user visible."
      :editor-field="model.id"
    >
      <PVInputText
        v-model="model.id.currentValue"
        :disabled="!!model.id.originalValue"
      />
    </FormEditorField>
    <FormEditorField
      help-text="An optional description of the organization or entity that is hosting this initiative."
      :editor-field="model.affiliation"
    >
      <PVInputText
        v-model="model.affiliation.currentValue"
      />
    </FormEditorField>
    <FormEditorField
      help-text="The description of the initiative that will be shown to the public. Newlines will be respected."
      :editor-field="model.publicDescription"
    >
      <PVTextarea
        v-model="model.publicDescription.currentValue"
        auto-resize
      />
    </FormEditorField>
    <FormEditorField
      help-text="The description of the initiative that will be shown to members of the inititiative. Newlines will be respected."
      :editor-field="model.internalDescription"
    >
      <PVTextarea
        v-model="model.internalDescription.currentValue"
        auto-resize
      />
    </FormEditorField>
    <FormEditorField
      help-text="When disabled, anyone can join this initiative. When enabled, initiative administrators can mint invitation codes that they can share with folks to allow them to join the project."
      :editor-field="model.requiresInvitationToJoin"
    >
      <ExplicitInputSwitch
        v-model:value="model.requiresInvitationToJoin.currentValue"
        on-label="Requires Invitation To Join"
        off-label="Anyone Can Join"
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="model.isAcceptingNewMembers"
      help-text="When enabled, new members can join the project through the joining mechanism selected above."
    >
      <ExplicitInputSwitch
        v-model:value="model.isAcceptingNewMembers.currentValue"
        on-label="Accepting New Members"
        off-label="Closed To New Members"
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="model.isAcceptingNewPortfolios"
      help-text="When enabled, initiative members can add new portfolios to the initiative."
    >
      <ExplicitInputSwitch
        v-model:value="model.isAcceptingNewPortfolios.currentValue"
        on-label="Accepting New Portfolios"
        off-label="Closed To New Portfolios"
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="model.language"
      help-text="What language should reports have when they are generated for this initiative?"
    >
      <LanguageSelector
        v-model:value="model.language.currentValue"
      />
    </FormEditorField>
    <FormEditorField
      help-text="What version of the PACTA algorithm should this initiative use to generate reports?"
      :editor-field="model.pactaVersion"
    >
      <PactaversionSelector
        v-model:value="model.pactaVersion.currentValue"
      />
    </FormEditorField>
  </div>
</template>
