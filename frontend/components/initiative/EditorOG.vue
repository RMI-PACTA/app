<script setup lang="ts">
import { type Initiative } from '@/openapi/generated/pacta'

const props = defineProps<{
  initiative: Initiative
}>()

const emit = defineEmits<(e: 'update:initiative', initiative: Initiative) => void>()

const model = computed({
  get: () => props.initiative,
  set: (initiative: Initiative) => { emit('update:initiative', initiative) },
})

const nameCompleted = computed(() => model.value.name.length > 0)
const idRegex = /^[a-zA-Z0-9_-]+$/
const idCompleted = computed(() => idRegex.test(model.value.id))
const publicDescriptionCompleted = computed(() => model.value.publicDescription.length > 0)
const languageCompleted = computed(() => !!model.value.language)
const incompleteFields = computed<string[]>(() => {
  const result: string[] = []
  if (!nameCompleted.value) { result.push('Initiative Name') }
  if (!idCompleted.value) { result.push('Initiative ID') }
  if (!publicDescriptionCompleted.value) { result.push('Public Description') }
  if (!languageCompleted.value) { result.push('Language') }
  return result
})
defineExpose({ incompleteFields })
</script>

<template>
  <div>
    <FormField
      label="Initiative Name"
      help-text="The name of the PACTA initiative."
      required
      :completed="nameCompleted"
    >
      <PVInputText
        v-model="model.name"
      />
    </FormField>
    <FormField
      label="Initiative ID"
      help-text="This is the immutable unique identifier for the initiative. It can only contain alphanumeric characters, underscores, and dashes. This value will be shown in URLs, but will typically not be user visible."
      required
      :completed="idCompleted"
    >
      <PVInputText
        v-model="model.id"
      />
    </FormField>
    <FormField
      label="Affiliation"
      help-text="An optional description of the organization or entity that is hosting this initiative."
    >
      <PVInputText
        v-model="model.affiliation"
      />
    </FormField>
    <FormField
      label="Public Description"
      help-text="The description of the initiative that will be shown to the public. Newlines will be respected."
      required
      :completed="publicDescriptionCompleted"
    >
      <PVTextarea
        v-model="model.publicDescription"
        auto-resize
      />
    </FormField>
    <FormField
      label="Internal Description"
      help-text="The description of the initiative that will be shown to members of the inititiative. Newlines will be respected."
    >
      <PVTextarea
        v-model="model.internalDescription"
        auto-resize
      />
    </FormField>
    <FormField
      label="Participation Mechanism"
      help-text="When disabled, anyone can join this initiative. When enabled, initiative administrators can mint invitation codes that they can share with folks to allow them to join the project."
    >
      <ExplicitInputSwitch
        v-model:value="model.requiresInvitationToJoin"
        on-label="Requires Invitation To Join"
        off-label="Anyone Can Join"
      />
    </FormField>
    <FormField
      label="Open To New Members"
      help-text="When enabled, new members can join the project through the joining mechanism selected above."
    >
      <ExplicitInputSwitch
        v-model:value="model.isAcceptingNewMembers"
        on-label="Accepting New Members"
        off-label="Closed To New Members"
      />
    </FormField>
    <FormField
      label="Open To New Portfolios"
      help-text="When enabled, initiative members can add new portfolios to the initiative."
    >
      <ExplicitInputSwitch
        v-model:value="model.isAcceptingNewPortfolios"
        on-label="Accepting New Portfolios"
        off-label="Closed To New Portfolios"
      />
    </FormField>
    <FormField
      label="Language"
      help-text="What language should reports have when they are generated for this initiative?"
      required
      :completed="languageCompleted"
    >
      <LanguageSelector
        v-model:value="model.language"
      />
    </FormField>
    <FormField
      label="PACTA Version"
      help-text="What version of the PACTA algorithm should this initiative use to generate reports?"
    >
      <PactaversionSelector
        v-model:value="model.pactaVersion"
      />
    </FormField>
  </div>
</template>
