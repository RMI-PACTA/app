<script setup lang="ts">
import { analysisEditor } from '@/lib/editor'
import { type Portfolio, type PortfolioGroup, type Initiative, type Analysis } from '@/openapi/generated/pacta'
import { selectedCountSuffix } from '@/lib/selection'

const { public: { apiServerURL } } = useRuntimeConfig()
const {
  humanReadableTimeFromStandardString,
} = useTime()
const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()
const i18n = useI18n()
const { t } = i18n

interface Props {
  portfolios: Portfolio[]
  portfolioGroups: PortfolioGroup[]
  initiatives: Initiative[]
  analyses: Analysis[]
  selectedPortfolioIds: string[]
  selectedPortfolioGroupIds: string[]
  selectedAnalysisIds: string[]
}
const props = defineProps<Props>()
interface Emits {
  (e: 'update:selectedPortfolioIds', value: string[]): void
  (e: 'update:selectedPortfolioGroupIds', value: string[]): void
  (e: 'update:selectedAnalysisIds', value: string[]): void
  (e: 'refresh'): void
}
const emit = defineEmits<Emits>()

const refresh = () => { emit('refresh') }

const selectedAnalysisIDs = computed({
  get: () => props.selectedAnalysisIds ?? [],
  set: (value: string[]) => { emit('update:selectedAnalysisIds', value) },
})

interface EditorObject extends ReturnType<typeof analysisEditor> {
  id: string
}

const prefix = 'components/analysis/ListView'
const tt = (s: string) => t(`${prefix}.${s}`)
const expandedRows = useState<EditorObject[]>(`${prefix}.expandedRows`, () => [])
const selectedRows = computed<EditorObject[]>({
  get: () => {
    const ids = selectedAnalysisIDs.value
    return editorObjects.value.filter((editorObject) => ids.includes(editorObject.id))
  },
  set: (value: EditorObject[]) => {
    const ids = value.map((row) => row.id)
    ids.sort()
    selectedAnalysisIDs.value = ids
  },
})

const editorObjects = computed<EditorObject[]>(() => props.analyses.map((item) => ({ ...analysisEditor(item, i18n), id: item.id })))

const selectedAnalyses = computed<Analysis[]>(() => selectedRows.value.map((row) => row.currentValue.value))

const saveChanges = (id: string) => {
  const index = editorObjects.value.findIndex((editorObject) => editorObject.id === id)
  const eo = presentOrFileBug(editorObjects.value[index])
  return withLoading(
    () => pactaClient.updateAnalysis(id, eo.changes.value).then(refresh),
    `${prefix}.saveChanges`,
  )
}

const deleteAnalysis = (id: string) => withLoading(
  () => pactaClient.deleteAnalysis(id),
  `${prefix}.deleteAnalysis`,
)
const deleteSelected = () => Promise.all([selectedRows.value.map((row) => deleteAnalysis(row.id))]).then(refresh)
</script>

<template>
  <div class="flex flex-column gap-3">
    <div class="flex gap-2 flex-wrap">
      <PVButton
        icon="pi pi-refresh"
        class="p-button-outlined p-button-secondary p-button-sm"
        :label="tt('Refresh')"
        @click="refresh"
      />
      <PVButton
        :disabled="!selectedRows || selectedRows.length === 0"
        icon="pi pi-trash"
        class="p-button-outlined p-button-danger p-button-sm"
        :label="tt('Delete') + selectedCountSuffix(selectedRows)"
        @click="deleteSelected"
      />
    </div>
    <PVDataTable
      v-model:selection="selectedRows"
      v-model:expanded-rows="expandedRows"
      :value="editorObjects"
      data-key="id"
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
        :header="tt('View')"
      >
        <template #body="slotProps">
          <LinkButton
            icon="pi pi-external-link"
            class="p-button-outlined p-button-xs"
            :label="tt('View')"
            :to="`${apiServerURL}/report/${slotProps.data.id}`"
            new-tab
          />
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
          <StandardDebug
            always
            :value="slotProps.data.currentValue.value"
            label="Raw Data"
          />
          <h2 class="mt-5">
            {{ tt('Editable Properties') }}
          </h2>
          <AnalysisEditor
            v-model:editor-values="slotProps.data.editorValues.value"
            :editor-fields="slotProps.data.editorFields.value"
          />
          <div class="flex gap-3 justify-content-between">
            <PVButton
              icon="pi pi-trash"
              class="p-button-danger p-button-outlined"
              :label="tt('Delete')"
              @click="() => deleteAnalysis(slotProps.data.id)"
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
        </div>
      </template>
    </PVDataTable>
    <div class="flex flex-wrap gap-3 w-full justify-content-between">
      <!-- TODO(grady) Hook this up to something. -->
      <PVButton
        class="p-button-outlined"
        :label="tt('How To Run a Report')"
        icon="pi pi-question-circle"
        icon-pos="right"
      />
    </div>
    <StandardDebug
      :value="selectedAnalyses"
      label="Selected Analyses"
    />
    <StandardDebug
      :value="props.analyses"
      label="All Analyses"
    />
  </div>
</template>
