<script setup lang="ts">
import { portfolioGroupEditor } from '@/lib/editor'

const prefix = 'portfolio/group/NewModal'
const pactaClient = usePACTA()
const { loading: { withLoading }, newPortfolioGroup: { newPortfolioGroupVisible } } = useModal()
const i18n = useI18n()
const { t } = i18n

interface Emits {
  (e: 'created'): void
}
const emit = defineEmits<Emits>()

const defaultPortfolioGroup = {
  id: '',
  name: '',
  description: '',
  createdAt: '',
  members: [],
}
const {
  editorFields,
  editorValues,
  currentValue: portfolioGroup,
  saveTooltip,
  canSave,
} = portfolioGroupEditor(defaultPortfolioGroup, i18n)
const discard = () => { newPortfolioGroupVisible.value = false }
const save = () => withLoading(
  () => pactaClient.createPortfolioGroup(portfolioGroup.value)
    .then(() => { emit('created'); newPortfolioGroupVisible.value = false }),
  `${prefix}.save`,
)
const tt = (key: string) => t(`components/portfolio/group/NewModal.${key}`)

const header = computed(() => tt('New Portfolio Group'))
const subHeader = computed(() => tt('You can do things here like create a portfolio group'))
</script>

<template>
  <StandardModal
    v-model:visible="newPortfolioGroupVisible"
    :header="header"
    :sub-header="subHeader"
  >
    <p>
      TODO(#80) Copy Goes Here
    </p>
    <PortfolioGroupEditor
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
  </StandardModal>
</template>
