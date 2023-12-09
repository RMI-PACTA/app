<script setup lang="ts">
import { type FileUploadUploaderEvent } from 'primevue/fileupload'
import { serializeError } from 'serialize-error'
import { formatFileSize } from '@/lib/filesize'

const pactaClient = usePACTA()
const { $axios } = useNuxtApp()
const { t } = useI18n()

const prefix = 'Upload'
const tt = (key: string) => t(`${prefix}.${key}`)

enum FileStatus {
  Selected = 'Selected',
  Waiting = 'Waiting',
  Uploading = 'Uploading',
  Uploaded = 'Uploaded',
  Validating = 'Validating',
  CleanUp = 'Cleaning Up',
  Done = 'Done',
  Error = 'Error',
}

interface FileState {
  file: File
  status: FileStatus
  shortName: string
  incompleteUploadId?: string
  errorMessage?: string
}

interface FileStateDetail extends FileState {
  index: number
  icon: string
  sizeStr: string
  key: string
  effectiveError?: string | undefined
}

const holdingsDate = useState<Date>(`${prefix}.holdingsDate`, () => new Date())
const errorCode = useState<string>(`${prefix}.errorCode`, () => '')
const errorMessage = useState<string>(`${prefix}.errorMessage`, () => '')
const startedProcessing = useState<boolean>(`${prefix}.startedProcessing`, () => false)
const isProcessing = useState<boolean>(`${prefix}.isProcessing`, () => false)
const fileStates = useState<FileState[]>(`${prefix}.fileState`, () => [])

const reset = () => {
  holdingsDate.value = new Date()
  errorCode.value = ''
  errorMessage.value = ''
  startedProcessing.value = false
  isProcessing.value = false
  fileStates.value = []
}

const removeFile = (index: number) => {
  fileStates.value.splice(index, 1)
}
const hasAnyState = (status: FileStatus): boolean => {
  return fileStates.value.some((fileState) => fileState.status === status)
}
const hasAllState = (status: FileStatus): boolean => {
  return fileStates.value.every((fileState) => fileState.status === status)
}
const setAllStates = (status: FileStatus) => {
  fileStates.value.forEach((fileState) => {
    fileState.status = status
  })
}

const fileStatesWithDetail = computed<FileStateDetail[]>(() => {
  const dupeKey = (fileState: FileState) => `${fileState.file.name}-${fileState.file.size}`
  const fileNameAndSizeCounts = fileStates.value.reduce((acc, fileState) => {
    const key = dupeKey(fileState)
    acc.set(key, (acc.get(key) ?? 0) + 1)
    return acc
  }, new Map<string, number>())
  const isDuplicate = (fileState: FileState) => (fileNameAndSizeCounts.get(dupeKey(fileState)) ?? 0) > 1
  return fileStates.value.map((fileState, index) => {
    let icon = 'pi pi-spinner pi-spin'
    if (fileState.status === FileStatus.Done) {
      icon = 'pi pi-check-circle text-success'
    } else if (fileState.status === FileStatus.Error) {
      icon = 'pi pi-times text-error'
    } else if (fileState.status === FileStatus.Selected) {
      icon = 'pi pi-circle'
    }
    let otherError: string | undefined
    // TODO(#79) validate this server side too.
    if (fileState.file.name.length > 1000) {
      otherError = 'Filename is too long (1000 characters max).'
    } else if (fileState.file.size > 1028 * 1028 * 100) {
      otherError = 'File is too large (100MB max).'
    } else if (!fileState.file.name.endsWith('.csv')) {
      otherError = 'File must be a csv file.'
    } else if (isDuplicate(fileState)) {
      otherError = 'This file may be a duplicate, consider removing it.'
    }
    return {
      ...fileState,
      index,
      icon,
      sizeStr: formatFileSize(fileState.file.size),
      statusStr: tt(fileState.status),
      key: `${fileState.file.name}-${index}`,
      effectiveError: fileState.errorMessage ?? otherError,
    }
  })
})
const fileUploaderProps = computed(() => ({
  disabled: isProcessing.value || allDone.value,
  mode: 'basic',
  auto: true,
  multiple: true,
  'custom-upload': true,
  'choose-label': 'Add File(s)',
}))
const actionButtonLabel = computed(() => {
  if (hasAnyState(FileStatus.Waiting)) {
    return 'Waiting...'
  }
  if (hasAnyState(FileStatus.Uploading)) {
    return 'Uploading...'
  }
  if (hasAnyState(FileStatus.Validating)) {
    return 'Validating...'
  }
  if (hasAnyState(FileStatus.CleanUp)) {
    return 'Cleaning Up...'
  }
  return 'Begin Upload'
})
const allDone = computed(() => hasAllState(FileStatus.Done) && fileStates.value.length > 0)

const onSelect = (e: FileUploadUploaderEvent) => {
  if (!e.files || e.files.length === 0) {
    return
  }
  const fs = Array.isArray(e.files) ? e.files : [e.files]
  fileStates.value.push(...fs.map((file) => ({
    file,
    fileName: file.name,
    shortName: file.name.substring(file.name.lastIndexOf('/') + 1),
    size: file.size,
    status: FileStatus.Selected,
  })))
  setAllStates(FileStatus.Selected)
}

const startUpload = async () => {
  errorCode.value = ''
  errorMessage.value = ''
  isProcessing.value = true
  startedProcessing.value = true
  setAllStates(FileStatus.Waiting)
  const startPortfolioUploadResp = await pactaClient.startPortfolioUpload({
    items: fileStates.value.map((fileState: FileState) => ({
      file_name: fileState.file.name,
      // TODO(#79) consider adding file size here as a validation step.
    })),
    holdings_date: {
      time: holdingsDate.value.toISOString(),
    },
  }).catch(e => {
    console.log('error starting upload', e, e.body)
    if (e.body?.error_id) {
      errorCode.value = e.body.error_id
    } else {
      errorCode.value = 'Unknown Error'
    }
    if (e.body?.message) {
      errorMessage.value = e.body.message
    } else {
      errorMessage.value = 'An unexpected error occurred - please file a bug to help us fix it.'
    }
    setAllStates(FileStatus.Selected)
  })
  if (!startPortfolioUploadResp) {
    isProcessing.value = false
    return
  }

  setAllStates(FileStatus.Uploading)
  const uploads = fileStatesWithDetail.value.map(async (fileState: FileStateDetail): Promise<void> => {
    const { file, index } = fileState
    const respItem = startPortfolioUploadResp.items.find((item) => item.file_name === file.name)
    if (!respItem) {
      throw new Error(`Could not find start item for file ${file.name} - something is probably wrong in the API bizlogic`)
    }
    fileStates.value[index].incompleteUploadId = respItem.incomplete_upload_id

    await $axios({
      method: 'PUT',
      url: respItem.upload_url,
      data: file,
      headers: {
        'Content-Type': file.type,
        'x-ms-blob-type': 'BlockBlob',
      },
    }).then(() => {
      fileStates.value[index].status = FileStatus.Uploaded
    }).catch((e) => {
      console.log('error uploading file', e)
      fileStates.value[index].status = FileStatus.Error
      fileStates.value[index].errorMessage = serializeError(e)
      errorCode.value = 'Upload Failed'
      errorMessage.value = 'One or more files could not be uploaded - please delete/resolve them and try again.'
    })
  })
  await Promise.all(uploads)

  let hadError = false
  for (const fileState of fileStates.value) {
    if (fileState.status === FileStatus.Uploaded) {
      continue
    } else if (fileState.status === FileStatus.Error) {
      hadError = true
    } else {
      throw new Error(`Unexpected file state ${fileState.status}`)
    }
  }
  if (!hadError) {
    await doParsing()
  }
  isProcessing.value = false
}

const doParsing = async () => {
  fileStates.value.forEach((_, i) => {
    fileStates.value[i].status = FileStatus.Validating
  })
  await pactaClient.completePortfolioUpload({
    items: fileStates.value.map((fileState) => ({
      incomplete_upload_id: presentOrFileBug(fileState.incompleteUploadId),
    })),
  })
  await waitForValidationToCompleteOrTimeout()
  await cleanUpIncompleteUploads()
}

const waitForValidationToCompleteOrTimeout = async () => {
  const timeout = 1000 * 60 * 5
  const start = Date.now()
  do {
    await refreshStateFromIncompleteUploads()
    if (Date.now() - start > timeout) {
      throw new Error('Timeout waiting for uploads to complete')
    }
    await new Promise((resolve) => setTimeout(resolve, 1000))
  } while (fileStates.value.some((fileState) => fileState.status === FileStatus.Validating))
  isProcessing.value = false
}

const refreshStateFromIncompleteUploads = async () => {
  const resp = await pactaClient.listIncompleteUploads()
  const incompleteUploads = resp.items
  for (const incompleteUpload of incompleteUploads) {
    const idx = fileStates.value.findIndex((fileState) => fileState.incompleteUploadId === incompleteUpload.id)
    // Note - this item might not be in the list if the user hasn't cleaned up prior incomplete portfolios.
    if (idx < 0) {
      continue
    }
    if (incompleteUpload.failureCode) {
      fileStates.value[idx].status = FileStatus.Error
      fileStates.value[idx].errorMessage = incompleteUpload.failureMessage
    } else if (incompleteUpload.completedAt) {
      fileStates.value[idx].status = FileStatus.CleanUp
    }
  }
}

const cleanUpIncompleteUploads = async () => {
  const fss = fileStates.value
  for (let i = 0; i < fss.length; i++) {
    const id = fss[i].incompleteUploadId
    if (id) {
      await pactaClient.deleteIncompleteUpload(id)
      fileStates.value[i].status = FileStatus.Done
    }
  }
}
</script>

<template>
  <StandardContent>
    <TitleBar title="Upload Portfolios" />
    <!-- TODO(#80) Finalize this copy -->
    <p>
      This is a page where you can upload portfolios to test out the PACTA platform.
      This Copy will need work, and will need to link to the documentation.
    </p>
    <FormField
      label="Holdings Date"
      help-text="The holdings date for the portfolios that will be uploaded"
    >
      <PVCalendar
        v-model="holdingsDate"
        view="month"
        date-format="mm/yy"
        :disabled="startedProcessing"
      />
    </FormField>
    <FormField
      label="Portfolio Files"
      class="w-full mb-0"
      help-text="This should include a link to documentation etc."
    >
      <PVFileUpload
        v-show="fileStatesWithDetail.length === 0"
        v-bind="fileUploaderProps"
        @uploader="onSelect"
      />
      <PVDataTable
        v-show="fileStatesWithDetail.length > 0"
        :value="fileStatesWithDetail"
        class="w-full"
        data-key="key"
      >
        <PVColumn>
          <template #header>
            <PVFileUpload
              v-bind="fileUploaderProps"
              @uploader="onSelect"
            />
          </template>
          <template #body="slotProps">
            <div class="flex gap-2 flex-wrap justify-content-between align-items-center">
              <div class="flex flex-column gap-2">
                <div class="font-bold">
                  {{ slotProps.data.shortName }}
                </div>
                <div class="flex gap-2 align-items-center">
                  <div>({{ slotProps.data.sizeStr }})</div>
                  <PVButton
                    class="p-button-danger p-button-text px-1 py-0 w-auto"
                    icon="pi pi-trash"
                    :disabled="isProcessing"
                    @click="() => removeFile(slotProps.data.index)"
                  />
                </div>
              </div>
              <PVMessage
                v-if="slotProps.data.effectiveError"
                severity="warn"
                :closable="false"
              >
                {{ slotProps.data.effectiveError }}
              </PVMessage>
              <div class="flex gap-2 align-items-center">
                <div><i :class="slotProps.data.icon" /></div>
                <div>{{ slotProps.data.status }}</div>
              </div>
            </div>
          </template>
        </PVColumn>
      </PVDataTable>
      <StandardDebug
        :value="fileStatesWithDetail"
        label="File States"
      />
    </FormField>
    <PVMessage
      v-show="!!errorCode"
      severity="error"
      class="m-0"
      :closable="false"
    >
      <div class="flex flex-column gap-2">
        <b>{{ errorCode }}</b>
        <div>{{ errorMessage }}</div>
      </div>
    </PVMessage>
    <PVButton
      v-if="!allDone"
      :label="actionButtonLabel"
      :loading="isProcessing"
      :disabled="isProcessing || fileStatesWithDetail.length === 0"
      icon="pi pi-arrow-right"
      icon-pos="right"
      @click="startUpload"
    />
    <template v-else>
      <PVMessage
        severity="success"
        class="m-0"
        :closable="false"
      >
        <!-- TODO(#80) Finalize This Copy -->
        Files have been uploaded, parsed, and translated to portfolios successfully.
      </PVMessage>
      <div class="flex gap-3">
        <PVButton
          label="Upload More"
          class="p-button-outlined"
          icon="pi pi-sync"
          @click="reset"
        />
        <LinkButton
          label="See Uploaded Portfolios"
          icon="pi pi-arrow-right"
          icon-pos="right"
          to="/portfolios"
        />
      </div>
    </template>
  </standardcontent>
</template>
