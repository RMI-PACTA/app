<script setup lang="ts">
import { type Portfolio, type AccessBlobContentResp } from '@/openapi/generated/pacta'

const { t } = useI18n()
const pactaClient = usePACTA()
const { error: { handleError } } = useModal()

interface Props {
  portfolio: Portfolio
}
const props = defineProps<Props>()

const prefix = 'components/portfolio/DownloadButton'
const statePrefix = `${prefix}[${useStateIDGenerator().id()}]`
const tt = (key: string) => t(`${prefix}.${key}`)

const downloadInProgress = useState<boolean>(`${statePrefix}.downloadInProgress`, () => false)
const doDownload = async () => {
  const blob = props.portfolio.blob
  if (!blob) {
    handleError(new Error(`Portfolio did not have blob associated with it, and could not be downloaded: ${props.portfolio.id}`))
    return
  }
  const blobId = blob.id
  downloadInProgress.value = true
  const resp: AccessBlobContentResp = await pactaClient.accessBlobContent({ items: [{ blobId }] })
  const response = await fetch(resp.items[0].downloadUrl)
  const data = await response.blob()
  const element = document.createElement('a')
  element.href = URL.createObjectURL(data)
  const fileName = `${blob.fileName}`
  element.download = fileName
  document.body.appendChild(element)
  element.click()
  document.body.removeChild(element)
  downloadInProgress.value = false
}
</script>

<template>
  <PVButton
    :loading="downloadInProgress"
    :icon="downloadInProgress ? 'pi pi-spinner pi-spin' : 'pi pi-download'"
    class="p-button-secondary p-button-outlined p-button-xs"
    :label="tt('Download')"
    @click="doDownload"
  />
</template>
