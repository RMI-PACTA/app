<script setup lang="ts">
import { type PactaVersion } from '@/openapi/generated/pacta'

const router = useRouter()
const { pactaClient } = useAPI()
const { error: { withLoadingAndErrorHandling, handleOAPIError } } = useModal()

const prefix = 'admin/pacta-version'
const pactaVersions = useState<PactaVersion[]>(`${prefix}.pactaVersions`, () => [])

const newPV = () => router.push('/admin/pacta-version/new')
const markDefault = (id: string) => withLoadingAndErrorHandling(
  () => pactaClient.markPactaVersionAsDefault(id)
    .then(handleOAPIError)
    .then(() => { pactaVersions.value = pactaVersions.value.map(pv => ({ ...pv, isDefault: id === pv.id })) }),
  `${prefix}.markPactaVersionAsDefault`
)
const deletePV = (id: string) => withLoadingAndErrorHandling(
  () => pactaClient.deletePactaVersion(id)
    .then(handleOAPIError)
    .then(() => { pactaVersions.value = pactaVersions.value.filter(pv => pv.id !== id) }),
  `${prefix}.deletePactaVersion`
)

// TODO(#13) Remove this from the on-mounted hook
onMounted(async () => {
  await withLoadingAndErrorHandling(
    () => pactaClient.listPactaVersions()
      .then(handleOAPIError)
      .then(pvs => { pactaVersions.value = pvs }),
    `${prefix}.getPactaVersions`
  )
})
</script>

<template>
  <StandardContent>
    <TitleBar title="PACTA Versions" />
    <p>
      General ideas about PACTA versions go here.
    </p>
    <PVDataTable
      :value="pactaVersions"
      class="w-full"
      sort-field="createdAt"
      :sort-order="-1"
    >
      <PVColumn
        field="name"
        header="Name"
        sortable
      />
      <PVColumn
        header="Created At"
        field="createdAt"
        data-type="date"
        sortable
      />
      <PVColumn header="Edit">
        <template #body="slotProps">
          <LinkButton
            icon="pi pi-arrow-right"
            :to="`/admin/pacta-version/${slotProps.data.id}`"
          />
        </template>
      </PVColumn>
      <PVColumn header="Default">
        <template #body="slotProps">
          <PVButton
            :icon="slotProps.data.isDefault ? 'pi pi-check-circle' : 'pi pi-circle'"
            class="p-button-success"
            :disabled="slotProps.data.isDefault"
            @click="() => markDefault(slotProps.data.id)"
          />
        </template>
      </PVColumn>
      <PVColumn header="Delete">
        <template #body="slotProps">
          <PVButton
            icon="pi pi-trash"
            class="p-button-danger p-button-outlined"
            @click="() => deletePV(slotProps.data.id)"
          />
        </template>
      </PVColumn>
    </PVDataTable>
    <PVButton
      label="New PACTA Version"
      icon="pi pi-plus"
      @click="newPV"
    />
    <StandardDebug :value="pactaVersions" />
  </StandardContent>
</template>
