<script setup lang="ts">
import { type PactaVersion, type PactaVersionChanges } from '@/openapi/generated/pacta'

const router = useRouter()
const { pactaClient } = useAPI()
const { error: { withLoadingAndErrorHandling, handleOAPIError } } = useModal()
const { fromParams } = useURLParams()

const id = presentOrCheckURL(fromParams('id'))

const prefix = `admin/pacta-version/${id}`
const pactaVersion = useState<PactaVersion>(`${prefix}.pactaVersion`)
const { data: persistedPactaVersion, error, refresh } = await useAsyncData(`${prefix}.getPactaVersion`, () => {
  return withLoadingAndErrorHandling(() => {
    return pactaClient.findPactaVersionById(id)
      .then(handleOAPIError)
  }, `${prefix}.getPactaVersion`)
})
if (error.value) {
  throw createError(error.value)
}
if (!persistedPactaVersion.value) {
  throw createError({ message: 'PACTA version not found' })
}
pactaVersion.value = { ...persistedPactaVersion.value }
const refreshPACTA = async () => {
  await refresh()
  if (persistedPactaVersion.value) {
    pactaVersion.value = { ...persistedPactaVersion.value }
  }
}

const changes = computed<PactaVersionChanges>(() => {
  const a = persistedPactaVersion.value
  const b = pactaVersion.value
  if (!a || !b) { return {} }
  return {
    ...(a.name !== b.name ? { name: b.name } : {}),
    ...(a.description !== b.description ? { description: b.description } : {}),
    ...(a.digest !== b.digest ? { digest: b.digest } : {})
  }
})
const hasChanges = computed<boolean>(() => Object.keys(changes.value).length > 0)

const markDefault = () => withLoadingAndErrorHandling(
  () => pactaClient.markPactaVersionAsDefault(id)
    .then(handleOAPIError)
    .then(refreshPACTA),
  `${prefix}.markPactaVersionAsDefault`
)
const deletePV = () => withLoadingAndErrorHandling(
  () => pactaClient.deletePactaVersion(id)
    .then(handleOAPIError)
    .then(() => router.push('/admin/pacta-version')),
  `${prefix}.deletePactaVersion`
)
const saveChanges = () => withLoadingAndErrorHandling(
  () => pactaClient.updatePactaVersion(id, changes.value)
    .then(handleOAPIError)
    .then(refreshPACTA)
    .then(() => router.push('/admin/pacta-version')),
  `${prefix}.saveChanges`
)
</script>

<template>
  <StandardContent v-if="pactaVersion">
    <TitleBar :title="`Editing PACTA Version: ${pactaVersion.name}`" />
    <div class="flex gap-3">
      <PVButton
        :disabled="pactaVersion.isDefault"
        class="p-button-success"
        :label="pactaVersion.isDefault ? 'Default Version' : 'Make Default Version'"
        :icon="pactaVersion.isDefault ? 'pi pi-check-circle' : 'pi pi-circle'"
        @click="markDefault"
      />
      <PVButton
        :disabled="pactaVersion.isDefault"
        icon="pi pi-trash"
        class="p-button-danger"
        label="Delete"
        @click="deletePV"
      />
    </div>
    <PactaversionEditor
      v-model:pactaVersion="pactaVersion"
    />
    <div class="flex gap-3">
      <LinkButton
        label="Discard Changes"
        icon="pi pi-arrow-left"
        class="p-button-secondary p-button-outlined"
        to="/admin/pacta-version"
      />
      <PVButton
        :disabled="!hasChanges"
        label="Save Changes"
        icon="pi pi-arrow-right"
        icon-pos="right"
        @click="saveChanges"
      />
    </div>
    <StandardDebug
      :value="persistedPactaVersion"
      label="Persisted PACTA Version"
    />
    <StandardDebug
      :value="pactaVersion"
      label="PACTA Version"
    />
    <StandardDebug
      :value="changes"
      label="PV Changes"
    />
  </StandardContent>
</template>
