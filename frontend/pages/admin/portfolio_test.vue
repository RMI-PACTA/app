<script setup lang="ts">
import { NewPortfolioAsset } from 'openapi/generated/pacta';
import { FileUploadUploaderEvent } from 'primevue/fileupload';

const pactaClient = await usePACTA()
const { $axios } = useNuxtApp()

interface Asset {
  resp: NewPortfolioAsset
  fileName: string
}

const uploadedAssets = useState<Asset[]>('Index.uploadedAssets', () => [])

const uploadedAssetNames = computed(() => {
  return uploadedAssets.value.map((npa) => npa.fileName.substring(npa.fileName.lastIndexOf('/') + 1))
})

const assetIDs = computed(() => {
  return uploadedAssets.value.map((npa) => npa.resp.asset_id)
})

const startProcessing = async () => {
  if (assetIDs.value?.length === 0) {
    return
  }
  const resp = await pactaClient.processPortfolio({asset_ids: assetIDs.value})
  alert(`TASK ID: ${resp.task_id}`)
}

const createPortfolioAsset = async (file: File) => {
  const resp = await pactaClient.createPortfolioAsset()
  void await $axios({
    method: 'PUT',
    url: resp.upload_url,
    data: file,
    headers: {
      'Content-Type': file.type,
      'x-ms-blob-type': 'BlockBlob',
    },
  })
  uploadedAssets.value.push({resp, fileName: file.name})
}

const onUpload = async (e: FileUploadUploaderEvent) => {
  if (!e.files || e.files.length === 0) {
    return
  }
  const files = Array.isArray(e.files) ? e.files : [e.files]
  for (const file of files) {
    await createPortfolioAsset(file)
  }
}
</script>

<template>
  <div>
    <div class="flex gap-3 m-2 justify-content-center">
      <PVFileUpload
        mode="basic"
        :auto="true"
        :multiple="true"
        custom-upload
        choose-label="Upload Portfolio(s)"
        @uploader="onUpload"
      />

      <PVButton
        label="Send Portfolio(s) for Processing"
        icon="pi pi-send"
        :disable="uploadedAssets.length == 0"
        @click="startProcessing"
      />
    </div>

    <template v-if="uploadedAssetNames.length > 0">
      <hr>
      <ul>
        <li
          v-for="asset in uploadedAssetNames"
          :key="asset"
        >
          {{ asset }}
        </li>
      </ul>
    </template>
  </div>
</template>
