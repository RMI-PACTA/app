<script setup lang="ts">
const { fromQueryReactive } = useURLParams()

const { t } = useI18n()
const tt = (s: string) => t(`pages/portfolios.${s}`)

const selectedPortfolioIdsQP = fromQueryReactive('pids')
const selectedPortfolioGroupIdsQP = fromQueryReactive('pgids')
const activeIndexQP = fromQueryReactive('activeIndex')

const selectedPortfolioIds = computed<string[]>({
  get: () => (selectedPortfolioIdsQP.value ?? '').split(','),
  set: (v: string[]) => {
    if (v.length === 0) {
      selectedPortfolioIdsQP.value = undefined
    } else {
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
      selectedPortfolioGroupIdsQP.value = v.join(',')
    }
  },
})
const activeIndex = computed<number>({
  get: () => parseInt(activeIndexQP.value ?? '0'),
  set: (v: number) => {
    activeIndexQP.value = v.toString()
  },
})
</script>

<template>
  <StandardContent>
    <TitleBar :title="tt('Portfolios')" />
    <p>
      TODO(#80) Add Copy Here
    </p>
    <PVTabView v-model:activeIndex="activeIndex">
      <PVTabPanel :header="tt('Portfolios')">
        <PortfolioListView
          v-model:selected-portfolio-group-ids="selectedPortfolioGroupIds"
          v-model:selected-portfolio-ids="selectedPortfolioIds"
        />
      </PVTabPanel>
      <PVTabPanel :header="tt('Portfolio Groups')">
        <PortfolioGroupListView
          v-model:selected-portfolio-group-ids="selectedPortfolioGroupIds"
          v-model:selected-portfolio-ids="selectedPortfolioIds"
        />
      </PVTabPanel>
    </PVTabView>
  </StandardContent>
</template>

<style lang="scss">
.portfolio-upload-table.p-datatable.p-datatable-sm {
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
</style>
