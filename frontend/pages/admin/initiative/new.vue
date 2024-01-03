<script setup lang="ts">
import { Language } from '@/openapi/generated/pacta'
import { initiativeEditor } from '@/lib/editor'

const localePath = useLocalePath()
const prefix = 'admin/initiative/new'
const router = useRouter()
const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()
const i18n = useI18n()
const { t } = i18n

const defaultInitiative = {
  id: '',
  name: '',
  affiliation: '',
  publicDescription: '',
  internalDescription: '',
  requiresInvitationToJoin: false,
  isAcceptingNewMembers: false,
  isAcceptingNewPortfolios: false,
  language: Language.LANGUAGE_EN,
  pactaVersion: undefined,
  createdAt: '',
  portfolioInitiativeMemberships: [],
}
const {
  editorFields,
  editorValues,
  currentValue: initiative,
  saveTooltip,
  canSave,
} = initiativeEditor(defaultInitiative, i18n)
const discard = () => router.push(localePath('/admin/initiative'))
const save = () => withLoading(
  () => pactaClient.createInitiative(initiative.value).then(() => router.push(localePath('/admin/initiative'))),
  `${prefix}.save`,
)
const tt = (key: string) => t(`pages/admin/initiative/new.${key}`)
</script>

<template>
  <StandardContent>
    <TitleBar :title="tt('New Initiative')" />
    <p>
      TODO(#80) Initiative Copy Goes Here
    </p>
    <InitiativeEditor
      v-model:editorValues="editorValues"
      :editor-fields="editorFields"
    />
    <div class="flex gap-3">
      <PVButton
        :label="tt('Discard')"
        icon="pi pi-arrow-left"
        class="p-button-secondary p-button-outlined"
        @click="discard"
      />
      <div v-tooltip.bottom="saveTooltip">
        <PVButton
          :disabled="!canSave"
          :label="tt('Save')"
          icon="pi pi-arrow-right"
          icon-pos="right"
          @click="save"
        />
      </div>
    </div>
    <StandardDebug
      :value="editorFields"
      label="Editor Fields"
    />
    <StandardDebug
      :value="editorValues"
      label="Editor Values"
    />
  </StandardContent>
</template>
