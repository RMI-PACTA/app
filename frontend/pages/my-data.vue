<script setup lang="ts">
import { type WritableComputedRef } from 'vue'
import {
  Tab,
  QueryParamTab,
  QueryParamSelectedPortfolioIds,
  QueryParamExpandedPortfolioIds,
  QueryParamSelectedPortfolioGroupIds,
  QueryParamExpandedPortfolioGroupIds,
  QueryParamSelectedAnalysisIds,
  QueryParamExpandedAnalysisIds,
} from '@/lib/mydata'

const prefix = 'pages/my-data'

const { fromQueryReactiveWithDefault } = useURLParams()
const { t } = useI18n()
const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()

const tt = (s: string) => t(`${prefix}.${s}`)

const joinedWithCommas = (r: WritableComputedRef<string>): WritableComputedRef<string[]> => computed({
  get: () => r.value.split(','),
  set: (v: string[]) => { r.value = v.join(',') },
})
const selectedPortfolioIds = joinedWithCommas(fromQueryReactiveWithDefault(QueryParamSelectedPortfolioIds, ''))
const expandedPortfolioIds = joinedWithCommas(fromQueryReactiveWithDefault(QueryParamExpandedPortfolioIds, ''))
const selectedPortfolioGroupIds = joinedWithCommas(fromQueryReactiveWithDefault(QueryParamSelectedPortfolioGroupIds, ''))
const expandedPortfolioGroupIds = joinedWithCommas(fromQueryReactiveWithDefault(QueryParamExpandedPortfolioGroupIds, ''))
const selectedAnalysisIds = joinedWithCommas(fromQueryReactiveWithDefault(QueryParamSelectedAnalysisIds, ''))
const expandedAnalysisIds = joinedWithCommas(fromQueryReactiveWithDefault(QueryParamExpandedAnalysisIds, ''))
const tabQP = fromQueryReactiveWithDefault(QueryParamTab, 'p')

const [
  { data: incompleteUploadsData, refresh: refreshIncompleteUploadsApi },
  { data: portfolioData, refresh: refreshPortfoliosApi },
  { data: portfolioGroupData, refresh: refreshPortfolioGroupsApi },
  { data: analysesData, refresh: refreshAnalysesApi },
  { data: initiativeData, refresh: refreshInitiativesApi },
] = await Promise.all([
  useSimpleAsyncData(`${prefix}.incompleteUploads`, () => pactaClient.listIncompleteUploads()),
  useSimpleAsyncData(`${prefix}.portfolios`, () => pactaClient.listPortfolios()),
  useSimpleAsyncData(`${prefix}.portfolioGroups`, () => pactaClient.listPortfolioGroups()),
  useSimpleAsyncData(`${prefix}.analyses`, () => pactaClient.listAnalyses()),
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
const refreshAnalyses = async () => {
  await withLoading(refreshAnalysesApi, `${prefix}.refreshAnalyses`)
}
const refreshInitiatives = async () => {
  await withLoading(refreshInitiativesApi, `${prefix}.refreshInitiatives`)
}
const refreshAll = () => Promise.all([
  refreshIncompleteUploads(),
  refreshPortfolios(),
  refreshPortfolioGroups(),
  refreshAnalyses(),
  refreshInitiatives(),
])

type TabToIndexMap = Record<Tab, number>
const tabToIndexMap = computed(() => {
  const result: TabToIndexMap = {
    [Tab.Portfolio]: -1,
    [Tab.PortfolioGroup]: -1,
    [Tab.IncompleteUpload]: -1,
    [Tab.Analysis]: -1,
  }
  let idx = 0
  if (incompleteUploadsData.value.items.length > 0) {
    result[Tab.IncompleteUpload] = idx
    idx++
  }
  result[Tab.Portfolio] = idx++
  result[Tab.PortfolioGroup] = idx++
  if (analysesData.value.items.length > 0) {
    result[Tab.Analysis] = idx
    idx++
  }
  return result
})
const activeIndex = computed<number>({
  get: () => {
    const tab = tabQP.value
    const result = tabToIndexMap.value[tab as Tab]
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
          v-model:selected-portfolio-ids="selectedPortfolioIds"
          v-model:expanded-portfolio-ids="expandedPortfolioIds"
          :portfolios="portfolioData.items"
          :portfolio-groups="portfolioGroupData.items"
          :analyses="analysesData.items"
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
          v-if="portfolioData && portfolioGroupData && initiativeData && analysesData"
          v-model:selected-portfolio-group-ids="selectedPortfolioGroupIds"
          v-model:expanded-portfolio-group-ids="expandedPortfolioGroupIds"
          :portfolios="portfolioData.items"
          :portfolio-groups="portfolioGroupData.items"
          :analyses="analysesData.items"
          @refresh="refreshAll"
        />
      </PVTabPanel>
      <PVTabPanel
        v-if="(analysesData?.items ?? []).length > 0"
      >
        <template #header>
          <div class="flex align-items-center gap-3">
            <i class="pi pi-book" />
            <div>{{ tt('Analyses') }}</div>
          </div>
        </template>
        <AnalysisListView
          v-if="portfolioData && portfolioGroupData && initiativeData && analysesData"
          v-model:selected-analysis-ids="selectedAnalysisIds"
          v-model:expanded-analysis-ids="expandedAnalysisIds"
          :portfolios="portfolioData.items"
          :portfolio-groups="portfolioGroupData.items"
          :analyses="analysesData.items"
          :initiatives="initiativeData"
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
