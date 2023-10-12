<script setup lang="ts">
import { pactaVersionEditor } from '@/lib/editor'

const router = useRouter()
const pactaClient = await usePACTA()
const { loading: { withLoading } } = useModal()
const { fromParams } = useURLParams()

const id = presentOrCheckURL(fromParams('id'))

const prefix = `admin/pacta-version/${id}`
const { data, refresh } = await useSimpleAsyncData(
  `${prefix}.getPactaVersion`,
  () => pactaClient.findPactaVersionById(id),
)

const {
  setEditorValue: setPactaVersion,
  editorObject: editorPactaVersion,
  changes,
  saveTooltip,
  canSave,
} = pactaVersionEditor(presentOrCheckURL(data.value, 'no PACTA version in response'))
const isDefault = computed(() => editorPactaVersion.value.isDefault.currentValue)

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
    .then(() => router.push('/admin/pacta-version')),
  `${prefix}.deletePactaVersion`,
)
const saveChanges = () => withLoading(
  () => pactaClient.updatePactaVersion(id, changes.value)
    .then(refreshPACTA)
    .then(() => router.push('/admin/pacta-version')),
  `${prefix}.saveChanges`,
)
</script>

<template>
  <StandardContent v-if="editorPactaVersion">
    <TitleBar :title="`Editing PACTA Version: ${editorPactaVersion.name.currentValue}`" />
    <div class="flex gap-3">
      <PVButton
        :disabled="isDefault"
        class="p-button-success"
        :label="isDefault ? 'Default Version' : 'Make Default Version'"
        :icon="isDefault ? 'pi pi-check-circle' : 'pi pi-circle'"
        @click="markDefault"
      />
      <PVButton
        :disabled="isDefault"
        icon="pi pi-trash"
        class="p-button-danger"
        label="Delete"
        @click="deletePV"
      />
    </div>
    <PactaversionEditor
      v-model:editorPactaVersion="editorPactaVersion"
    />
    <div class="flex gap-3 align-items-center">
      <LinkButton
        label="Discard Changes"
        icon="pi pi-arrow-left"
        class="p-button-secondary p-button-outlined"
        to="/admin/pacta-version"
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
      :value="editorPactaVersion"
      label="PACTA Version"
    />
    <StandardDebug
      :value="changes"
      label="PV Changes"
    />
  </StandardContent>
</template>
