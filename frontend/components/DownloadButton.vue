<script setup lang="ts">
import { computed } from 'vue'

const { t } = useI18n()

interface Props {
  value: string
  cta?: string | undefined
  fileName: string
}
const props = defineProps<Props>()

const prefix = 'components/DownloadButton'
const tt = (key: string) => t(`${prefix}.${key}`)

const statePrefix = `${prefix}[${useStateIDGenerator().id()}]`
const downloaded = useState<boolean>(`${statePrefix}.downloaded`, () => false)
const message = computed(() => downloaded.value ? tt('Downloaded') : props.cta)
const icon = computed(() => downloaded.value ? 'pi pi-check' : 'pi pi-download')

const download = () => {
  downloaded.value = true
  const a = document.createElement('a')
  const file = new Blob([props.value])
  a.href = URL.createObjectURL(file)
  a.download = props.fileName
  a.click()
  a.remove()
  setTimeout(() => { downloaded.value = false }, 5000)
}
</script>

<template>
  <PVButton
    :disabled="downloaded"
    :label="message"
    :icon="icon"
    icon-pos="right"
    class="text-sm"
    @click="download"
  />
</template>
