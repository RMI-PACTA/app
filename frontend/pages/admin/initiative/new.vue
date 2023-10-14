<script setup lang="ts">
import { Initiative } from '@/openapi/generated/pacta'
import { initiativeEditor } from '@/lib/editor'

const localePath = useLocalePath()
const prefix = 'admin/initiative/new'
const router = useRouter()
const pactaClient = await usePACTA()
const { loading: { withLoading } } = useModal()

const {
  editorObject: editorInitiative,
  currentValue: initiative,
  saveTooltip,
  canSave,
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
const discard = () => router.push(localePath('/admin/initiative'))
const save = () => withLoading(
  () => pactaClient.createInitiative(initiative.value).then(() => router.push(localePath('/admin/initiative'))),
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
          :disabled="!canSave"
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
