<script setup lang="ts">
import { portfolioEditor } from '@/lib/editor'
import { type Portfolio, type PortfolioGroup, type Initiative, type Analysis } from '@/openapi/generated/pacta'
import { selectedCountSuffix } from '@/lib/selection'

const {
  humanReadableTimeFromStandardString,
} = useTime()
const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()
const localePath = useLocalePath()
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

const selectedPortfolioIDs = computed({
  get: () => props.selectedPortfolioIds ?? [],
  set: (value: string[]) => { emit('update:selectedPortfolioIds', value) },
})

interface EditorObject extends ReturnType<typeof portfolioEditor> {
  id: string
  analyses: Analysis[]
}

const prefix = 'components/portfolio/ListView'
const tt = (s: string) => t(`${prefix}.${s}`)
const expandedRows = useState<EditorObject[]>(`${prefix}.expandedRows`, () => [])
const selectedRows = computed<EditorObject[]>({
  get: () => {
    const ids = selectedPortfolioIDs.value
    return editorObjects.value.filter((editorObject) => ids.includes(editorObject.id))
  },
  set: (value: EditorObject[]) => {
    const ids = value.map((row) => row.id)
    ids.sort()
    selectedPortfolioIDs.value = ids
  },
})

const editorObjects = computed<EditorObject[]>(() => props.portfolios.map((item) => ({
  id: item.id,
  ...portfolioEditor(item, i18n),
  analyses: props.analyses.filter((a) => a.portfolioSnapshot.portfolio?.id === item.id),
})))

const selectedPortfolios = computed<Portfolio[]>(() => selectedRows.value.map((row) => row.currentValue.value))

const saveChanges = (id: string) => {
  const index = editorObjects.value.findIndex((editorObject) => editorObject.id === id)
  const eo = presentOrFileBug(editorObjects.value[index])
  return withLoading(
    () => pactaClient.updatePortfolio(id, eo.changes.value).then(refresh),
    `${prefix}.saveChanges`,
  )
}

const deletePortfolio = (id: string) => withLoading(
  () => pactaClient.deletePortfolio(id),
  `${prefix}.deletePortfolio`,
)
const deleteThisPortfolio = (id: string) => deletePortfolio(id).then(refresh)
const deleteSelected = () => Promise.all([selectedRows.value.map((row) => deletePortfolio(row.id))]).then(refresh)
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
      <PortfolioGroupMembershipMenuButton
        :selected-portfolios="selectedPortfolios"
        :portfolio-groups="props.portfolioGroups"
        @changed-memberships="refresh"
        @changed-groups="refresh"
      />
      <PortfolioInitiativeMembershipMenuButton
        :selected-portfolios="selectedPortfolios"
        :initiatives="props.initiatives"
        @changed-memberships="refresh"
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
      <template #empty>
        <PVMessage severity="info">
          {{ tt('No Uploaded Portfolios Message') }}
        </PVMessage>
      </template>
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
      <PVColumn
        field="currentValue.value.name"
        sortable
        :header="tt('Name')"
      />
      <PVColumn
        :header="tt('Memberships')"
      >
        <template #body="slotProps">
          <div class="flex flex-column gap-2">
            <div
              v-if="slotProps.data.currentValue.value.groups.length > 0"
              class="flex gap-1 align-items-center flex-wrap"
            >
              <span>{{ tt('Groups') }}:</span>
              <LinkButton
                v-for="membership in slotProps.data.currentValue.value.groups"
                :key="membership.portfolioGroup.id"
                class="p-button-outlined p-button-xs"
                icon="pi pi-table"
                :label="membership.portfolioGroup.name"
                :to="localePath(`/my-data?tab=pg&pgids=${membership.portfolioGroup.id}`)"
              />
            </div>
            <div
              v-if="slotProps.data.currentValue.value.initiatives.length > 0"
              class="flex gap-1 align-items-center flex-wrap"
            >
              <span>{{ tt('Initiatives') }}:</span>
              <LinkButton
                v-for="membership in slotProps.data.currentValue.value.initiatives"
                :key="membership.initiative.id"
                class="p-button-xs"
                :label="membership.initiative.name"
                icon="pi pi-arrow-right"
                icon-pos="right"
                :to="localePath(`/initiative/${membership.initiative.id}`)"
              />
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
          <StandardDebug
            always
            :value="slotProps.data.currentValue.value"
            label="Raw Data"
          />
          <h2 class="mt-5">
            {{ tt('Memberships') }}
          </h2>
          <div class="flex flex-column gap-2">
            <PortfolioGroupMembershipMenuButton
              :selected-portfolios="[slotProps.data.currentValue.value]"
              :portfolio-groups="props.portfolioGroups"
              @changed-memberships="refresh"
              @changed-groups="refresh"
            />
            <PortfolioInitiativeMembershipMenuButton
              :selected-portfolios="[slotProps.data.currentValue.value]"
              :initiatives="props.initiatives"
              @changed-memberships="refresh"
            />
          </div>
          <h2 class="mt-5">
            {{ tt('Analysis') }}
          </h2>
          <AnalysisContextualListView
            :analyses="slotProps.data.analyses"
            :name="slotProps.data.currentValue.value.name"
            :portfolio-id="slotProps.data.id"
            @refresh="refresh"
          />
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
              @click="() => deleteThisPortfolio(slotProps.data.id)"
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
      <LinkButton
        :class="props.portfolios.length > 0 ? 'p-button-outlined' : ''"
        icon="pi pi-arrow-right"
        icon-pos="right"
        to="/upload"
        :label="tt('Upload New Portfolios')"
      />
      <!-- TODO(grady) Hook this up to something. -->
      <PVButton
        v-if="props.portfolios.length > 0"
        class="p-button-outlined"
        :label="tt('How To Run a Report')"
        icon="pi pi-question-circle"
        icon-pos="right"
      />
    </div>
    <StandardDebug
      :value="selectedPortfolios"
      label="Selected Portfolios"
    />
    <StandardDebug
      :value="props.portfolios"
      label="All Portfolios"
    />
  </div>
</template>
