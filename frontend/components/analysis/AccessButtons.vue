<script setup lang="ts">
import { type Analysis, type AccessBlobContentReqItem, type AccessBlobContentResp, AnalysisType } from '@/openapi/generated/pacta'
import JSZip from 'jszip'

const { t } = useI18n()
const { public: { apiServerURL } } = useRuntimeConfig()
const pactaClient = usePACTA()
const { getMaybeMe } = useSession()
const { isAdmin, isSuperAdmin, maybeMeOwnerId } = await getMaybeMe()

interface Props {
  analysis: Analysis
}
const props = defineProps<Props>()

const prefix = 'components/analysis/AccessButtons'
const statePrefix = `${prefix}[${useStateIDGenerator().id()}]`
const tt = (key: string) => t(`${prefix}.${key}`)

const canAccessAsPublic = computed(() => props.analysis.artifacts.every((asset) => asset.sharedToPublic))
const canAccessAsAdmin = computed(() => {
  if (isAdmin.value || isSuperAdmin.value) {
    return props.analysis.artifacts.every((asset) => asset.adminDebugEnabled)
  }
  return false
})
const canAccessAsOwner = computed(() => {
  if (maybeMeOwnerId.value) {
    return maybeMeOwnerId.value === props.analysis.ownerId
  }
  return false
})
const canAccess = computed(() => {
  return canAccessAsPublic.value || canAccessAsAdmin.value || canAccessAsOwner.value
})
const isCompleted = computed(() => {
  return props.analysis.completedAt !== undefined && props.analysis.failureMessage === undefined
})
const isViewable = computed(() => {
  return isCompleted.value && (props.analysis.analysisType === AnalysisType.ANALYSIS_TYPE_DASHBOARD || props.analysis.analysisType === AnalysisType.ANALYSIS_TYPE_REPORT)
})
const downloadInProgress = useState<boolean>(`${statePrefix}.downloadInProgress`, () => false)
const doDownload = async () => {
  downloadInProgress.value = true
  // Just download the audit_file.csv for audits
  const artifactsToDownload = props.analysis.analysisType === AnalysisType.ANALYSIS_TYPE_AUDIT ? props.analysis.artifacts.filter((a) => a.blob.fileName === 'audit_file.csv') : props.analysis.artifacts
  const response: AccessBlobContentResp = await pactaClient.accessBlobContent({
    items: artifactsToDownload.map((asset): AccessBlobContentReqItem => ({
      blobId: asset.blob.id,
    })),
  })

  let content: Blob | null = null
  let fileName: string | null = null
  if (response.items.length === 1) {
    const resp = await fetch(response.items[0].downloadUrl)
    content = await resp.blob()
    fileName = artifactsToDownload[0].blob.fileName
  } else {
    const zip = new JSZip()
    await Promise.all(response.items.map(
      async (item): Promise<void> => {
        const response = await fetch(item.downloadUrl)
        const data = await response.blob()
        const blob = presentOrFileBug(props.analysis.artifacts.find((artifact) => artifact.blob.id === item.blobId)).blob
        const fileName = `${blob.fileName}`
        zip.file(fileName, data)
      }),
    )
    content = await zip.generateAsync({ type: 'blob' })
    fileName = `${props.analysis.name}.zip`
  }
  const element = document.createElement('a')
  element.href = URL.createObjectURL(content)
  element.download = fileName
  document.body.appendChild(element)
  element.click()
  document.body.removeChild(element)
  downloadInProgress.value = false
}

const openReport = () => navigateTo(`${apiServerURL}/report/${props.analysis.id}/`, {
  open: {
    target: '_blank',
  },
  external: true,
})
</script>

<template>
  <div
    v-tooltip="canAccess ? undefined : tt('Denied')"
    class="flex gap-1 align-items-center w-fit"
  >
    <PVButton
      v-if="isViewable"
      icon="pi pi-external-link"
      :disabled="!canAccess"
      class="p-button-secondary p-button-outlined p-button-xs"
      :label="tt('View')"
      @click="openReport"
    />
    <PVButton
      v-if="isCompleted"
      v-tooltip="canAccess ? tt('Download') : ''"
      :disabled="downloadInProgress || !canAccess"
      :loading="downloadInProgress"
      :icon="downloadInProgress ? 'pi pi-spinner pi-spin' : 'pi pi-download'"
      class="p-button-secondary p-button-text p-button-xs"
      @click="doDownload"
    />
  </div>
</template>
