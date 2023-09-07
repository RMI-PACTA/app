<script setup lang="ts">
import { type PactaVersion } from '@/openapi/generated/pacta'

const router = useRouter()
const { pactaClient } = useAPI()
const { error: { withLoadingAndErrorHandling, handleOAPIError } } = useModal()

const prefix = 'admin/pacta-version'
const pactaVersions = useState<PactaVersion[]>(`${prefix}.pactaVersions`, () => [])

const deletePV = (id: string) => withLoadingAndErrorHandling(
  () => pactaClient.deletePactaVersion(id)
    .then(handleOAPIError)
    .then(() => { pactaVersions.value = pactaVersions.value.filter(pv => pv.id !== id) }),
  `${prefix}.deletePactaVersion`
)
const newPV = () => router.push('/admin/pacta-version/new')
</script>

<template>
  <StandardContent>
    <TitleBar title="PACTA Versions" />
    <p>
      General ideas about pacta versions go here.
    </p>
    <PVDataTable
      :value="pactaVersions"
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
      <PVColumn>
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
  </StandardContent>
</template>
