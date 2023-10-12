<script setup lang="ts">
import { Initiative } from '@/openapi/generated/pacta'
import { initiativeEditor } from '@/lib/editor'

const prefix = 'admin/initiative/new'
const router = useRouter()
const pactaClient = await usePACTA()
const { loading: { withLoading } } = useModal()

const {
  editorInitiative,
  incompleteFields,
  hasChanges,
  isIncomplete,
  initiative,
} = initiativeEditor({
  id: '',
  name: '',
  affiliation: '',
  publicDescription: '',
  internalDescription: '',
  requiresInvitationToJoin: false,
  isAcceptingNewMembers: false,
  isAcceptingNewPortfolios: false,
  language: Initiative.language.EN,
  pactaVersion: undefined,
  createdAt: '',
})
const saveTooltip = computed<string | undefined>(() => {
  if (!hasChanges.value) { return 'All changes saved' }
  if (isIncomplete.value) { return `Cannot save with incomplete fields: ${incompleteFields.value.join(', ')}` }
  return undefined
})
const saveDisabled = computed<boolean>(() => saveTooltip.value !== undefined)

const discard = () => router.push('/admin/initiative')
const save = () => withLoading(
  () => pactaClient.createInitiative(initiative.value).then(() => router.push('/admin/initiative')),
  `${prefix}.save`,
)
</script>

<template>
  <StandardContent>
    <TitleBar title="New Initiative" />
    <p>
      TODO(#38) Initiative Copy Goes Here
    </p>
    <InitiativeEditor
      v-model:editorInitiative="editorInitiative"
    />
    <div class="flex gap-3">
      <PVButton
        label="Discard"
        icon="pi pi-arrow-left"
        class="p-button-secondary p-button-outlined"
        @click="discard"
      />
      <div v-tooltip.bottom="saveTooltip">
        <PVButton
          :disabled="saveDisabled"
          label="Save"
          icon="pi pi-arrow-right"
          icon-pos="right"
          @click="save"
        />
      </div>
    </div>
    <StandardDebug
      label="Editor Initiative"
      :value="editorInitiative"
    />
  </StandardContent>
</template>
