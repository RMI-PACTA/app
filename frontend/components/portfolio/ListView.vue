<script setup lang="ts">
import { portfolioEditor } from '@/lib/editor'
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

interface EditorObject extends ReturnType<typeof portfolioEditor> {
  id: string
}

const prefix = 'components/portfolio/ListView'
const tt = (s: string) => t(`${prefix}.${s}`)
const expandedRows = useState<EditorObject[]>(`${prefix}.expandedRows`, () => [])
const selectedRows = computed<EditorObject[]>({
  get: () => {
    return editorObjects.filter((editorObject) => selectedPortfolioIDs.value.includes(editorObject.id))
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
    editorObjects = portfolioData.value.items.map((item) => ({ ...portfolioEditor(item, i18n), id: item.id }))
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

let editorObjects = portfolioData.value.items.map((item) => ({ ...portfolioEditor(item, i18n), id: item.id }))

const portfolioGroups = computed<PortfolioGroup[]>(() => portfolioGroupData.value.items)
const selectedPortfolios = computed<Portfolio[]>(() => selectedRows.value.map((row) => row.currentValue.value))

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
  <div class="flex flex-column gap-3">
    <PortfolioGroupMembershipMenuButton
      :selected-portfolios="selectedPortfolios"
      :portfolio-groups="portfolioGroups"
      @changed-memberships="refreshPortfolios"
      @changed-groups="refreshPortfolioGroups"
    />
    <PVDataTable
      v-model:selection="selectedRows"
      v-model:expanded-rows="expandedRows"
      :value="editorObjects"
      data-key="id"
      class="portfolio-upload-table"
      size="small"
      sort-field="editorValues.value.createdAt.originalValue"
      :sort-order="-1"
    >
      <PVColumn selection-mode="multiple" />
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
        field="editorValues.value.name.originalValue"
        sortable
        :header="tt('Name')"
      />
      <PVColumn
        :header="tt('Memberships')"
      >
        <template #body="slotProps">
          <div class="flex flex-column gap-2">
            <div class="flex gap-1 align-items-center flex-wrap">
              <span>{{ tt('Groups') }}:</span>
              <span
                v-for="membership in slotProps.data.editorValues.value.groups.originalValue"
                :key="membership.portfolioGroup.id"
                class="p-tag p-tag-rounded"
              >
                {{ membership.portfolioGroup.name }}
              </span>
            </div>
            <div class="flex gap-2 align-items-center flex-wrap">
              <span>{{ tt('Initiatives') }}:</span>
              <span
                v-for="membership in slotProps.data.editorValues.value.memberships"
                :key="membership"
                class="p-tag p-tag-rounded"
              >
                {{ membership }}
              </span>
            </div>
          </div>
        </template>
      </PVColumn>
      <PVColumn
        expander
        :header="tt('Details')"
      />
      <template
        #expansion="slotProps"
      >
        <div class="surface-100 p-3">
          <h2 class="mt-0">
            {{ tt('Metadata') }}
          </h2>
          <div class="flex flex-column gap-2 w-fit">
            <div class="flex gap-2 justify-content-between">
              <span>{{ tt('Created At') }}</span>
              <b>{{ humanReadableTimeFromStandardString(slotProps.data.editorValues.value.createdAt.originalValue).value }}</b>
            </div>
            <div class="flex gap-2 justify-content-between">
              <span>{{ tt('Number of Rows') }}</span>
              <b>{{ slotProps.data.editorValues.value.numberOfRows.originalValue }}</b>
            </div>
            <div class="flex gap-2 justify-content-between">
              <span>{{ tt('Holdings Date') }}</span>
              <b>{{ humanReadableDateFromStandardString(slotProps.data.editorValues.value.holdingsDate.originalValue.time).value }}</b>
            </div>
          </div>
          <h2 class="mt-5">
            {{ tt('Editable Properties') }}
          </h2>
          <PortfolioEditor
            v-model:editor-values="slotProps.data.editorValues.value"
            :editor-fields="slotProps.data.editorFields.value"
          />
          <div class="flex gap-3 justify-content-between">
            <PVButton
              icon="pi pi-trash"
              class="p-button-danger p-button-outlined"
              :label="tt('Delete')"
              @click="() => deletePortfolio(slotProps.data.id)"
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
          <!--
          <StandardDebug
            :value="slotProps.data.editorFields.value"
            label="Editor Fields"
          />
          <StandardDebug
            :value="slotProps.data.editorValues.value"
            label="Editor Values"
          />
-->
        </div>
      </template>
    </PVDataTable>
    <div class="flex flex-wrap gap-3 w-full justify-content-between">
      <LinkButton
        class="p-button-outlined"
        icon="pi pi-arrow-left"
        to="/upload"
        :label="tt('Upload New Portfolios')"
      />
      <!-- TODO(grady) Hook this up to something. -->
      <PVButton
        class="p-button-outlined"
        :label="tt('How To Run a Report')"
        icon="pi pi-question-circle"
        icon-pos="right"
      />
    </div>
    <!--
    <StandardDebug
      :value="portfolioData"
      label="Portfolio Data"
    />
    -->
    <StandardDebug
      :value="selectedPortfolios"
      label="Selected Portfolios"
    /><!--
    <StandardDebug
      :value="portfolioGroupData"
      label="PG Data"
    />
    -->
  </div>
</template>
