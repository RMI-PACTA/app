<script setup lang="ts">
import { portfolioEditor } from '@/lib/editor'
import { type Portfolio, AuditLogQuerySortBy, type PortfolioGroup, type Initiative, type Analysis } from '@/openapi/generated/pacta'
import { selectedCountSuffix } from '@/lib/selection'
import { createURLAuditLogQuery } from '@/lib/auditlogquery'
import { type WritableComputedRef } from 'vue'

const { linkToPortfolioGroup } = useMyDataURLs()
const { humanReadableTimeFromStandardString } = useTime()
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
  expandedPortfolioIds: string[]
  expandedSections: Map<string, number[]>
}
const props = defineProps<Props>()
interface Emits {
  (e: 'update:selectedPortfolioIds', value: string[]): void
  (e: 'update:expandedPortfolioIds', value: string[]): void
  (e: 'update:expandedSections', value: Map<string, number[]>): void
  (e: 'refresh'): void
}
const emit = defineEmits<Emits>()

const refresh = () => { emit('refresh') }

const selectedPortfolioIdsModel = computed({
  get: () => props.selectedPortfolioIds ?? [],
  set: (value: string[]) => { emit('update:selectedPortfolioIds', value) },
})
const expandedPortfolioIdsModel = computed({
  get: () => props.expandedPortfolioIds ?? [],
  set: (value: string[]) => { emit('update:expandedPortfolioIds', value) },
})
const expandedSectionsModel = computed({
  get: () => props.expandedSections ?? new Map<string, number[]>(),
  set: (value: Map<string, number[]>) => { emit('update:expandedSections', value) },
})

interface EditorObject extends ReturnType<typeof portfolioEditor> {
  id: string
  analyses: Analysis[]
  expandedSections: WritableComputedRef<number[]>
}

const prefix = 'components/portfolio/ListView'
const tt = (s: string) => t(`${prefix}.${s}`)

const selectedRows = computed<EditorObject[]>({
  get: () => {
    const ids = selectedPortfolioIdsModel.value
    return editorObjects.value.filter((editorObject) => ids.includes(editorObject.id))
  },
  set: (value: EditorObject[]) => {
    const ids = value.map((row) => row.id)
    ids.sort()
    selectedPortfolioIdsModel.value = ids
  },
})
const readyToExpand = useState<boolean>(`${prefix}.readyToExpand`, () => false)
onMounted(() => {
  readyToExpand.value = true
})
const expandedRows = computed<EditorObject[]>({
  get: () => {
    if (!readyToExpand.value) {
      return []
    }
    const ids = expandedPortfolioIdsModel.value
    const result = editorObjects.value.filter((editorObject) => ids.includes(editorObject.id))
    return result
  },
  set: (value: EditorObject[]) => {
    const ids = value.map((row) => row.id)
    ids.sort()
    expandedPortfolioIdsModel.value = ids
  },
})

const editorObjects = computed<EditorObject[]>(() => props.portfolios.map((item) => {
  const expandedSectionSuffix = item.id.substring(item.id.length - 4)
  return ({
    ...portfolioEditor(item, i18n),
    id: item.id,
    analyses: props.analyses.filter((a) => a.portfolioSnapshot.portfolio?.id === item.id),
    expandedSections: computed<number[]>({
      get: () => expandedSectionsModel.value.get(expandedSectionSuffix) ?? [],
      set: (value: number[]) => {
        const m = new Map<string, number[]>(expandedSectionsModel.value)
        m.set(expandedSectionSuffix, value)
        expandedSectionsModel.value = m
      },
    }),
  })
}))

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

const auditLogURL = (id: string) => {
  return createURLAuditLogQuery(
    localePath,
    {
      sorts: [{ by: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_CREATED_AT, ascending: false }],
      wheres: [{ inTargetId: [id] }],
    },
  )
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
                :to="linkToPortfolioGroup(membership.portfolioGroup.id)"
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
          <h2 class="mb-3 mt-0">
            {{ tt('Portfolio') }}: {{ slotProps.data.currentValue.value.name }}
          </h2>
          <PVAccordion
            v-model:activeIndex="slotProps.data.expandedSections.value"
            :multiple="true"
          >
            <PVAccordionTab>
              <template #header>
                <CommonAccordionHeader
                  :heading="tt('EditHeading')"
                  :sub-heading="tt('EditSubHeading')"
                  icon="pi pi-pencil"
                />
              </template>
              <PortfolioEditor
                v-model:editor-values="slotProps.data.editorValues.value"
                :editor-fields="slotProps.data.editorFields.value"
              />
              <div class="flex gap-2">
                <PVButton
                  :disabled="!slotProps.data.canSave.value"
                  :label="tt('Discard Changes')"
                  icon="pi pi-refresh"
                  class="p-button-secondary p-button-outlined"
                  @click="slotProps.data.resetEditor"
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
            </PVAccordionTab>
            <PVAccordionTab>
              <template #header>
                <CommonAccordionHeader
                  :heading="tt('MembershipsHeading')"
                  :sub-heading="tt('MembershipsSubHeading')"
                >
                  <div class="flex gap-1 justify-content-center">
                    <PVInlineMessage
                      severity="info"
                      icon="pi pi-table"
                    >
                      {{ slotProps.data.currentValue.value.groups.length }}
                    </PVInlineMessage>
                    <PVInlineMessage
                      severity="info"
                      icon="pi pi-sitemap"
                    >
                      {{ slotProps.data.currentValue.value.initiatives.length }}
                    </PVInlineMessage>
                  </div>
                </CommonAccordionHeader>
              </template>
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
            </PVAccordionTab>
            <PVAccordionTab>
              <template #header>
                <CommonAccordionHeader
                  :heading="tt('AnalysesHeading')"
                  :sub-heading="tt('AnalysesSubHeading')"
                >
                  <PVInlineMessage
                    v-if="slotProps.data.analyses.length === 0"
                    severity="success"
                    icon="pi pi-copy"
                  >
                    {{ tt('AnalysesComeHereChip') }}
                  </PVInlineMessage>
                  <div
                    v-else
                    class="bg-red-500"
                  >
                    {{ slotProps.data.analyses.map((a: Analysis) => a.analysisType) }}
                  </div>
                </CommonAccordionHeader>
              </template>
              <AnalysisContextualListView
                :analyses="slotProps.data.analyses"
                :name="slotProps.data.currentValue.value.name"
                :portfolio-id="slotProps.data.id"
                @refresh="refresh"
              />
            </PVAccordionTab>
            <PVAccordionTab>
              <template #header>
                <CommonAccordionHeader
                  :heading="tt('MoreHeading')"
                  :sub-heading="tt('MoreSubHeading')"
                  icon="pi pi-plus"
                />
              </template>
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
              </FormField>
              <FormField
                :label="tt('Raw Portfolio Metadata')"
                :help-text="tt('RawPortfolioMetadataHelpText')"
              >
                <StandardDebug
                  always
                  :value="slotProps.data.currentValue.value"
                  :label="`${tt('Portfolio Metadata')}: ${slotProps.data.currentValue.value.name}`"
                />
              </FormField>
              <FormField
                :label="tt('Delete Portfolio')"
                :help-text="tt('DeletePortfolioHelpText')"
              >
                <PVButton
                  icon="pi pi-trash"
                  class="p-button-danger p-button-outlined align-self-start"
                  :label="tt('Delete')"
                  @click="() => deleteThisPortfolio(slotProps.data.id)"
                />
              </FormField>
            </PVAccordionTab>
          </PVAccordion>
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
