<script setup lang="ts">
import { portfolioEditor } from '@/lib/editor'

const {
  humanReadableTimeFromStandardString,
  humanReadableDateFromStandardString,
} = useTime()
const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()
const i18n = useI18n()

interface EditorObject extends ReturnType<typeof portfolioEditor> {
  id: string
}

const prefix = 'pages/portfolios'
const expandedRows = useState<EditorObject[]>(`${prefix}.expandedRows`, () => [])
const selectedRows = useState<EditorObject[]>(`${prefix}.selectedRows`, () => [])

const [
  { data },
] = await Promise.all([
  useSimpleAsyncData(`${prefix}.portfolios`, () => pactaClient.listPortfolios()),
])

let editorObjects = data.value.items.map((item) => ({ ...portfolioEditor(item, i18n), id: item.id }))

const deletePortfolio = (id: string) => withLoading(
  () => pactaClient.deletePortfolio(id).then(() => {
    editorObjects = editorObjects.filter((editorObject) => editorObject.id !== id)
    expandedRows.value = expandedRows.value.filter((row) => row.id !== id)
  }),
  `${prefix}.deletePortfolio`,
)
const saveChanges = (id: string) => {
  const index = editorObjects.findIndex((editorObject) => editorObject.id === id)
  const eo = presentOrFileBug(editorObjects[index])
  return withLoading(
    () => pactaClient.updatePortfolio(id, eo.changes.value)
      .then(() => pactaClient.findPortfolioById(id))
      .then((portfolio) => {
        editorObjects[index] = { ...portfolioEditor(portfolio, i18n), id }
      }),
    `${prefix}.saveChanges`,
  )
}
</script>

<template>
  <StandardContent>
    <TitleBar title="Portfolios" />
    <p>
      This page shows your logical portfolios.
    </p>
    <PVDataTable
      v-model:selection="selectedRows"
      v-model:expanded-rows="expandedRows"
      :value="editorObjects"
      data-key="id"
      class="portfolio-upload-table"
      size="small"
    >
      <PVColumn selection-mode="multiple" />
      <PVColumn
        field="editorValues.value.createdAt.originalValue"
        header="Created At"
        sortable
      >
        <template #body="slotProps">
          {{ humanReadableTimeFromStandardString(slotProps.data.editorValues.value.createdAt.originalValue).value }}
        </template>
      </PVColumn>
      <PVColumn
        field="editorValues.value.name.originalValue"
        sortable
        header="Name"
      />
      <PVColumn
        expander
        header="Details"
      />
      <template
        #expansion="slotProps"
      >
        <div class="surface-100 p-3">
          <h2 class="mt-0">
            Metadata
          </h2>
          <div class="flex flex-column gap-2 w-fit">
            <div class="flex gap-2 justify-content-between">
              <span>Created At</span>
              <b>{{ humanReadableTimeFromStandardString(slotProps.data.editorValues.value.createdAt.originalValue).value }}</b>
            </div>
            <div class="flex gap-2 justify-content-between">
              <span>Number of Rows</span>
              <b>{{ slotProps.data.editorValues.value.numberOfRows.originalValue }}</b>
            </div>
            <div class="flex gap-2 justify-content-between">
              <span>Holdings Date</span>
              <b>{{ humanReadableDateFromStandardString(slotProps.data.editorValues.value.holdingsDate.originalValue.time).value }}</b>
            </div>
          </div>
          <h2 class="mt-5">
            Editable Properties
          </h2>
          <PortfolioEditor
            v-model:editor-values="slotProps.data.editorValues.value"
            :editor-fields="slotProps.data.editorFields"
          />
          <div class="flex gap-3 justify-content-between">
            <PVButton
              icon="pi pi-trash"
              class="p-button-danger p-button-outlined"
              label="Delete"
              @click="() => deletePortfolio(slotProps.data.id)"
            />
            <div v-tooltip.bottom="slotProps.data.saveTooltip">
              <PVButton
                :disabled="!slotProps.data.canSave.value"
                label="Save Changes"
                icon="pi pi-save"
                icon-pos="right"
                @click="() => saveChanges(slotProps.data.id)"
              />
            </div>
          </div>
          <StandardDebug
            :value="slotProps.data.editorFields.value"
            label="Editor Fields"
          />
          <StandardDebug
            :value="slotProps.data.editorValues.value"
            label="Editor Values"
          />
        </div>
      </template>
    </PVDataTable>
    <div class="flex flex-wrap gap-3 w-full justify-content-between">
      <LinkButton
        class="p-button-outlined"
        icon="pi pi-arrow-left"
        to="/upload"
        label="Upload New Portfolios"
      />
      <!-- TODO(grady) Hook this up to something. -->
      <PVButton
        class="p-button-outlined"
        label="How To Run a Report"
        icon="pi pi-question-circle"
        icon-pos="right"
      />
    </div>
    <StandardDebug
      :value="data"
      label="Editor Objects"
    />
  </standardcontent>
</template>

<style lang="scss">
.portfolio-upload-table.p-datatable.p-datatable-sm {
  width: 100%;

  .p-datatable-row-expansion td {
    padding: 0 0.5rem;

  }

  .p-checkbox {
    width: 1.25rem;
    height: 1.25rem;

    .p-checkbox-box {
      height: 100%;
      width: 100%;
    }
  }
}
</style>
