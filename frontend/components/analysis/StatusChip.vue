<script setup lang="ts">
import { type Analysis } from '@/openapi/generated/pacta'

const { t } = useI18n()
const prefix = 'components/analysis/StatusChip'
const tt = (s: string) => t(`${prefix}.${s}`)

interface Props {
  analysis: Analysis
}
const props = defineProps<Props>()

const isStale = (a: Analysis): boolean => new Date().getTime() - new Date(a.createdAt).getTime() > 1000 * 60 * 10
const status = computed(() => {
  if (props.analysis.completedAt !== undefined) {
    if (props.analysis.failureMessage !== undefined) {
      return tt('Failed')
    } else {
      return tt('Completed')
    }
  }
  if (isStale(props.analysis)) {
    return tt('PotentialTimeout')
  }
  return tt('Running')
})
const severity = computed(() => {
  if (props.analysis.completedAt !== undefined) {
    if (props.analysis.failureMessage !== undefined) {
      return 'danger'
    } else {
      return 'success'
    }
  }
  if (isStale(props.analysis)) {
    return 'danger'
  }
  return 'warn'
})
</script>

<template>
  <PVInlineMessage
    :severity="severity"
  >
    {{ status }}
  </PVInlineMessage>
</template>
