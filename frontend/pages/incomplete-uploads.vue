<script setup lang="ts">
import { incompleteUploadEditor } from '@/lib/editor'

const { humanReadableTimeFromStandardString } = useTime()
const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()
const i18n = useI18n()

const prefix = 'pages/my-data'
const expandedRows = useState(`${prefix}.expandedRows`, () => [])

const [
  { data },
] = await Promise.all([
  useSimpleAsyncData(`${prefix}.incompleteUploads`, () => pactaClient.listIncompleteUploads()),
])

let editorObjects = data.value.items.map((item) => ({ ...incompleteUploadEditor(item, i18n), id: item.id }))

const deleteIncompleteUpload = (id: string) => withLoading(
  () => pactaClient.deleteIncompleteUpload(id).then(() => {
    editorObjects = editorObjects.filter((editorObject) => editorObject.id !== id)
  }),
  `${prefix}.deleteIncompleteUpload[${id}]`,
)
const saveChanges = (id: string) => {
  const index = editorObjects.findIndex((editorObject) => editorObject.id === id)
  const eo = presentOrFileBug(editorObjects[index])
  return withLoading(
    () => pactaClient.updateIncompleteUpload(id, eo.changes.value)
      .then(() => pactaClient.findIncompleteUploadById(id))
      .then((incompleteUpload) => {
        editorObjects[index] = { ...incompleteUploadEditor(incompleteUpload, i18n), id }
      }),
    `${prefix}.saveChanges`,
  )
}
const deleteAll = () => withLoading(
  async () => {
    for (const id of editorObjects.map(eo => eo.id)) {
      await deleteIncompleteUpload(id)
    }
    editorObjects = []
  },
  `${prefix}.deleteAll`,
)
</script>

<template>
  <StandardContent>
    <TitleBar title="Incomplete Uploads" />
    <p>
      This page shows uploads that haven't been parsed successfully. This typically happens
      because they are missing required fields, have invalid values, or are otherwise not
      recognized by the system.
    </p>
    <p>
      In the row-expansion for each upload, you can see the status, failure reason (if any),
      and change properties of the file like name and description, which can be helpful for
      recordkeeping.
    </p>
    <p>
      If you want to fully abandon the upload, just delete it, and the data will be removed from
      our systems. If you'd instead like an admin to take a look at the failure to help you diagnose
      the issue, you can enable the "Admin Debugging" option, and file a support ticket.
    </p>
    <PVDataTable
      v-model:expanded-rows="expandedRows"
      size="small"
      :value="editorObjects"
      data-key="id"
      class="incomplete-upload-table"
    >
      <PVColumn
        field="currentValue.value.createdAt"
        header="Created At"
        sortable
      >
        <template #body="slotProps">
          {{ humanReadableTimeFromStandardString(slotProps.data.currentValue.value.createdAt).value }}
        </template>
      </PVColumn>
      <PVColumn
        field="currentValue.value.name"
        sortable
        header="Name"
      />
      <PVColumn
        expander
        header="Details"
      />
      <template
        #expansion="slotProps"
      >
        <div class="surface-100 p-3">
          <h2 class="mt-0">
            Metadata
          </h2>
          <div class="flex flex-column gap-2 w-fit">
            <div class="flex gap-2 justify-content-between">
              <span>Created At</span>
              <b>{{ humanReadableTimeFromStandardString(slotProps.data.currentValue.value.createdAt).value }}</b>
            </div>
            <div class="flex gap-2 justify-content-between">
              <span>Ran At</span>
              <b>{{ slotProps.data.currentValue.value.ranAt ? humanReadableTimeFromStandardString(slotProps.data.currentValue.value.ranAt).value : 'N/A' }}</b>
            </div>
            <div class="flex gap-2 justify-content-between">
              <span>Completed At</span>
              <b>{{ slotProps.data.currentValue.value.completedAt ? humanReadableTimeFromStandardString(slotProps.data.currentValue.value.completedAt).value : 'N/A' }}</b>
            </div>
            <div class="flex gap-2 justify-content-between">
              <span>Failure Code</span>
              <b>{{ slotProps.data.currentValue.value.failureCode ?? 'N/A' }}</b>
            </div>
            <div class="flex gap-2 justify-content-between">
              <span>Failure Message</span>
              <b>{{ slotProps.data.currentValue.value.failureMessage ?? 'N/A' }}</b>
            </div>
            <div class="flex gap-2 justify-content-between">
              <span>ID</span>
              <b>{{ slotProps.data.id }}</b>
            </div>
          </div>
          <h2 class="mt-5">
            Editable Properties
          </h2>
          <IncompleteuploadEditor
            v-model:editor-values="slotProps.data.editorValues.value"
            :editor-fields="slotProps.data.editorFields.value"
          />
          <div class="flex gap-3 justify-content-between">
            <PVButton
              icon="pi pi-trash"
              class="p-button-danger p-button-outlined"
              label="Delete"
              @click="() => deleteIncompleteUpload(slotProps.data.id)"
            />
            <div v-tooltip.bottom="slotProps.data.saveTooltip">
              <PVButton
                :disabled="!slotProps.data.canSave.value"
                label="Save Changes"
                icon="pi pi-save"
                icon-pos="right"
                @click="() => saveChanges(slotProps.data.id)"
              />
            </div>
          </div>
          <StandardDebug
            :value="slotProps.data.editorFields.value"
            label="Editor Fields"
          />
          <StandardDebug
            :value="slotProps.data.editorValues.value"
            label="EditorValues"
          />
        </div>
      </template>
    </PVDataTable>
    <PVButton
      v-if="editorObjects.length > 0"
      label="Delete All Incomplete Uploads"
      icon="pi pi-trash"
      class="p-button-danger p-button-outlined"
      @click="deleteAll"
    />
    <StandardDebug
      :value="data"
      label="Editor Objects"
    />
  </standardcontent>
</template>

<style lang="scss">
.incomplete-upload-table.p-datatable {
  width: 100%;

  .p-datatable-row-expansion td {
    padding: 0;
  }
}
</style>
