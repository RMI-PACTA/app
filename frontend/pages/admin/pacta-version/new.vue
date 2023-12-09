<script setup lang="ts">
import { pactaVersionEditor } from '@/lib/editor'

const prefix = 'admin/pacta-version/new'
const router = useRouter()
const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()
const localePath = useLocalePath()
const i18n = useI18n()

const defaultPactaVersion = {
  id: '',
  name: '',
  description: '',
  digest: '',
  createdAt: '',
  isDefault: false,
}
const {
  editorFields,
  editorValues,
  currentValue: pactaVersion,
  canSave,
  saveTooltip,
} = pactaVersionEditor(defaultPactaVersion, i18n)
const discard = () => router.push(localePath('/admin/pacta-version'))
const save = () => withLoading(
  () => pactaClient.createPactaVersion(pactaVersion.value).then(() => router.push(localePath('/admin/pacta-version'))),
  `${prefix}.save`,
)
</script>

<template>
  <StandardContent>
    <TitleBar title="New PACTA Version" />
    <p>
      Pacta version info goes here
    </p>
    <PactaversionEditor
      v-model:editorValues="editorValues"
      :editor-fields="editorFields"
    />
    <div class="flex gap-3">
      <PVButton
        label="Discard"
        icon="pi pi-arrow-left"
        class="p-button-secondary p-button-outlined"
        @click="discard"
      />
      <div v-tooltip="saveTooltip">
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
      label="Editor Fields"
      :value="editorFields"
    />
    <StandardDebug
      label="Editor Values"
      :value="editorValues"
    />
  </StandardContent>
</template>
