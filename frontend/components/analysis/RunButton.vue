<script setup lang="ts">
import { type RunAnalysisReq, type Analysis, type AnalysisType } from '@/openapi/generated/pacta'
import { useConfirm } from 'primevue/useconfirm'

const { linkToAnalysis } = useMyDataURLs()
const { require: confirm } = useConfirm()
const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()
const i18n = useI18n()
const { t } = i18n

const prefix = 'components/analysis/RunButton'
const statePrefix = `${prefix}[${useStateIDGenerator().id()}]`
const tt = (key: string) => t(`${prefix}.${key}`)

interface Props {
  analysisType: AnalysisType
  name: string
  warnForDuplicate?: boolean
  portfolioGroupId?: string
  portfolioId?: string
  initiativeId?: string
}
const props = defineProps<Props>()
interface Emits {
  (e: 'started'): void
  (e: 'finished'): void
}
const emit = defineEmits<Emits>()

const clicked = useState<boolean>(`${statePrefix}.clicked`, () => false)
const analysisId = useState<string | null>(`${statePrefix}.analysisId`, () => null)
const analysis = useState<Analysis | null>(`${statePrefix}.analysis`, () => null)

const request = computed<RunAnalysisReq>(() => {
  const common = {
    analysisType: props.analysisType,
    name: `${tt(props.analysisType)}: ${props.name}`,
    description: `${tt(props.analysisType)} run at ${new Date().toLocaleString()}`,
  }
  if (props.portfolioId) {
    return {
      ...common,
      portfolioId: props.portfolioId,
    }
  } else if (props.portfolioGroupId) {
    return {
      ...common,
      portfolioGroupId: props.portfolioGroupId,
    }
  } else if (props.initiativeId) {
    return {
      ...common,
      initiativeId: props.initiativeId,
    }
  } else {
    throw new Error('No portfolio, portfolio group or initiative ID provided')
  }
})

const startRun = async () => {
  emit('started')
  clicked.value = true
  const resp = await withLoading(
    () => pactaClient.runAnalysis(request.value),
    `${prefix}.runAnalysis`,
  )
  analysisId.value = resp.analysisId
  void refreshAnalysisState()
}
const runAnalysis = async () => {
  if (!props.warnForDuplicate) {
    await startRun()
    return
  }
  confirm({
    header: tt('ConfirmationHeader'),
    message: tt('ConfirmationMessage'),
    icon: 'pi pi-copy',
    position: 'center',
    blockScroll: true,
    reject: () => { clicked.value = false },
    rejectLabel: tt('Cancel Run'),
    rejectIcon: 'pi pi-times',
    acceptLabel: tt('Run Anyway'),
    accept: startRun,
    acceptIcon: 'pi pi-check',
  })
}
const refreshAnalysisState = async () => {
  const aid = analysisId.value
  if (!aid) {
    console.warn('No analysis ID set, but refresh requested')
    return
  }
  const resp = await pactaClient.findAnalysisById(aid)
  analysis.value = resp
  if (analysis.value?.completedAt) {
    emit('finished')
    return
  }
  setTimeout(() => { void refreshAnalysisState }, 2000)
}
const analysisCompleted = computed(() => analysis.value?.completedAt)
const runBtnVisible = computed(() => !analysisCompleted.value)
const runBtnDisabled = computed(() => analysisId.value !== null || clicked.value)
const runBtnLoading = computed(() => runBtnDisabled.value)
const runBtnIcon = computed(() => analysisId.value ? 'pi pi-spin pi-spinner' : 'pi pi-play')
const runBtnLabel = computed(() => {
  if (!clicked.value) {
    return tt('Run') + ' ' + tt(props.analysisType)
  }
  if (!analysisId.value) {
    return tt('Starting') + ' ' + tt(props.analysisType) + '...'
  }
  if (analysis.value) {
    return tt('Running') + ' ' + tt(props.analysisType) + '...'
  }
  return 'Should not happen'
})
const completeBtnTo = computed(() => {
  if (!analysisId.value) {
    return ''
  }
  return linkToAnalysis(analysisId.value)
})
const completeBtnLabel = computed(() => {
  if (!analysisId.value) {
    return ''
  }
  return tt(props.analysisType) + ' ' + tt('Completed')
})
</script>

<template>
  <PVButton
    v-if="runBtnVisible"
    :disabled="runBtnDisabled"
    :loading="runBtnLoading"
    :icon="runBtnIcon"
    :label="runBtnLabel"
    @click="runAnalysis"
  />
  <LinkButton
    v-else
    :to="completeBtnTo"
    icon="pi pi-check"
    :label="completeBtnLabel"
  />
</template>
