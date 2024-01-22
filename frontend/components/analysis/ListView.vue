<script setup lang="ts">
import { analysisEditor } from '@/lib/editor'
import { AuditLogQuerySortBy, type Portfolio, type AnalysisArtifact, type AnalysisArtifactChanges, type PortfolioGroup, type Initiative, type Analysis } from '@/openapi/generated/pacta'
import { selectedCountSuffix } from '@/lib/selection'
import { createURLAuditLogQuery } from '@/lib/auditlogquery'

const { humanReadableTimeFromStandardString } = useTime()
const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()
const i18n = useI18n()
const localePath = useLocalePath()
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
  artifactsADEState: boolean
  artifactsSTPState: boolean
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

const editorObjects = computed<EditorObject[]>(() => props.analyses.map((item: Analysis) => ({
  ...analysisEditor(item, i18n),
  id: item.id,
  artifactsADEState: item.artifacts.every(a => a.adminDebugEnabled),
  artifactsSTPState: item.artifacts.every(a => a.sharedToPublic),
})))

const selectedAnalyses = computed<Analysis[]>(() => selectedRows.value.map((row) => row.currentValue.value))

const saveChanges = (id: string) => {
  const index = editorObjects.value.findIndex((editorObject) => editorObject.id === id)
  const eo = presentOrFileBug(editorObjects.value[index])
  return withLoading(
    () => pactaClient.updateAnalysis(id, eo.changes.value).then(refresh),
    `${prefix}.saveChanges`,
  )
}

const auditLogURL = (id: string) => {
  return createURLAuditLogQuery(
    localePath,
    {
      sorts: [{ by: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_CREATED_AT, ascending: false }],
      wheres: [{ inTargetId: [id] }],
    },
  )
}

const changeAllArtifacts = (id: string, changes: AnalysisArtifactChanges) => {
  const artifacts: AnalysisArtifact[] = editorObjects.value.find((eo) => eo.id === id)?.currentValue.value.artifacts ?? []
  void withLoading(async () => {
    for (const artifact of artifacts) {
      await pactaClient.updateAnalysisArtifact(artifact.id, changes)
    }
    refresh()
  }, `${prefix}.changeAllArtifacts`)
}

const setAllADEOnArtifacts = (id: string, value: boolean) => {
  changeAllArtifacts(id, { adminDebugEnabled: value })
}
const setAllSTPOnArtifacts = (id: string, value: boolean) => {
  changeAllArtifacts(id, { sharedToPublic: value })
}
const deleteAnalysis = (id: string) => withLoading(
  () => pactaClient.deleteAnalysis(id),
  `${prefix}.deleteAnalysis`,
)
const deleteSelected = () => Promise.all([selectedRows.value.map((row) => deleteAnalysis(row.id))]).then(refresh)
const deleteSpecificAnalysis = async (id: string) => {
  await deleteAnalysis(id)
  refresh()
}
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
      sort-field="currentValue.value.createdAt"
      :sort-order="-1"
    >
      <PVColumn selection-mode="multiple" />
      <PVColumn
        field="currentValue.value.createdAt"
        :header="tt('Created At')"
        sortable
      >
        <template #body="slotProps">
          {{ humanReadableTimeFromStandardString(slotProps.data.currentValue.value.createdAt).value }}
        </template>
      </PVColumn>
      <PVColumn :header="tt('Status')">
        <template #body="slotProps">
          <AnalysisStatusChip
            :analysis="slotProps.data.currentValue.value"
          />
        </template>
      </PVColumn>
      <PVColumn
        field="currentValue.value.name"
        sortable
        :header="tt('Name')"
      />
      <PVColumn
        :header="tt('View')"
      >
        <template #body="slotProps">
          <AnalysisAccessButtons
            :analysis="slotProps.data.currentValue.value"
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
              @click="() => deleteSpecificAnalysis(slotProps.data.id)"
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
          <h2 class="mt-5">
            {{ tt('Access Controls') }}
          </h2>
          <FormField
            :label="tt('Admin Debugging Enabled')"
            :help-text="tt('ADEHelpText')"
          >
            <AdminDebugEnabledToggleButton
              :value="slotProps.data.artifactsADEState"
              @update:value="(newValue: boolean) => setAllADEOnArtifacts(slotProps.data.id, newValue)"
            />
          </FormField>
          <FormField
            :label="tt('Shared To Public')"
            :help-text="tt('STPHelpText')"
          >
            <SharedToPublicToggleButton
              :value="slotProps.data.artifactsSTPState"
              @update:value="(newValue: boolean) => setAllSTPOnArtifacts(slotProps.data.id, newValue)"
            />
          </FormField>
          <FormField
            :label="tt('Audit Logs')"
            :help-text="tt('AuditLogsHelpText')"
          >
            <LinkButton
              :label="tt('View Audit Logs')"
              :to="auditLogURL(slotProps.data.id)"
              icon="pi pi-arrow-right"
              class="p-button-outlined align-self-start"
              icon-pos="right"
            />
          </formfield>
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
