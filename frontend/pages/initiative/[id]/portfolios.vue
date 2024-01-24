<script setup lang="ts">
import { type PortfolioInitiativeMembershipPortfolio, type InitiativeAllData } from '@/openapi/generated/pacta'
import JSZip from 'jszip'

const { fromParams } = useURLParams()
const id = presentOrCheckURL(fromParams('id'))
const localePath = useLocalePath()
const { humanReadableTimeFromStandardString } = useTime()
const { t } = useI18n()
const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()
const { initiative, refreshInitiative } = await useInitiativeData(id)

const prefix = 'pages/initiative/id/portfolios'
const tt = (key: string) => t(`${prefix}.${key}`)

const expandedRows = useState<PortfolioInitiativeMembershipPortfolio[]>(`${prefix}.expandedRows`, () => [])
const members = computed(() => initiative.value.portfolioInitiativeMemberships)

const refresh = () => { void refreshInitiative() }

const deletePortfolioMembership = (portfolioId: string) => withLoading(
  async () => {
    await pactaClient.deleteInitiativePortfolioRelationship(id, portfolioId)
    await refreshInitiative()
  },
`${prefix}.deletePortfolioMembership`)

const downloadAll = async () => {
  await withLoading(async () => {
    const iad: InitiativeAllData = await pactaClient.allInitiativeData(id)
    const zip = new JSZip()
    await Promise.all(iad.items.map(async (item) => {
      const response = await fetch(item.downloadUrl)
      const data = await response.blob()
      const fileName = `${item.name}`
      zip.file(fileName, data)
    }))
    const content = await zip.generateAsync({ type: 'blob' })
    const element = document.createElement('a')
    element.href = URL.createObjectURL(content)
    const fileName = `${initiative.value.name} - All Data - ${new Date().toISOString()}.zip`
    element.download = fileName
    document.body.appendChild(element)
    element.click()
    document.body.removeChild(element)
  }, `${prefix}.downloadAll`)
}
</script>

<template>
  <div class="flex flex-column gap-3">
    <div class="flex gap-2">
      <PVButton
        icon="pi pi-refresh"
        class="p-button-outlined p-button-secondary p-button-sm"
        :label="tt('Refresh')"
        @click="refresh"
      />
      <PVButton
        icon="pi pi-download"
        :label="tt('Download All')"
        class="p-button-outlined p-button-secondary p-button-sm"
        @click="downloadAll"
      />
    </div>
    <PVDataTable
      v-model:expanded-rows="expandedRows"
      :value="members"
      class="align-self-stretch"
    >
      <PVColumn
        field="createdAt"
        sortable
        :header="tt('Added At')"
      >
        <template #body="slotProps">
          {{ humanReadableTimeFromStandardString(slotProps.data.createdAt).value }}
        </template>
      </PVColumn>
      <PVColumn
        field="addedByUserId"
        sortable
        :header="tt('Added By User')"
      >
        <template #body="slotProps">
          <div
            v-if="slotProps.data.addedByUserId"
            class="flex flex-column gap-1"
          >
            {{ slotProps.data.addedByUserId }}
            <LinkButton
              :label="tt('View User')"
              class="p-button-xs p-button-outlined"
              :to="localePath(`/user/${slotProps.data.addedByUserId}`)"
              icon="pi pi-external-link"
              icon-pos="right"
            />
          </div>
          <div v-else>
            User is deleted
          </div>
        </template>
      </PVColumn>
      <PVColumn
        field="portfolio.id"
        sortable
        :header="tt('Portfolio')"
      >
        <template #body="slotProps">
          <div class="flex flex-column gap-1">
            <span class="text-lg font-bold">{{ slotProps.data.portfolio.name }}</span>
            <span class="font-light">{{ slotProps.data.portfolio.id }}</span>
          </div>
        </template>
      </PVColumn>
      <PVColumn :header="tt('Remove')">
        <template #body="slotProps">
          <PVButton
            class="p-button-xs p-button-danger p-button-outlined"
            :label="tt('Remove')"
            icon="pi pi-trash"
            @click="deletePortfolioMembership(slotProps.data.portfolio.id)"
          />
        </template>
      </PVColumn>
      <PVColumn :header="tt('Download')">
        <template #body="slotProps">
          <PortfolioDownloadButton :portfolio="slotProps.data.portfolio" />
        </template>
      </PVColumn>
      <PVColumn
        expander
        :header="tt('Details')"
      />
      <template #expansion="slotProps">
        <StandardDebug
          :value="slotProps.data"
          :label="slotProps.data.portfolio.name"
          always
        />
      </template>
    </PVDataTable>
    <StandardDebug :value="initiative" />
  </div>
</template>
