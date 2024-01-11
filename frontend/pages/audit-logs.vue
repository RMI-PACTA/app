<script setup lang="ts">
import { type AuditLogQuerySort, type AuditLogQueryReq, AuditLogTargetType, AuditLogQuerySortBy, type AuditLog } from '@/openapi/generated/pacta'
import { urlReactiveAuditLogQuery } from '@/lib/auditlogquery'
import { type DataTableSortMeta } from 'primevue/datatable'

const prefix = 'pages/audit-logs'

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
const selectedColumnsQuery = fromQueryReactiveWithDefault('cols', 'createdAt,actor,action,primaryTarget,secondaryTarget')

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

const refreshAuditLogs = () => withLoading(
  () => pactaClient.listAuditLogs(auditLogQuery.value).then((resp) => { dataResponse.value = resp }),
  'audit-logs.refresh',
)
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
/* When we support filtering, we'll need this
const wheres = computed<AuditLogQueryWhere[]>({
  get: () => auditLogQuery.value.wheres,
  set: (v: AuditLogQueryWhere[]) => { doUpdate((alq: AuditLogQueryReq) => { alq.wheres = v }) },
})
*/
const sorts = computed<AuditLogQuerySort[]>({
  get: () => auditLogQuery.value.sorts ?? [],
  set: (v: AuditLogQuerySort[]) => { doUpdate((alq: AuditLogQueryReq) => { alq.sorts = v }) },
})
const cursor = computed<string | undefined>({
  get: () => auditLogQuery.value.cursor,
  set: (v: string | undefined) => { doUpdate((alq: AuditLogQueryReq) => { alq.cursor = v }) },
})

interface Column {
  field: string
  header: string
  sortBy?: AuditLogQuerySortBy | undefined
  customBody?: boolean
}

const allColumns: Column[] = [
  { field: 'id', header: tt('ID') },
  {
    field: 'createdAt',
    header: tt('Time'),
    sortBy: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_CREATED_AT,
    customBody: true,
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
  },
  {
    field: 'actorOwnerId',
    header: tt('Actor Owner ID'),
    sortBy: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_ACTOR_OWNER_ID,
  },
  {
    field: 'primaryTargetType',
    header: tt('Primary Target Type'),
    sortBy: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_PRIMARY_TARGET_TYPE,
  },
  {
    field: 'primaryTargetId',
    header: tt('Primary Target ID'),
    sortBy: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_PRIMARY_TARGET_ID,
  },
  {
    field: 'primaryTargetOwner',
    header: tt('Primary Target Owner'),
    sortBy: AuditLogQuerySortBy.AUDIT_LOG_QUERY_SORT_BY_PRIMARY_TARGET_OWNER_ID,
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
      // TODO(grady) make this directly addressable
      return localePath(`/my-data/${id}`)
    case AuditLogTargetType.AUDIT_LOG_TARGET_TYPE_PORTFOLIO_GROUP:
      // TODO(grady) make this directly addressable
      return localePath(`/my-data/${id}`)
    case AuditLogTargetType.AUDIT_LOG_TARGET_TYPE_INCOMPLETE_UPLOAD:
      // TODO(grady) make this directly addressable
      return localePath(`/my-data/${id}`)
    case AuditLogTargetType.AUDIT_LOG_TARGET_TYPE_ANALYSIS:
    case AuditLogTargetType.AUDIT_LOG_TARGET_TYPE_ANALYSIS_ARTIFACT:
      // TODO(grady) make this directly addressable
      return localePath(`/my-data/${id}`)
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

onMounted(() => {
  if (auditLogQuery.value.wheres.length === 0) {
    const mm = maybeMe.value
    auditLogQuery.value = {
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
  void waitForURLToUpdate().then(refreshAuditLogs)
})
</script>

<template>
  <StandardContent class="portfolios-page">
    <TitleBar title="Audit Logs" />
    <p>
      TODO(#80) Add Copy Here
    </p>
    <PVButton
      icon="pi pi-refresh"
      label="Refresh"
      @click="refreshAuditLogs"
    />
    <PVMultiSelect
      v-model="selectedColumns"
      display="chip"
      :options="allColumns"
      option-label="header"
      placeholder="Select Columns"
      class="better-multiselect-layout"
    />
    <StandardDebug
      label="Query Sort"
      :value="sorts"
    />
    <StandardDebug
      label="Table Sort Model"
      :value="tableSortModel"
    />
    <PVDataTable
      v-model:multiSortMeta="tableSortModel"
      v-model:expanded-rows="expandedRows"
      sort-mode="multiple"
      removable-sort
      :value="auditLogs"
      size="small"
    >
      <PVColumn
        v-for="column in selectedColumns"
        :key="column.field"
        :sortable="column.sortBy !== undefined"
        :field="column.field"
        :header="column.header"
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
          <span>
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
    <StandardDebug
      label="Data Response Cursor"
      :value="dataResponse.cursor"
    />
    <StandardDebug
      label="Audit Log Query"
      :value="auditLogQuery"
    />
  </StandardContent>
</template>

<style lang="scss">
.better-multiselect-layout.p-multiselect.p-multiselect-chip {
  .p-multiselect-label {
    padding: 0.25rem;
    display: flex;
    gap: 0.25rem;
    flex-wrap: wrap;

    .p-multiselect-token {
      margin: 0;
    }
  }
}
</style>
