<script setup lang="ts">
import { type PactaVersion } from '@/openapi/generated/pacta'

const router = useRouter()
const { pactaClient } = useAPI()
const { loading: { withLoading } } = useModal()

const prefix = 'admin/pacta-version/new'
const pactaVersion = useState<PactaVersion>(`${prefix}.pactaVersion`, () => ({
  id: '',
  name: '',
  description: '',
  digest: '',
  createdAt: '',
  isDefault: false
}))

const discard = () => router.push('/admin/pacta-version')
const save = () => withLoading(
  () => pactaClient.createPactaVersion(pactaVersion.value).then(() => router.push('/admin/pacta-version')),
  `${prefix}.save`
)
</script>

<template>
  <StandardContent>
    <TitleBar title="New PACTA Version" />
    <p>
      Pacta version info goes here
    </p>
    <PactaversionEditor
      v-model:pactaVersion="pactaVersion"
    />
    <div class="flex gap-3">
      <PVButton
        label="Discard"
        icon="pi pi-arrow-left"
        class="p-button-secondary p-button-outlined"
        @click="discard"
      />
      <PVButton
        label="Save"
        icon="pi pi-arrow-right"
        icon-pos="right"
        @click="save"
      />
    </div>
    <StandardDebug
      label="PACTA Version"
      :value="pactaVersion"
    />
  </StandardContent>
</template>
