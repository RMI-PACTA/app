<script setup lang="ts">
import { type PactaVersion } from 'openapi/generated/pacta'

const pactaVersion = useState<PactaVersion>('admin.pactaVersion', () => ({
  name: '',
  description: '',
  digest: '',
  createdAt: '',
  isDefault: false
}))
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
  </StandardContent>
</template>
