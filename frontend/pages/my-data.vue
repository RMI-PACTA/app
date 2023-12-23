<script setup lang="ts">
const prefix = 'pages/my-data'

const { fromQueryReactive } = useURLParams()
const { t } = useI18n()
const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()

const tt = (s: string) => t(`${prefix}.${s}`)

const selectedPortfolioIdsQP = fromQueryReactive('pids')
const selectedPortfolioGroupIdsQP = fromQueryReactive('pgids')
const tabQP = fromQueryReactive('tab')

const [
  { data: incompleteUploadsData, refresh: refreshIncompleteUploadsApi },
  { data: portfolioData, refresh: refreshPortfoliosApi },
  { data: portfolioGroupData, refresh: refreshPortfolioGroupsApi },
  { data: initiativeData, refresh: refreshInitiativesApi },
] = await Promise.all([
  useSimpleAsyncData(`${prefix}.incompleteUploads`, () => pactaClient.listIncompleteUploads()),
  useSimpleAsyncData(`${prefix}.portfolios`, () => pactaClient.listPortfolios()),
  useSimpleAsyncData(`${prefix}.portfolioGroups`, () => pactaClient.listPortfolioGroups()),
  useSimpleAsyncData(`${prefix}.initiatives`, () => pactaClient.listInitiatives()),
])
const refreshIncompleteUploads = async () => {
  await withLoading(refreshIncompleteUploadsApi, `${prefix}.refreshIncompleteUploads`)
}
const refreshPortfolios = async () => {
  await withLoading(refreshPortfoliosApi, `${prefix}.refreshPortfolios`)
}
const refreshPortfolioGroups = async () => {
  await withLoading(refreshPortfolioGroupsApi, `${prefix}.refreshPortfolioGroups`)
}
const refreshInitiatives = async () => {
  await withLoading(refreshInitiativesApi, `${prefix}.refreshInitiatives`)
}
const refreshAll = async () => {
  console.log('refreshing all')
  await Promise.all([
    refreshIncompleteUploads(),
    refreshPortfolios(),
    refreshPortfolioGroups(),
    refreshInitiatives(),
  ])
}

const selectedPortfolioIds = computed<string[]>({
  get: () => (selectedPortfolioIdsQP.value ?? '').split(','),
  set: (v: string[]) => {
    if (v.length === 0) {
      selectedPortfolioIdsQP.value = undefined
    } else {
      v.sort()
      selectedPortfolioIdsQP.value = v.join(',')
    }
  },
})
const selectedPortfolioGroupIds = computed<string[]>({
  get: () => (selectedPortfolioGroupIdsQP.value ?? '').split(','),
  set: (v: string[]) => {
    if (v.length === 0) {
      selectedPortfolioGroupIdsQP.value = undefined
    } else {
      v.sort()
      selectedPortfolioGroupIdsQP.value = v.join(',')
    }
  },
})
interface TabToIndexMap {
  iu: number
  p: number
  pg: number
}
const tabToIndexMap = computed(() => {
  const result: TabToIndexMap = {
    iu: -1,
    p: -1,
    pg: -1,
  }
  let idx = 0
  if (incompleteUploadsData.value.items.length > 0) {
    result.iu = idx
    idx++
  } else {
    result.iu = -1
  }
  result.p = idx++
  result.pg = idx++
  return result
})
const activeIndex = computed<number>({
  get: () => {
    const tab = (tabQP.value ?? 'p') as keyof TabToIndexMap
    const result = tabToIndexMap.value[tab]
    if (result === undefined) {
      console.error(`Unknown tab ${tab}`)
      return 0
    }
    return result
  },
  set: (vv: number) => {
    const ttim = tabToIndexMap.value
    const tab = Object.entries(ttim).find(([k, v]) => v === vv)
    if (!tab) {
      console.error(`Unknown tab index ${vv}`)
      return
    }
    tabQP.value = tab[0]
  },
})
</script>

<template>
  <StandardContent class="portfolios-page">
    <TitleBar :title="tt('My Data')" />
    <p>
      TODO(#80) Add Copy Here
    </p>
    <PVTabView
      v-model:activeIndex="activeIndex"
      scrollable
      class="w-full"
    >
      <PVTabPanel
        v-if="(incompleteUploadsData?.items ?? []).length > 0"
        :pt="{'headerAction': { class: 'bg-yellow-100' }}"
      >
        <template #header>
          <div class="flex align-items-center gap-3">
            <i class="pi pi-exclamation-triangle" />
            <div>{{ tt('Incomplete Uploads') }}</div>
          </div>
        </template>
        <IncompleteuploadListView
          :incomplete-uploads="incompleteUploadsData.items"
          @refresh="refreshIncompleteUploads"
        />
      </PVTabPanel>
      <PVTabPanel>
        <template #header="slotProps">
          <div class="flex align-items-center gap-3">
            <i class="pi pi-th-large" />
            <div>{{ tt('Portfolios') }}</div>
            {{ slotProps.data }}
          </div>
        </template>
        <PortfolioListView
          v-if="portfolioData && portfolioGroupData && initiativeData"
          v-model:selected-portfolio-group-ids="selectedPortfolioGroupIds"
          v-model:selected-portfolio-ids="selectedPortfolioIds"
          :portfolios="portfolioData.items"
          :portfolio-groups="portfolioGroupData.items"
          :initiatives="initiativeData"
          @refresh="refreshAll"
        />
      </PVTabPanel>
      <PVTabPanel>
        <template #header>
          <div class="flex align-items-center gap-3">
            <i class="pi pi-table" />
            <div>{{ tt('Portfolio Groups') }}</div>
          </div>
        </template>
        <PortfolioGroupListView
          v-if="portfolioData && portfolioGroupData && initiativeData"
          v-model:selected-portfolio-group-ids="selectedPortfolioGroupIds"
          v-model:selected-portfolio-ids="selectedPortfolioIds"
          :portfolios="portfolioData.items"
          :portfolio-groups="portfolioGroupData.items"
          @refresh="refreshAll"
        />
      </PVTabPanel>
    </PVTabView>
  </StandardContent>
</template>

<style lang="scss">
.portfolios-page {
  .p-tabview .p-tabview-nav li.p-highlight .p-tabview-nav-link {
    border: 2px solid;
    padding-top: calc(1rem - 1px);
    padding-bottom: calc(1rem - 1px);
  }

  .p-datatable.p-datatable-sm {
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
}
</style>
