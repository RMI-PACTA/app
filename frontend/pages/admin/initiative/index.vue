<script setup lang="ts">
const { humanReadableTimeFromStandardString } = useTime()
const router = useRouter()
const pactaClient = await usePACTA()
const { loading: { withLoading } } = useModal()
const localePath = useLocalePath()

const prefix = 'admin/initiative'
const { data: initiatives, refresh } = await useSimpleAsyncData(`${prefix}.get`, () => pactaClient.listInitiatives())

const newInitiative = () => router.push(localePath('/admin/initiative/new'))
const deleteInitiative = (id: string) => withLoading(
  () => pactaClient.deleteInitiative(id)
    .then(() => refresh()),
  `${prefix}.delete`,
)
</script>

<template>
  <StandardContent>
    <TitleBar title="Initiatives" />
    <p>
      TODO(#38) Add I18n
      General information about initiatives here.
    </p>
    <PVDataTable
      :value="initiatives"
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
      >
        <template #body="slotProps">
          {{ humanReadableTimeFromStandardString(slotProps.data.createdAt) }}
        </template>
      </PVColumn>
      <PVColumn header="View">
        <template #body="slotProps">
          <LinkButton
            icon="pi pi-arrow-right"
            :to="localePath(`/initiative/${slotProps.data.id}`)"
          />
        </template>
      </PVColumn>
      <PVColumn header="Delete">
        <template #body="slotProps">
          <PVButton
            icon="pi pi-trash"
            class="p-button-danger p-button-outlined"
            @click="() => deleteInitiative(slotProps.data.id)"
          />
        </template>
      </PVColumn>
    </PVDataTable>
    <PVButton
      label="New Initiative"
      icon="pi pi-plus"
      @click="newInitiative"
    />
    <StandardDebug :value="initiatives" />
  </StandardContent>
</template>
