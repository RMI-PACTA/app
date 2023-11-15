<script setup lang="ts">
import { type EditorUser } from '@/lib/editor'

const { t } = useI18n()
const { getMaybeMe } = await useSession()
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
const profileDescription = computed(() => isMe.value ? tt('you') : tt('this user'))
const nameHelpText = computed(() => {
  const pd = profileDescription.value
  const pre = tt('The name that will be associated with')
  const post = tt('may be public')
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
      :help-text="tt('language help text')"
    >
      <LanguageSelector
        v-model:value="eu.preferredLanguage.currentValue"
      />
    </FormEditorField>
    <FormEditorField
      :help-text="tt('admin help text')"
      :editor-field="eu.admin"
    >
      <ExplicitInputSwitch
        v-model:value="eu.admin.currentValue"
        :on-label="tt('is admin')"
        :off-label="tt('is not admin')"
      />
    </FormEditorField>
    <FormEditorField
      :help-text="tt('super admin help text')"
      :editor-field="eu.superAdmin"
    >
      <ExplicitInputSwitch
        v-model:value="eu.superAdmin.currentValue"
        :on-label="tt('is super admin')"
        :off-label="tt('is not super admin')"
      />
    </FormEditorField>
  </div>
</template>
