<script setup lang="ts">
import { Blob } from '@/openapi/generated/pacta'
import { computed } from 'vue'

const { t } = useI18n()
const pactaClient = usePACTA()

interface Props {
  blobs: Blob[]
  cta: string
}
const props = defineProps<Props>()

const prefix = 'components/download/BlobButton'
const tt = (key: string) => t(`${prefix}.${key}`)

enum Stage {
  Inactive,
  GettingURLs,
  Downloading,
  Done,
  Error,
}

const statePrefix = `${prefix}[${useStateIDGenerator().id()}]`
const stage = useState<Stage>(`${statePrefix}.stage`, () => Stage.Inactive)
const downloadingPercentages = useState<number[]>(`${statePrefix}.downloadingPercentages`, () => [])

const message = computed(() => downloaded.value ? tt('Downloaded') : props.cta)
const disabled = computed(() => stage.value !== Stage.Inactive)
const downloadingPercentage = computed(() => {
  if (downloadingPercentages.value.length === 0) {
    return 0
  }
  return downloadingPercentages.value.reduce((a, b) => a + b, 0) / downloadingPercentages.value.length
})

const progressPath = computed(() => {
  const radius = 45
  const rotate = -90
  const startX = 50
  const startY = 50 - radius
  const endAngle = rotate - 360 * downloadingPercentage.value / 100
  const endX = 50 + radius * Math.cos(endAngle * Math.PI / 180)
  const endY = 50 + radius * Math.sin(endAngle * Math.PI / 180)
  return `M ${startX} ${startY} A ${radius}${radius} 0 0 1 0 ${endX} ${endY}`
})

interface BlobWithUrl extends Blob {
  url: string
}

const download = async () => {
  stage.value = Stage.GettingURLs
  downloadingPercentages.value = props.blobs.map(() => 0)
  const abcResp = await pactaClient.accessBlobContent({
    items: props.blobs.map(b => ({ blobId: b.id })),
  })
  const blobsWithUrls: BlobWithUrl[] = props.blobs.map(b => ({ ...b, url: '' }))
  for (const item of abcResp.items) {
    presentOrFileBug(blobsWithUrls.find(b => b.id === item.blobId)).url = item.url
  }
  stage.value = Stage.Downloading
  const downloadPromises = blobsWithUrls.map((b, i) => downloadFile(b.url, b.name, (percentage) => {
    downloadingPercentages.value[i] = percentage
  }))
  await Promise.all(downloadPromises)
  stage.value = Stage.Done
}

const downloadFile = async (url: string, filename: string, progressCallback: (percentage: number) => void): Promise<void> => {
  await new Promise<void>((resolve, reject) => {
    const xhr = new XMLHttpRequest()
    xhr.responseType = 'blob'
    xhr.onprogress = (event) => {
      if (event.lengthComputable) {
        const percentage = Math.round((event.loaded / event.total) * 100)
        progressCallback(percentage)
      }
    }
    xhr.onload = () => {
      if (xhr.status === 200) {
        const blob = new Blob([xhr.response], { type: 'application/octet-stream' })
        const downloadUrl = URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = downloadUrl
        a.download = filename
        document.body.appendChild(a)
        a.click()
        URL.revokeObjectURL(downloadUrl)
        a.remove()
        resolve()
      } else {
        reject(new Error(`Download failed: ${xhr.statusText}`))
      }
    }
    xhr.onerror = (err) => {
      reject(err)
    }
    xhr.open('GET', url)
    xhr.send()
  })
}
</script>

<template>
  <PVButton
    :disabled="disabled"
    class="text-sm"
    @click="download"
  >
    <div class="flex gap-2">
      <i
        v-if="stage === Stage.Inactive"
        class="pi pi-download"
      />
      <i
        v-else-if="stage === Stage.GettingURLs"
        class="pi pi-spin pi-spinner"
      />
      <svg
        v-else-if="stage === Stage.Downloading"
        width="100"
        height="100"
        viewBox="0 0 100 100"
      >
        <path
          :d="progressPath"
          fill="none"
          stroke="green"
          stroke-width="10"
          stroke-linecap="round"
        />
      </svg>
      <i
        v-else-if="stage === Stage.Done"
        class="pi pi-check"
      />
      <span>{{ message }}<span v-if="props.blobs.length > 1">({{ props.blobs.length }})</span></span>
    </div>
  </PVButton>
</template>
