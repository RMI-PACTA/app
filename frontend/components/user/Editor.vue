<script setup lang="ts">
import { type EditorUser } from '@/lib/editor'

const { getMaybeMe } = useSession()
const { t } = useI18n()

const { maybeMe } = await getMaybeMe()

interface Props {
  editorUser: EditorUser
}
interface Emits {
  (e: 'update:editorUser', eu: EditorUser): void
}
const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const prefix = 'UserEditor'
const tt = (key: string) => t(`${prefix}.${key}`)

const eu = computed({
  get: () => props.editorUser,
  set: (eu) => { emit('update:editorUser', eu) },
})
const isMe = computed(() => maybeMe.value?.id === eu.value.id.currentValue)
const profileDescription = computed(() => isMe.value ? tt('your profile') : tt('this user profile'))
const nameHelpText = computed(() => {
  const pd = profileDescription.value
  const pre = tt('The name that will be associated with')
  const post = tt('Note that this name will be accessible to a public audience.')
  return `${pre} ${pd}. ${post}`
})
</script>

<template>
  <div>
    <FormEditorField
      :editor-field="eu.name"
      :help-text="nameHelpText"
    >
      <PVInputText
        v-model="eu.name.currentValue"
      />
    </FormEditorField>
    <FormEditorField
      :editor-field="eu.preferredLanguage"
      help-text="What language should platform freatures that support internationalization default to for this user?"
    >
      <LanguageSelector
        v-model:value="eu.preferredLanguage.currentValue"
      />
    </FormEditorField>
    <FormEditorField
      help-text="If enabled, this user will have administrator priveledges."
      :editor-field="eu.admin"
    >
      <ExplicitInputSwitch
        v-model:value="eu.admin.currentValue"
        on-label="Is an Administrator"
        off-label="Is not an Administrator"
      />
    </FormEditorField>
    <FormEditorField
      help-text="If enabled, this user will have super-administrator priveledges."
      :editor-field="eu.superAdmin"
    >
      <ExplicitInputSwitch
        v-model:value="eu.superAdmin.currentValue"
        on-label="Is a Super Administrator"
        off-label="Is not a Super Administrator"
      />
    </FormEditorField>
  </div>
</template>
