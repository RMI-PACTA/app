<script setup lang="ts">
import { type EditorInitiative, isComplete } from '@/lib/editor'

const props = defineProps<{
  editorInitiative: EditorInitiative
}>()

const emit = defineEmits<(e: 'update:editorInitiative', ei: EditorInitiative) => void>()

const model = computed({
  get: () => props.editorInitiative,
  set: (editorInitiative: EditorInitiative) => { emit('update:editorInitiative', editorInitiative) },
})
</script>

<template>
  <div>
    <FormField
      label="Initiative Name"
      help-text="The name of the PACTA initiative."
      :required="model.name.isRequired"
      :completed="isComplete(model.name)"
    >
      <PVInputText
        v-model="model.name.currentValue"
      />
    </FormField>
    <FormField
      label="Initiative ID"
      help-text="This is the immutable unique identifier for the initiative. It can only contain alphanumeric characters, underscores, and dashes. This value will be shown in URLs, but will typically not be user visible."
      :required="model.id.isRequired"
      :completed="isComplete(model.id)"
    >
      <PVInputText
        v-model="model.id.currentValue"
        :disabled="!!model.id.originalValue"
      />
    </FormField>
    <FormField
      label="Affiliation"
      help-text="An optional description of the organization or entity that is hosting this initiative."
    >
      <PVInputText
        v-model="model.affiliation.currentValue"
      />
    </FormField>
    <FormField
      label="Public Description"
      help-text="The description of the initiative that will be shown to the public. Newlines will be respected."
      :required="model.publicDescription.isRequired"
      :completed="isComplete(model.publicDescription)"
    >
      <PVTextarea
        v-model="model.publicDescription.currentValue"
        auto-resize
      />
    </FormField>
    <FormField
      label="Internal Description"
      help-text="The description of the initiative that will be shown to members of the inititiative. Newlines will be respected."
    >
      <PVTextarea
        v-model="model.internalDescription.currentValue"
        auto-resize
      />
    </FormField>
    <FormField
      label="Participation Mechanism"
      help-text="When disabled, anyone can join this initiative. When enabled, initiative administrators can mint invitation codes that they can share with folks to allow them to join the project."
    >
      <ExplicitInputSwitch
        v-model:value="model.requiresInvitationToJoin.currentValue"
        on-label="Requires Invitation To Join"
        off-label="Anyone Can Join"
      />
    </FormField>
    <FormField
      label="Open To New Members"
      help-text="When enabled, new members can join the project through the joining mechanism selected above."
    >
      <ExplicitInputSwitch
        v-model:value="model.isAcceptingNewMembers.currentValue"
        on-label="Accepting New Members"
        off-label="Closed To New Members"
      />
    </FormField>
    <FormField
      label="Open To New Portfolios"
      help-text="When enabled, initiative members can add new portfolios to the initiative."
    >
      <ExplicitInputSwitch
        v-model:value="model.isAcceptingNewPortfolios.currentValue"
        on-label="Accepting New Portfolios"
        off-label="Closed To New Portfolios"
      />
    </FormField>
    <FormField
      ref="fields"
      label="Language"
      help-text="What language should reports have when they are generated for this initiative?"
      :required="model.language.isRequired"
      :completed="isComplete(model.language)"
    >
      <LanguageSelector
        v-model:value="model.language.currentValue"
      />
    </FormField>
    <FormField
      label="PACTA Version"
      help-text="What version of the PACTA algorithm should this initiative use to generate reports?"
      :required="model.pactaVersion.isRequired"
      :completed="isComplete(model.pactaVersion)"
    >
      <PactaversionSelector
        v-model:value="model.pactaVersion.currentValue"
      />
    </FormField>
  </div>
</template>
