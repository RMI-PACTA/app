<script setup lang="ts">
import { type PactaVersion, type PactaVersionChanges } from '@/openapi/generated/pacta'

const router = useRouter()
const { pactaClient } = useAPI()
const { error: { withLoadingAndErrorHandling, handleOAPIError } } = useModal()
const { fromParams } = useURLParams()

const id = presentOrCheckURL(fromParams('id'))

const prefix = 'admin/pacta-version/[id]'
const persistedPactaVersion = useState<PactaVersion>('totally different???')
const pactaVersion = useState<PactaVersion>(`${prefix}.pactaVersion`)

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
    .then(() => { pactaVersion.value.isDefault = true }),
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
    .then(() => { persistedPactaVersion.value = pactaVersion.value })
    .then(() => router.push('/admin/pacta-version')),
  `${prefix}.saveChanges`
)

// TODO(#13) Remove this from the on-mounted hook
onMounted(async () => {
  await withLoadingAndErrorHandling(
    () => pactaClient.findPactaVersionById(id)
      .then(handleOAPIError)
      .then(pv => {
        pactaVersion.value = { ...pv }
        persistedPactaVersion.value = { ...pv }
      }),
    `${prefix}.getPactaVersion`
  )
})
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
