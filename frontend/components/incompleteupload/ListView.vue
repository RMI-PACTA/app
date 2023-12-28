<script setup lang="ts">
import { incompleteUploadEditor } from '@/lib/editor'
import { type IncompleteUpload } from '@/openapi/generated/pacta'
import { selectedCountSuffix } from '@/lib/selection'

const {
  humanReadableTimeFromStandardString,
} = useTime()
const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()
const i18n = useI18n()
const { t } = i18n

interface Props {
  incompleteUploads: IncompleteUpload[]
}
const props = defineProps<Props>()
interface Emits {
  (e: 'refresh'): void
}
const emit = defineEmits<Emits>()

interface EditorObject extends ReturnType<typeof incompleteUploadEditor> {
  id: string
}

const prefix = 'components/incompleteupload/ListView'
const tt = (s: string) => t(`${prefix}.${s}`)

const editorObjects = computed<EditorObject[]>(() => props.incompleteUploads.map((item) => ({ ...incompleteUploadEditor(item, i18n), id: item.id })))

const expandedRows = useState<EditorObject[]>(`${prefix}.expandedRows`, () => [])
const selectedRows = useState<EditorObject[]>(`${prefix}.selectedRows`, () => [])

const deleteIncompleteUpload = (id: string) => withLoading(
  () => pactaClient.deleteIncompleteUpload(id).then(() => {
    expandedRows.value = expandedRows.value.filter((row) => row.id !== id)
  }),
  `${prefix}.deletePortfolioGroup`,
)
const deleteSelected = async () => {
  await Promise.all(
    selectedRows.value.map((row) => deleteIncompleteUpload(row.id)),
  ).then(() => {
    emit('refresh')
    selectedRows.value = []
  })
}
const saveChanges = (id: string) => {
  const index = editorObjects.value.findIndex((editorObject) => editorObject.id === id)
  const eo = presentOrFileBug(editorObjects.value[index])
  return withLoading(
    () => pactaClient.updateIncompleteUpload(id, eo.changes.value).then(() => { emit('refresh') }),
    `${prefix}.saveChanges`,
  )
}
</script>

<template>
  <div class="flex flex-column gap-3">
    <p>
      TODO(#80) Write some copy about what Incomplete Uploads are, and direct users toward deleting them.
    </p>
    <div class="flex gap-2 flex-wrap">
      <PVButton
        icon="pi pi-refresh"
        class="p-button-outlined p-button-secondary p-button-sm"
        :label="tt('Refresh')"
        @click="() => emit('refresh')"
      />
      <PVButton
        :disabled="!selectedRows || selectedRows.length === 0"
        icon="pi pi-trash"
        class="p-button-outlined p-button-danger p-button-sm"
        :label="tt('Delete') + selectedCountSuffix(selectedRows)"
        @click="deleteSelected"
      />
    </div>
    <PVDataTable
      v-model:selection="selectedRows"
      v-model:expanded-rows="expandedRows"
      :value="editorObjects"
      data-key="id"
      class="w-full"
      size="small"
      sort-field="editorValues.value.createdAt.originalValue"
      :sort-order="-1"
    >
      <PVColumn selection-mode="multiple" />
      <PVColumn
        field="editorValues.value.name.originalValue"
        sortable
        :header="tt('Name')"
      />
      <PVColumn
        field="editorValues.value.createdAt.originalValue"
        :header="tt('Created At')"
        sortable
      >
        <template #body="slotProps">
          {{ humanReadableTimeFromStandardString(slotProps.data.editorValues.value.createdAt.originalValue).value }}
        </template>
      </PVColumn>
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
              <b>{{ humanReadableTimeFromStandardString(slotProps.data.editorValues.value.createdAt.originalValue).value }}</b>
            </div>
          </div>
          <!-- TODO(grady) add failure information here. -->
          <StandardDebug
            :value="slotProps.data.editorValues.value"
            label="Editor Values"
          />
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
              :label="tt('Delete')"
              @click="async () => { await deleteIncompleteUpload(slotProps.data.id); emit('refresh') }"
            />
            <div v-tooltip.bottom="slotProps.data.saveTooltip">
              <PVButton
                :disabled="!slotProps.data.canSave.value"
                :label="tt('Save Changes')"
                icon="pi pi-save"
                icon-pos="right"
                @click="() => saveChanges(slotProps.data.id)"
              />
            </div>
          </div>
        </div>
      </template>
    </PVDataTable>
  </div>
</template>
