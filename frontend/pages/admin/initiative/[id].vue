<script setup lang="ts">
import { initiativeEditor } from '@/lib/editor'

const router = useRouter()
const pactaClient = await usePACTA()
const { loading: { withLoading } } = useModal()
const { fromParams } = useURLParams()

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
    .then(() => router.push('/admin/initiative')),
  `${prefix}.deleteInitiative`,
)
const saveChanges = () => withLoading(
  () => pactaClient.updateInitiative(id, changes.value)
    .then(() => router.push('/admin/initiative')),
  `${prefix}.saveChanges`,
)
</script>

<template>
  <StandardContent v-if="editorInitiative">
    <TitleBar :title="`Editing Initiative: ${editorInitiative.name.currentValue}`" />
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
        to="/admin/initiative"
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
      label="Editor Initiative"
    />
    <StandardDebug
      :value="changes"
      label="Initiative Changes"
    />
  </StandardContent>
</template>
