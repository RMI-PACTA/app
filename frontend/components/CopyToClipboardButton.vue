<script setup lang="ts">
import { computed } from 'vue'

const { t } = useI18n()

interface Props {
  value: string
  cta: string
}
const props = defineProps<Props>()

const prefix = 'CopyToClipboardButton'
const tt = (key: string) => t(`${prefix}.${key}`)

const statePrefix = `${prefix}[${useStateIDGenerator().id()}]`
const copiedToClipboard = useState<boolean>(`${statePrefix}.copiedToClipboard`, () => false)
const message = computed(() => copiedToClipboard.value ? tt('Copied') : props.cta)
const icon = computed(() => copiedToClipboard.value ? 'pi pi-check' : 'pi pi-copy')

const copyToClipboard = async () => {
  await navigator.clipboard.writeText(props.value)
  copiedToClipboard.value = true
  setTimeout(() => { copiedToClipboard.value = false }, 5000)
}
</script>

<template>
  <PVButton
    :disabled="copiedToClipboard"
    :label="message"
    :icon="icon"
    icon-pos="right"
    class="text-sm"
    @click="copyToClipboard"
  />
</template>
