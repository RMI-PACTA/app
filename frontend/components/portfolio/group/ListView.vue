<script setup lang="ts">
import { portfolioGroupEditor } from '@/lib/editor'
import { type Portfolio, type PortfolioGroup, type PortfolioGroupMembershipPortfolio, type Analysis } from '@/openapi/generated/pacta'
import { selectedCountSuffix } from '@/lib/selection'

const {
  humanReadableTimeFromStandardString,
} = useTime()
const pactaClient = usePACTA()
const { loading: { withLoading }, newPortfolioGroup: { newPortfolioGroupVisible } } = useModal()
const i18n = useI18n()
const localePath = useLocalePath()
const { t } = i18n

interface Props {
  portfolios: Portfolio[]
  portfolioGroups: PortfolioGroup[]
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

const selectedPortfolioGroupIDs = computed({
  get: () => props.selectedPortfolioGroupIds ?? [],
  set: (value: string[]) => { emit('update:selectedPortfolioGroupIds', value) },
})

interface EditorObject extends ReturnType<typeof portfolioGroupEditor> {
  id: string
}

const prefix = 'components/portfolio/group/ListView'
const tt = (s: string) => t(`${prefix}.${s}`)

const editorObjects = computed<EditorObject[]>(() => props.portfolioGroups.map((item) => ({ ...portfolioGroupEditor(item, i18n), id: item.id })))

const expandedRows = useState<EditorObject[]>(`${prefix}.expandedRows`, () => [])
const selectedRows = computed<EditorObject[]>({
  get: () => {
    return editorObjects.value.filter((editorObject) => selectedPortfolioGroupIDs.value.includes(editorObject.id))
  },
  set: (value: EditorObject[]) => {
    selectedPortfolioGroupIDs.value = value.map((row) => row.id)
  },
})

const deletePortfolioGroup = (id: string) => withLoading(
  () => pactaClient.deletePortfolioGroup(id).then(() => {
    expandedRows.value = expandedRows.value.filter((row) => row.id !== id)
    emit('refresh')
  }),
  `${prefix}.deletePortfolioGroup`,
)
const deleteSelected = () => Promise.all([selectedRows.value.map((row) => deletePortfolioGroup(row.id))]).then(() => { emit('refresh') })
const saveChanges = (id: string) => {
  const index = editorObjects.value.findIndex((editorObject) => editorObject.id === id)
  const eo = presentOrFileBug(editorObjects.value[index])
  return withLoading(
    () => pactaClient.updatePortfolioGroup(id, eo.changes.value)
      .then(() => pactaClient.findPortfolioGroupById(id))
      .then((portfolio) => {
        editorObjects.value[index] = { ...portfolioGroupEditor(portfolio, i18n), id }
      }),
    `${prefix}.saveChanges`,
  )
}
const editorObjectToIds = (editorObject: EditorObject): string[] => {
  return (editorObject.currentValue.value.members ?? []).map((m: PortfolioGroupMembershipPortfolio) => m.portfolio.id)
}
</script>

<template>
  <div class="flex flex-column gap-3">
    <div class="flex gap-2 flex-wrap">
      <PVButton
        icon="pi pi-refresh"
        class="p-button-outlined p-button-secondary p-button-sm"
        :label="tt('Refresh')"
        @click="() => emit('refresh')"
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
      class="w-full"
      size="small"
      sort-field="currentValue.value.createdAt"
      :sort-order="-1"
    >
      <template #empty>
        <PVMessage severity="info">
          {{ tt('No Portfolio Groups Message') }}
        </PVMessage>
      </template>
      <PVColumn selection-mode="multiple" />
      <PVColumn
        field="currentValue.value.name"
        sortable
        :header="tt('Name')"
      />
      <PVColumn
        field="currentValue.value.createdAt"
        :header="tt('Created At')"
        sortable
      >
        <template #body="slotProps">
          {{ humanReadableTimeFromStandardString(slotProps.data.currentValue.value.createdAt).value }}
        </template>
      </PVColumn>
      <PVColumn
        :header="tt('Number of Members')"
      >
        <template #body="slotProps">
          <LinkButton
            :disabled="editorObjectToIds(slotProps.data).length === 0"
            :to="localePath(`/my-data?tab=p&pids=${ editorObjectToIds(slotProps.data).join(',')}`)"
            :label="`${editorObjectToIds(slotProps.data).length}`"
            icon="pi pi-th-large"
            class="py-1 px-2 p-button-outlined p-button-secondary"
          />
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
          <StandardDebug
            always
            :value="slotProps.data.currentValue.value"
            label="Raw Data"
          />
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
    <div class="flex flex-wrap gap-3 w-full justify-content-between">
      <PVButton
        :class="portfolioGroups.length > 0 ? 'p-button-outlined' : ''"
        icon="pi pi-plus"
        :label="tt('New Portfolio Group')"
        @click="() => newPortfolioGroupVisible = true"
      />
      <!-- TODO(grady) Hook this up to something. -->
      <PVButton
        v-if="portfolioGroups.length > 0"
        class="p-button-outlined"
        :label="tt('How To Run a Report')"
        icon="pi pi-question-circle"
        icon-pos="right"
      />
    </div>
  </div>
</template>
