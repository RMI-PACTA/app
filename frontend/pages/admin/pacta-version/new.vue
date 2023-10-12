<script setup lang="ts">
import { pactaVersionEditor } from '@/lib/editor'

const prefix = 'admin/pacta-version/new'
const router = useRouter()
const pactaClient = await usePACTA()
const { loading: { withLoading } } = useModal()

const {
  incompleteFields,
  hasChanges,
  isIncomplete,
  editorPactaVersion,
  pactaVersion,
} = pactaVersionEditor({
  id: '',
  name: '',
  description: '',
  digest: '',
  createdAt: '',
  isDefault: false,
})

const saveTooltip = computed<string | undefined>(() => {
  if (!hasChanges.value) { return 'All changes saved' }
  if (isIncomplete.value) { return `Cannot save with incomplete fields: ${incompleteFields.value.join(', ')}` }
  return undefined
})
const saveDisabled = computed<boolean>(() => saveTooltip.value !== undefined)

const discard = () => router.push('/admin/pacta-version')
const save = () => withLoading(
  () => pactaClient.createPactaVersion(pactaVersion.value).then(() => router.push('/admin/pacta-version')),
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
      v-model:editorPactaVersion="editorPactaVersion"
    />
    <div class="flex gap-3">
      <PVButton
        label="Discard"
        icon="pi pi-arrow-left"
        class="p-button-secondary p-button-outlined"
        @click="discard"
      />
      <PVButton
        :disabled="saveDisabled"
        label="Save"
        icon="pi pi-arrow-right"
        icon-pos="right"
        @click="save"
      />
    </div>
    <StandardDebug
      label="PACTA Version"
      :value="editorPactaVersion"
    />
  </StandardContent>
</template>
