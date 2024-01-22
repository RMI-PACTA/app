<script setup lang="ts">
import { AnalysisType, type Analysis } from '@/openapi/generated/pacta'
import { linkToAnalysis } from '@/lib/mydata'

const { humanReadableTimeFromStandardString } = useTime()
const i18n = useI18n()
const localePath = useLocalePath()
const { t } = i18n
const prefix = 'components/analysis/ContextualListView'
const tt = (s: string) => t(`${prefix}.${s}`)

interface Props {
  analyses: Analysis[]
  name: string
  portfolioId?: string | undefined
  portfolioGroupId?: string | undefined
  initiativeId?: string | undefined
}
const props = defineProps<Props>()

interface Emits {
  (e: 'refresh'): void
}
const emit = defineEmits<Emits>()

const hasAnyAudits = computed(() => props.analyses.some((a) => a.analysisType === AnalysisType.ANALYSIS_TYPE_AUDIT))
const hasAnyReports = computed(() => props.analyses.some((a) => a.analysisType === AnalysisType.ANALYSIS_TYPE_REPORT))

const showPromptToRunAudit = computed(() => {
  return !hasAnyAudits.value && !hasAnyReports.value
})
const showPromptToRunReport = computed(() => {
  return hasAnyAudits.value && !hasAnyReports.value
})
const auditButtonClasses = computed(() => {
  return !hasAnyAudits.value && !hasAnyReports.value ? '' : 'p-button-outlined'
})
const reportButtonClasses = computed(() => {
  return !hasAnyReports.value && hasAnyAudits.value ? '' : 'p-button-outlined'
})
</script>

<template>
  <div class="flex flex-column gap-2">
    <PVDataTable
      :value="props.analyses"
      data-key="id"
      size="medium"
      sort-field="createdAt"
      :sort-order="-1"
    >
      <template #empty>
        <div class="p-3">
          {{ tt('No Analyses Message') }}
        </div>
      </template>
      <PVColumn
        field="createdAt"
        :header="tt('Ran At')"
      >
        <template #body="slotProps">
          {{ humanReadableTimeFromStandardString(slotProps.data.createdAt).value }}
        </template>
      </PVColumn>
      <PVColumn
        :header="tt('Status')"
      >
        <template #body="slotProps">
          <AnalysisStatusChip
            :analysis="slotProps.data"
          />
        </template>
      </PVColumn>
      <PVColumn
        field="analysisType"
        :header="tt('Type')"
      >
        <template #body="slotProps">
          {{ tt(slotProps.data.analysisType) }}
        </template>
      </PVColumn>
      <PVColumn
        field="name"
        :header="tt('Name')"
      />
      <PVColumn
        field="Access"
        :header="tt('Access')"
      >
        <template #body="slotProps">
          <AnalysisAccessButtons
            class="py-2"
            :analysis="slotProps.data"
          />
        </template>
      </PVColumn>
      <PVColumn :header="tt('Details')">
        <template #body="slotProps">
          <LinkButton
            class="p-button-outlined p-button-xs p-button-secondary"
            icon="pi pi-arrow-right"
            :to="linkToAnalysis(localePath, slotProps.data.id)"
          />
        </template>
      </PVColumn>
    </PVDataTable>
    <PVMessage
      v-if="showPromptToRunAudit"
      severity="info"
      class="m-0"
      :closable="false"
    >
      {{ tt('No Audits Message') }}
    </PVMessage>
    <PVMessage
      v-if="showPromptToRunReport"
      severity="info"
      class="m-0"
      :closable="false"
    >
      {{ tt('No Reports Message') }}
    </PVMessage>
    <div class="flex gap-2">
      <AnalysisRunButton
        :analysis-type="AnalysisType.ANALYSIS_TYPE_AUDIT"
        :name="props.name"
        :class="auditButtonClasses"
        :portfolio-id="props.portfolioId"
        :portfolio-group-id="props.portfolioGroupId"
        :initiative-id="props.initiativeId"
        class="p-button-sm"
        @started="() => emit('refresh')"
      />
      <AnalysisRunButton
        :analysis-type="AnalysisType.ANALYSIS_TYPE_REPORT"
        :name="props.name"
        :class="reportButtonClasses"
        :portfolio-id="props.portfolioId"
        :portfolio-group-id="props.portfolioGroupId"
        :initiative-id="props.initiativeId"
        class="p-button-sm"
        @started="() => emit('refresh')"
      />
    </div>
  </div>
</template>
