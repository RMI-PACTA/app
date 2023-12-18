<script setup lang="ts">
import { portfolioGroupEditor } from '@/lib/editor'
import { type Portfolio, type PortfolioGroup } from '@/openapi/generated/pacta'

const {
  humanReadableTimeFromStandardString,
  humanReadableDateFromStandardString,
} = useTime()
const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()
const i18n = useI18n()
const { t } = i18n

interface Props {
  selectedPortfolioIds: string[]
  selectedPortfolioGroupIds: string[]
}
const props = defineProps<Props>()
interface Emits {
  (e: 'update:selectedPortfolioIds', value: string[]): void
  (e: 'update:selectedPortfolioGroupIds', value: string[]): void
}
const emit = defineEmits<Emits>()

const selectedPortfolioIDs = computed({
  get: () => props.selectedPortfolioIds ?? [],
  set: (value: string[]) => { emit('update:selectedPortfolioIds', value) },
})
const selectedPortfolioGroupIDs = computed({
  get: () => props.selectedPortfolioIds ?? [],
  set: (value: string[]) => { emit('update:selectedPortfolioGroupIds', value) },
})

interface EditorObject extends ReturnType<typeof portfolioGroupEditor> {
  id: string
}

const prefix = 'components/portfolio/group/ListView'
const tt = (s: string) => t(`${prefix}.${s}`)
const expandedRows = useState<EditorObject[]>(`${prefix}.expandedRows`, () => [])
const selectedRows = computed<EditorObject[]>({
  get: () => {
    return editorObjects.filter((editorObject) => selectedPortfolioGroupIDs.value.includes(editorObject.id))
  },
  set: (value: EditorObject[]) => {
    selectedPortfolioGroupIDs.value = value.map((row) => row.id)
  },
})

const [
  { data: portfolioData, refresh: refreshPortfoliosApi },
  { data: portfolioGroupData, refresh: refreshPortfolioGroupsApi },
] = await Promise.all([
  useSimpleAsyncData(`${prefix}.portfolios`, () => pactaClient.listPortfolios()),
  useSimpleAsyncData(`${prefix}.portfolioGroups`, () => pactaClient.listPortfolioGroups()),
])
const refreshPortfolios = async () => {
  await withLoading(async () => {
    await refreshPortfoliosApi()
    editorObjects = portfolioData.value.items.map((item) => ({ ...portfolioGroupEditor(item, i18n), id: item.id }))
  }, `${prefix}.refreshPortfolios`)
}
const refreshPortfolioGroups = async () => {
  await withLoading(async () => {
    await refreshPortfolioGroupsApi()
  }, `${prefix}.refreshPortfolioGroups`)
}
const refreshAll = async () => {
  await Promise.all([
    refreshPortfolios(),
    refreshPortfolioGroups(),
  ])
}

let editorObjects = portfolioGroupData.value.items.map((item) => ({ ...portfolioGroupEditor(item, i18n), id: item.id }))

const deletePortfolioGroup = (id: string) => withLoading(
  () => pactaClient.deletePortfolioGroup(id).then(() => {
    editorObjects = editorObjects.filter((editorObject) => editorObject.id !== id)
    expandedRows.value = expandedRows.value.filter((row) => row.id !== id)
  }),
  `${prefix}.deletePortfolioGroup`,
)
const saveChanges = (id: string) => {
  const index = editorObjects.findIndex((editorObject) => editorObject.id === id)
  const eo = presentOrFileBug(editorObjects[index])
  return withLoading(
    () => pactaClient.updatePortfolioGroup(id, eo.changes.value)
      .then(() => pactaClient.findPortfolioGroupById(id))
      .then((portfolio) => {
        editorObjects[index] = { ...portfolioGroupEditor(portfolio, i18n), id }
      }),
    `${prefix}.saveChanges`,
  )
}
</script>

<template>
  <div class="flex flex-column gap-3">
    <PVDataTable
      v-model:selection="selectedRows"
      v-model:expanded-rows="expandedRows"
      :value="editorObjects"
      data-key="id"
      class="portfolio-group-upload-table"
      size="small"
      sort-field="editorValues.value.createdAt.originalValue"
      :sort-order="-1"
    >
      <PVColumn selection-mode="multiple" />
      <PVColumn
        field="editorValues.value.name.originalValue"
        sortable
        :header="tt('Name')"
      />
      <PVColumn
        field="editorValues.value.createdAt.originalValue"
        :header="tt('Created At')"
        sortable
      >
        <template #body="slotProps">
          {{ humanReadableTimeFromStandardString(slotProps.data.editorValues.value.createdAt.originalValue).value }}
        </template>
      </PVColumn>
      <PVColumn
        :header="tt('Number of Members')"
      >
        <template #body="slotProps">
          {{ slotProps.data.editorValues.value.members.originalValue.length }}
        </template>
      </PVColumn>
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
          </div>
          <h2 class="mt-5">
            Editable Properties
          </h2>
          <PortfolioGroupEditor
            v-model:editor-values="slotProps.data.editorValues.value"
            :editor-fields="slotProps.data.editorFields.value"
          />
          <div class="flex gap-3 justify-content-between">
            <PVButton
              icon="pi pi-trash"
              class="p-button-danger p-button-outlined"
              :label="tt('Delete')"
              @click="() => deletePortfolioGroup(slotProps.data.id)"
            />
            <div v-tooltip.bottom="slotProps.data.saveTooltip">
              <PVButton
                :disabled="!slotProps.data.canSave.value"
                :label="tt('Save Changes')"
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
    <StandardDebug
      :value="selectedRows"
      label="Selected Rows"
    />
  </div>
</template>
