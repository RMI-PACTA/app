<script setup lang="ts">
import { userEditor } from '@/lib/editor'

const pactaClient = usePACTA()
const { fromParams } = useURLParams()
const { loading: { withLoading } } = useModal()
const router = useRouter()
const localePath = useLocalePath()
const i18n = useI18n()

const id = presentOrCheckURL(fromParams('id'))
const prefix = `user/[${id}]`

const { data } = await useSimpleAsyncData(`${prefix}.getUser`, () => pactaClient.findUserById(id))
const {
  editorValues,
  editorFields,
  changes,
  saveTooltip,
  canSave,
} = userEditor(presentOrCheckURL(data.value, 'no user in response'), i18n)

const deleteUser = () => withLoading(
  () => pactaClient.deleteUser(id).then(() => router.push(localePath('/'))),
  `${prefix}.deleteUser`,
)
const saveChanges = () => withLoading(
  () => pactaClient.updateUser(id, changes.value)
    .then(() => router.push(localePath(`/user/${id}`))),
  `${prefix}.saveChanges`,
)
</script>

<template>
  <div class="flex flex-column gap-3">
    <UserEditor
      v-model:editorValues="editorValues"
      :editor-fields="editorFields"
    />
    <div class="flex gap-3">
      <PVButton
        icon="pi pi-trash"
        class="p-button-danger"
        label="Delete"
        @click="deleteUser"
      />
      <LinkButton
        label="Discard Changes"
        icon="pi pi-arrow-left"
        class="p-button-secondary p-button-outlined"
        :to="localePath(`/user/${id}`)"
      />
      <div v-tooltip.bottom="saveTooltip">
        <PVButton
          :disabled="!canSave"
          label="Save Changes"
          icon="pi pi-arrow-right"
          icon-pos="right"
          @click="saveChanges"
        />
      </div>
    </div>
    <StandardDebug
      :value="editorFields"
      label="Editor Fields"
    />
    <StandardDebug
      :value="editorValues"
      label="Editor Values"
    />
    <StandardDebug
      :value="changes"
      label="Changes"
    />
  </div>
</template>
