<script setup lang="ts">
import { initiativeEditor } from '@/lib/editor'

const router = useRouter()
const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()
const { fromParams } = useURLParams()
const localePath = useLocalePath()

const id = presentOrCheckURL(fromParams('id'))

const prefix = `admin/initiative/${id}`
const { data } = await useSimpleAsyncData(`${prefix}.getInitiative`, () => pactaClient.findInitiativeById(id))
const {
  editorObject: editorInitiative,
  changes,
  saveTooltip,
  canSave,
} = initiativeEditor(presentOrCheckURL(data.value, 'no initiative in response'))

const deleteInitiative = () => withLoading(
  () => pactaClient.deleteInitiative(id)
    .then(() => router.push(localePath('/admin/initiative'))),
  `${prefix}.deleteInitiative`,
)
const saveChanges = () => withLoading(
  () => pactaClient.updateInitiative(id, changes.value)
    .then(() => router.push(localePath(`/initiative/${id}`))),
  `${prefix}.saveChanges`,
)
</script>

<template>
  <div class="flex flex-column gap-3">
    <InitiativeEditor
      v-model:editorInitiative="editorInitiative"
    />
    <div class="flex gap-3">
      <PVButton
        icon="pi pi-trash"
        class="p-button-danger"
        label="Delete"
        @click="deleteInitiative"
      />
      <LinkButton
        label="Discard Changes"
        icon="pi pi-arrow-left"
        class="p-button-secondary p-button-outlined"
        :to="localePath('/admin/initiative')"
      />
      <div v-tooltip.bottom="saveTooltip">
        <PVButton
          :disabled="!canSave"
          label="Save Changes"
          icon="pi pi-arrow-right"
          icon-pos="right"
          @click="saveChanges"
        />
      </div>
    </div>
    <StandardDebug
      :value="editorInitiative"
      label="Edit Initiative"
    />
    <StandardDebug
      :value="changes"
      label="Edit Initiative Changes"
    />
  </div>
</template>
