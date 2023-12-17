<script setup lang="ts">
import { pactaVersionEditor } from '@/lib/editor'

const router = useRouter()
const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()
const { fromParams } = useURLParams()
const localePath = useLocalePath()
const i18n = useI18n()
const { t } = i18n

const id = presentOrCheckURL(fromParams('id'))

const tt = (key: string) => t(`pages/admin/pacta-version/id/${key}`)
const prefix = `admin/pacta-version/${id}`
const { data, refresh } = await useSimpleAsyncData(
  `${prefix}.getPactaVersion`,
  () => pactaClient.findPactaVersionById(id),
)

const {
  setEditorValue: setPactaVersion,
  editorFields,
  editorValues,
  changes,
  saveTooltip,
  canSave,
} = pactaVersionEditor(presentOrCheckURL(data.value, 'no PACTA version in response'), i18n)
const isDefault = computed(() => editorValues.value.isDefault.currentValue)

const refreshPACTA = async () => {
  await refresh()
  setPactaVersion(presentOrCheckURL(data.value, 'no PACTA version in response after refresh'))
}

const markDefault = () => withLoading(
  () => pactaClient.markPactaVersionAsDefault(id)
    .then(refreshPACTA),
  `${prefix}.markPactaVersionAsDefault`,
)
const deletePV = () => withLoading(
  () => pactaClient.deletePactaVersion(id)
    .then(() => router.push(localePath('/admin/pacta-version'))),
  `${prefix}.deletePactaVersion`,
)
const saveChanges = () => withLoading(
  () => pactaClient.updatePactaVersion(id, changes.value)
    .then(refreshPACTA)
    .then(() => router.push(localePath('/admin/pacta-version'))),
  `${prefix}.saveChanges`,
)
</script>

<template>
  <StandardContent v-if="editorValues">
    <TitleBar :title="`${tt('Editing PACTA Version')}: ${editorValues.name.currentValue}`" />
    <div class="flex gap-3">
      <PVButton
        :disabled="isDefault"
        class="p-button-success"
        :label="isDefault ? tt('Default Version') : tt('Make Default Version')"
        :icon="isDefault ? 'pi pi-check-circle' : 'pi pi-circle'"
        @click="markDefault"
      />
      <PVButton
        :disabled="isDefault"
        icon="pi pi-trash"
        class="p-button-danger"
        :label="tt('Delete')"
        @click="deletePV"
      />
    </div>
    <PactaversionEditor
      v-model:editorValues="editorValues"
      :editor-fields="editorFields"
    />
    <div class="flex gap-3 align-items-center">
      <LinkButton
        :label="tt('Discard Changes')"
        icon="pi pi-arrow-left"
        class="p-button-secondary p-button-outlined"
        :to="localePath('/admin/pacta-version')"
      />
      <div v-tooltip.bottom="saveTooltip">
        <PVButton
          :disabled="!canSave"
          :label="tt('Save Changes')"
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
      label="PV Changes"
    />
  </StandardContent>
</template>
