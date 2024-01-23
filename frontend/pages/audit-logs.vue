<script setup lang="ts">
import { type AuditLogQuerySort, type AuditLogQueryWhere, type AuditLogQueryReq, AuditLogAction, AuditLogTargetType, AuditLogQuerySortBy, type AuditLog } from '@/openapi/generated/pacta'
import { urlReactiveAuditLogQuery } from '@/lib/auditlogquery'
import { type DataTableSortMeta } from 'primevue/datatable'
import { FilterMatchMode } from 'primevue/api'

const prefix = 'pages/audit-logs'

const { linkToPortfolio, linkToPortfolioGroup, linkToAnalysis, linkToIncompleteUploadList } = useMyDataURLs()
const localePath = useLocalePath()
const { fromQueryReactiveWithDefault, waitForURLToUpdate } = useURLParams()
const { t } = useI18n()
const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()
const { getMaybeMe } = useSession()
const { humanReadableTimeFromStandardString } = useTime()

const { maybeMe } = await getMaybeMe()
const tt = (s: string) => t(`${prefix}.${s}`)

const auditLogQuery = urlReactiveAuditLogQuery(fromQueryReactiveWithDefault)
const defaultSelectedColumns = 'createdAt,actor,action,primaryTarget,secondaryTarget'
const selectedColumnsQuery = fromQueryReactiveWithDefault('cols', defaultSelectedColumns)

const resetColumns = () => { selectedColumnsQuery.value = defaultSelectedColumns }

interface DataResponse {
  auditLogs: AuditLog[]
  cursor: string | undefined
  hasNextPage: boolean
}
const dataResponse = useState<DataResponse>(`${prefix}.dataResponse`, () => ({
  auditLogs: [] as AuditLog[],
  cursor: undefined,
  hasNextPage: true,
}))
const expandedRows = useState<AuditLog[]>(`${prefix}.expandedRows`, () => [])
const refreshAuditLogs = async (): Promise<void> => {
  await withLoading(
    () => pactaClient.listAuditLogs(auditLogQuery.value).then((resp) => { dataResponse.value = resp }),
    'audit-logs.refresh',
  )
}
const lastLoadedQuery = useState<string>(`${prefix}.currentlyLoading`, () => '')
const activelyRefreshing = useState<boolean>(`${prefix}.activelyRefreshing`, () => false)
const refreshAuditLogsIfQueryIsStale = async (): Promise<void> => {
  if (activelyRefreshing.value) {
    return
  }
  activelyRefreshing.value = true
  const queryValue = auditLogQuery.value
  const current = JSON.stringify(queryValue)
  if (lastLoadedQuery.value !== current) {
    lastLoadedQuery.value = current
    await refreshAuditLogs()
  }
  activelyRefreshing.value = false
}
watch(auditLogQuery, refreshAuditLogsIfQueryIsStale, { deep: true })
const auditLogs = computed<AuditLog[]>(() => dataResponse.value.auditLogs)
const hasNextPage = computed<boolean>(() => dataResponse.value.hasNextPage)
const hasPrevPage = computed<boolean>(() => showingRange.value.start > 1)
const showingRange = computed<{ start: number, end: number }>(() => {
  const start = auditLogQuery.value.cursor ? parseInt(auditLogQuery.value.cursor) : 1
  const end = start + auditLogs.value.length
  return { start, end }
})
const clickPrev = () => {
  if (hasPrevPage.value) {
    cursor.value = `${Math.max(showingRange.value.start - (auditLogQuery.value.limit ?? 100), 0)}`
    void withLoading(refreshAuditLogs, `${prefix}.`)
  } else {
    console.warn('clickPrev called when hasPrevPage is false')
  }
}
const clickNext = () => {
  if (hasNextPage.value) {
    cursor.value = dataResponse.value.cursor
    void withLoading(refreshAuditLogs, `${prefix}.next`)
  } else {
    console.warn('clickNext called when hasNext is false')
  }
}

const doUpdate = (fn: (v: AuditLogQueryReq) => void): void => {
  const newQuery = { ...auditLogQuery.value }
  fn(newQuery)
  auditLogQuery.value = newQuery
}
const wheres = computed<AuditLogQueryWhere[]>({
  get: () => auditLogQuery.value.wheres,
  set: (v: AuditLogQueryWhere[]) => { doUpdate((alq: AuditLogQueryReq) => { alq.wheres = v }) },
})
const sorts = computed<AuditLogQuerySort[]>({
  get: () => auditLogQuery.value.sorts ?? [],
  set: (v: AuditLogQuerySort[]) => { doUpdate((alq: AuditLogQueryReq) => { alq.sorts = v }) },
})
const cursor = computed<string | undefined>({
  get: () => auditLogQuery.value.cursor,
  set: (v: string | undefined) => { doUpdate((alq: AuditLogQueryReq) => { alq.cursor = v }) },
})

enum FilterType {
  StringBased,
  EnumBased,
  TimeBased,
}
interface OurDataTableFilterMeta {
  id: {
    value: string[] | undefined
    matchMode: string
  }
  createdAt: {
    value: Array<string | null> | undefined
    matchMode: string
  }
  actor: {
    value: string[] | undefined
    matchMode: string
  }
  action: {
    value: AuditLogAction[] | undefined
    matchMode: string
  }
  primaryTarget: {
    value: string[] | undefined
    matchMode: string
  }
  secondaryTarget: {
    value: string[] | undefined
    matchMode: string
  }
  actorId: {
    value: string[] | undefined
    matchMode: string
  }
  actorOwnerId: {
    value: string[] | undefined
    matchMode: string
  }
  primaryTargetType: {
    value: AuditLogTargetType[] | undefined
    matchMode: string
  }
  primaryTargetId: {
    value: string[] | undefined
    matchMode: string
  }
  primaryTargetOwner: {
    value: string[] | undefined
    matchMode: string
  }
  secondaryTargetType: {
    value: AuditLogTargetType[] | undefined
    matchMode: string
  }
  secondaryTargetId: {
    value: string[] | undefined
    matchMode: string
  }
  secondaryTargetOwner: {
    value: string[] | undefined
    matchMode: string
  }
}
interface Column {
  field: keyof OurDataTableFilterMeta
  header: string
  sortBy?: AuditLogQuerySortBy | undefined
  customBody?: boolean
  filterType?: FilterType | undefined
  enumValues?: Array<{ value: string, label: string }> | undefined
}

const allColumns: Column[] = [
  {
    field: 'id',
    header: tt('ID'),
    filterType: FilterType.StringBased,
  },
  {
    field: 'createdAt',
    header: tt('Time'),
    sortBy: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_CREATED_AT,
    customBody: true,
    filterType: FilterType.TimeBased,
  },
  {
    field: 'actor',
    header: tt('Actor'),
    sortBy: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_ACTOR_ID,
    customBody: true,
  },
  {
    field: 'action',
    header: tt('Action'),
    customBody: true,
    filterType: FilterType.EnumBased,
    enumValues: Object.values(AuditLogAction).map((a: AuditLogAction) => ({
      label: a.replace('AuditLogAction', ''),
      value: a,
    })).sort((a, b) => a.label.localeCompare(b.label)),
  },
  {
    field: 'primaryTarget',
    header: tt('Primary Target'),
    sortBy: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_PRIMARY_TARGET_TYPE,
    customBody: true,
  },
  {
    field: 'secondaryTarget',
    header: tt('Secondary Target'),
    sortBy: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_SECONDARY_TARGET_TYPE,
    customBody: true,
  },
  {
    field: 'actorId',
    header: tt('Actor ID'),
    sortBy: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_ACTOR_ID,
    filterType: FilterType.StringBased,
  },
  {
    field: 'actorOwnerId',
    header: tt('Actor Owner ID'),
    sortBy: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_ACTOR_OWNER_ID,
    filterType: FilterType.StringBased,

  },
  {
    field: 'primaryTargetType',
    header: tt('Primary Target Type'),
    sortBy: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_PRIMARY_TARGET_TYPE,
    filterType: FilterType.EnumBased,
    enumValues: Object.values(AuditLogTargetType).map((a: AuditLogTargetType) => ({
      label: a.replace('AuditLogTargetType', ''),
      value: a,
    })).sort((a, b) => a.label.localeCompare(b.label)),
  },
  {
    field: 'primaryTargetId',
    header: tt('Primary Target ID'),
    sortBy: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_PRIMARY_TARGET_ID,
    filterType: FilterType.StringBased,
  },
  {
    field: 'primaryTargetOwner',
    header: tt('Primary Target Owner'),
    sortBy: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_PRIMARY_TARGET_OWNER_ID,
    filterType: FilterType.StringBased,
  },
  {
    field: 'secondaryTargetType',
    header: tt('Secondary Target Type'),
    sortBy: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_SECONDARY_TARGET_TYPE,
  },
  {
    field: 'secondaryTargetId',
    header: tt('Secondary Target ID'),
    sortBy: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_SECONDARY_TARGET_ID,
  },
  {
    field: 'secondaryTargetOwner',
    header: tt('Secondary Target Owner'),
    sortBy: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_SECONDARY_TARGET_OWNER_ID,
  },
]
const selectedColumns = computed<Column[]>({
  get: () => {
    const selectedColumns = selectedColumnsQuery.value.split(',')
    return allColumns.filter(c => selectedColumns.includes(c.field))
  },
  set: (v: Column[]) => { selectedColumnsQuery.value = v.map(c => c.field).join(',') },
})

const tableFilterModel = computed<OurDataTableFilterMeta>({
  get: () => {
    const partial: Partial<OurDataTableFilterMeta> = {}
    for (const column of allColumns) {
      switch (column.filterType) {
        case FilterType.StringBased:
        case FilterType.EnumBased:
          partial[column.field] = { value: undefined, matchMode: FilterMatchMode.IN }
          break
        case FilterType.TimeBased:
          partial[column.field] = { value: undefined, matchMode: FilterMatchMode.BETWEEN }
          break
        case undefined:
          break
      }
    }
    const result = partial as OurDataTableFilterMeta
    const ws = wheres.value
    for (const w of ws) {
      if (w.inId) {
        result.id.value = w.inId
      }
      if (w.minCreatedAt) {
        if (!result.createdAt.value) {
          result.createdAt.value = [null, null]
        }
        presentOrFileBug(result.createdAt.value)[0] = w.minCreatedAt
      }
      if (w.maxCreatedAt) {
        if (!result.createdAt.value) {
          result.createdAt.value = [null, null]
        }
        presentOrFileBug(result.createdAt.value)[1] = w.maxCreatedAt
      }
      if (w.inAction) {
        result.action.value = w.inAction
      }
      if (w.inActorId) {
        result.actorId.value = w.inActorId
      }
      if (w.inActorOwnerId) {
        result.actorOwnerId.value = w.inActorOwnerId
      }
      if (w.inTargetType) {
        result.primaryTargetType.value = w.inTargetType
      }
      if (w.inTargetId) {
        result.primaryTargetId.value = w.inTargetId
      }
      if (w.inTargetOwnerId) {
        result.primaryTargetOwner.value = w.inTargetOwnerId
      }
    }
    return result
  },
  set: (v: OurDataTableFilterMeta) => {
    const result: AuditLogQueryWhere[] = []
    for (const [fieldName, filterState] of Object.entries(v)) {
      if (filterState.value === undefined || filterState.value === null) {
        continue
      }
      switch (fieldName) {
        case 'id':
          result.push({
            inId: filterState.value,
          })
          break
        case 'createdAt':
          result.push({
            minCreatedAt: (filterState.value)[0],
          }, {
            maxCreatedAt: (filterState.value)[1],
          })
          break
        case 'action':
          result.push({
            inAction: filterState.value,
          })
          break
        case 'actorId':
          result.push({
            inActorId: filterState.value,
          })
          break
        case 'actorOwnerId':
          result.push({
            inActorOwnerId: filterState.value,
          })
          break
        case 'primaryTargetType':
          result.push({
            inTargetType: filterState.value,
          })
          break
        case 'primaryTargetId':
          result.push({
            inTargetId: filterState.value,
          })
          break
        case 'primaryTargetOwner':
          result.push({
            inTargetOwnerId: filterState.value,
          })
          break
        default:
          console.log(`Unsupported field name ${fieldName}`)
          break
      }
    }
    wheres.value = result
  },
})
const tableSortModel = computed<DataTableSortMeta[]>({
  get: () => {
    return sorts.value.map((s): DataTableSortMeta | null => {
      const column = allColumns.find((c) => c.sortBy === s.by)
      if (!column) {
        console.warn(`Could not find column for sort ${s.by}`)
        return null
      }
      return {
        field: column.field,
        order: s.ascending ? 1 : -1,
      }
    }).filter((s) => s !== null).map((s) => s as DataTableSortMeta)
  },
  set: (v: DataTableSortMeta[]) => {
    const result = v.map((s) => {
      const column = allColumns.find((c) => c.field === s.field)
      if (!column) {
        console.warn(`Could not find column for field ${s.field}`)
        return null
      }
      return {
        by: column.sortBy,
        ascending: s.order === 1,
      }
    }).filter((s) => (s !== null)).map((s) => s as AuditLogQuerySort)
    sorts.value = result
  },
})

const getTargetLink = (t: AuditLogTargetType, id: string): string => {
  switch (t) {
    case AuditLogTargetType.AUDIT_LOG_TARGET_TYPE_PORTFOLIO:
      return linkToPortfolio(id)
    case AuditLogTargetType.AUDIT_LOG_TARGET_TYPE_PORTFOLIO_GROUP:
      return linkToPortfolioGroup(id)
    case AuditLogTargetType.AUDIT_LOG_TARGET_TYPE_INCOMPLETE_UPLOAD:
      return linkToIncompleteUploadList()
    case AuditLogTargetType.AUDIT_LOG_TARGET_TYPE_ANALYSIS:
    case AuditLogTargetType.AUDIT_LOG_TARGET_TYPE_ANALYSIS_ARTIFACT:
      return linkToAnalysis(id)
    case AuditLogTargetType.AUDIT_LOG_TARGET_TYPE_USER:
      return localePath(`/user/${id}`)
    case AuditLogTargetType.AUDIT_LOG_TARGET_TYPE_INITIATIVE:
      return localePath(`/initiative/${id}`)
    case AuditLogTargetType.AUDIT_LOG_TARGET_TYPE_PACTA_VERSION:
      return localePath(`/admin/pacta-version/${id}`)
  }
  console.warn(`Unknown target type ${t}`)
  return '#'
}

const columnControlOverlay = ref<{ toggle: (e: Event) => void }>()
const toggleColumnControl = (e: Event) => {
  const c = columnControlOverlay.value
  if (c) {
    c.toggle(e)
  } else {
    console.warn('columnControlOverlay is not available')
  }
}

const defaultAuditLogQuery = (): AuditLogQueryReq => {
  const mm = maybeMe.value
  return {
    limit: 100,
    cursor: undefined,
    wheres: [{
      inActorId: mm ? [mm.id] : [],
    }],
    sorts: [{
      by: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_CREATED_AT,
      ascending: false,
    }],
  }
}
const resetFilters = () => {
  auditLogQuery.value = defaultAuditLogQuery()
}
onMounted(() => {
  if (auditLogQuery.value.wheres.length === 0) {
    auditLogQuery.value = defaultAuditLogQuery()
  }
  void waitForURLToUpdate().then(refreshAuditLogs)
})
</script>

<template>
  <StandardContent class="portfolios-page">
    <TitleBar title="Audit Logs" />
    <p>
      {{ tt('DescriptionParagraph1') }}
    </p>
    <p>
      {{ tt('DescriptionParagraph2') }}
    </p>
    <div class="flex gap-2 flex-wrap">
      <StandardDebug
        label="Data Response"
        :value="dataResponse"
      />
      <StandardDebug
        label="Audit Log Query"
        :value="auditLogQuery"
      />
      <StandardDebug
        label="Filters"
        :value="tableFilterModel"
      />
      <StandardDebug
        label="Sorts"
        :value="tableSortModel"
      />
    </div>
    <div class="flex gap-2">
      <PVButton
        icon="pi pi-cog"
        :label="tt('Columns')"
        class="p-button-outlined"
        @click="toggleColumnControl"
      />
      <PVButton
        icon="pi pi-refresh"
        :label="tt('Reset')"
        class="p-button-outlined p-button-secondary"
        @click="resetFilters"
      />
    </div>
    <PVOverlayPanel
      ref="columnControlOverlay"
      class="max-w-screen w-25rem"
    >
      <div class="flex flex-column gap-2">
        <div class="flex justify-content-between align-items-center">
          <div class="flex flex-column gap-1">
            <div class="text-lg font-bold">
              {{ tt('Column Selection') }}
            </div>
            <span class="text-sm">
              {{ tt('Select or Remove Columns To Display') }}
            </span>
          </div>
          <div class="flex gap-1">
            <PVButton
              v-tooltip="'Reset'"
              icon="pi pi-refresh"
              class="p-button-text p-button-secondary"
              rounded
              @click="resetColumns"
            />
            <PVButton
              v-tooltip="'Done'"
              icon="pi pi-check"
              class="p-button-text p-button-secondary"
              rounded
              @click="toggleColumnControl"
            />
          </div>
        </div>
        <PVMultiSelect
          v-model="selectedColumns"
          display="chip"
          :options="allColumns"
          option-label="header"
          placeholder="Select Columns"
        />
      </div>
    </PVOverlayPanel>
    <StandardFullWidth>
      <PVDataTable
        v-model:filters="tableFilterModel"
        v-model:multiSortMeta="tableSortModel"
        v-model:expanded-rows="expandedRows"
        filter-display="menu"
        sort-mode="multiple"
        removable-sort
        :value="auditLogs"
        size="small"
        class="audit-log-data-table"
      >
        <PVColumn
          v-for="column in selectedColumns"
          :key="column.field"
          :sortable="column.sortBy !== undefined"
          :field="column.field"
          :header="column.header"
          :show-filter-match-modes="false"
        >
          <template
            v-if="column.customBody"
            #body="slotProps"
          >
            <template v-if="column.field === 'createdAt'">
              {{ humanReadableTimeFromStandardString(slotProps.data.createdAt).value }}
            </template>
            <template
              v-else-if="column.field === 'actor'"
            >
              <div class="flex flex-column gap-1">
                <span class="text-sm">
                  {{ slotProps.data.actorId }}
                </span>
                <span class="text-sm">
                  {{ slotProps.data.actorOwnerId }}
                </span>
                <div class="flex align-items-center gap-2">
                  <span>
                    Acting as <b>{{ `${slotProps.data.actorType}`.replace("AuditLogActorType", '') }}</b>
                  </span>
                  <LinkButton
                    :to="getTargetLink(AuditLogTargetType.AUDIT_LOG_TARGET_TYPE_USER, slotProps.data.actorId)"
                    class="p-button-xs p-button-outlined p-button-secondary"
                    icon="pi pi-external-link"
                  />
                </div>
              </div>
            </template>
            <template v-else-if="column.field === 'action'">
              {{ `${slotProps.data.action}`.replace("AuditLogAction", '') }}
            </template>
            <template v-else-if="column.field === 'primaryTarget'">
              <div class="flex flex-column gap-1">
                <div class="flex align-items-center gap-2">
                  <b>{{ slotProps.data.primaryTargetType.replace("AuditLogTargetType", '') }}</b>
                  <LinkButton
                    :to="getTargetLink(AuditLogTargetType.AUDIT_LOG_TARGET_TYPE_USER, slotProps.data.actorId)"
                    class="p-button-xs p-button-outlined p-button-secondary"
                    icon="pi pi-external-link"
                  />
                </div>
                <span class="text-sm">
                  {{ slotProps.data.primaryTargetId }}
                </span>
                <span class="text-sm">
                  {{ slotProps.data.primaryTargetOwner }}
                </span>
              </div>
            </template>
            <template v-else-if="column.field === 'secondaryTarget'">
              <div
                v-if="slotProps.data.secondaryTargetId"
                class="flex flex-column gap-1"
              >
                <b>
                  {{ slotProps.data.secondaryTargetType.replace("AuditLogTargetType", '') }}
                  <LinkButton
                    :to="getTargetLink(AuditLogTargetType.AUDIT_LOG_TARGET_TYPE_USER, slotProps.data.actorId)"
                    class="p-button-xs p-button-outlined p-button-secondary"
                    icon="pi pi-external-link"
                  />
                </b>
                <span>
                  {{ slotProps.data.secondaryTargetId }}
                </span>
                <span>
                  {{ slotProps.data.secondaryTargetOwner }}
                </span>
              </div>
            </template>
          </template>
          <template
            v-if="column.filterType !== undefined"
            #filter="{ filterModel }"
          >
            <div class="flex flex-column gap-3 w-26rem max-w-screen">
              <div class="text-lg font-bold">
                Filter by {{ column.header }}
              </div>
              <PVMultiSelect
                v-if="column.filterType === FilterType.EnumBased"
                v-model="filterModel.value"
                option-label="label"
                option-value="value"
                :options="column.enumValues"
                class="p-column-filter"
                display="chip"
                :placeholder="tt('Filter To Specific Values')"
              />
              <PVChips
                v-if="column.filterType === FilterType.StringBased"
                v-model="filterModel.value"
                class="p-column-filter"
                :placeholder="tt('Filter To Specific Values')"
              />
              <template
                v-if="column.filterType === FilterType.TimeBased"
              >
                <CalendarRangeWithISOModel
                  v-model:value="filterModel.value"
                />
              </template>
            </div>
          </template>
        </PVColumn>
        <PVColumn
          expander
          header="Details"
        />
        <template #expansion="slotProps">
          <div class="code-block my-2">
            {{ slotProps.data }}
          </div>
        </template>
        <template
          #footer
        >
          <div class="flex align-items-center gap-3 flex-wrap">
            <PVButton
              v-if="hasPrevPage"
              label="Previous Page"
              icon="pi pi-caret-left"
              class="p-button-outlined p-button-xs p-button-secondary"
              @click="clickPrev"
            />
            <span v-if="showingRange">
              Showing Logs
              {{ showingRange.start }}
              -
              {{ showingRange.end }}
            </span>
            <PVButton
              v-if="hasNextPage"
              label="Next Page"
              icon="pi pi-caret-right"
              class="p-button-outlined p-button-xs p-button-secondary"
              icon-pos="right"
              @click="clickNext"
            />
          </div>
        </template>
      </PVDataTable>
    </StandardFullWidth>
  </StandardContent>
</template>

<style lang="scss">
.p-column-filter-overlay-menu .p-column-filter-buttonbar{
  padding-top: 0;
}
</style>
