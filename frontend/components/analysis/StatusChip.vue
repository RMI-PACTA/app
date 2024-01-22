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
const bg = computed(() => {
  if (props.analysis.completedAt !== undefined) {
    if (props.analysis.failureMessage !== undefined) {
      return 'bg-red-400'
    } else {
      return 'bg-green-200'
    }
  }
  if (isStale(props.analysis)) {
    return 'bg-red-200'
  }
  return 'bg-yellow-200'
})
</script>

<template>
  <div
    :class="bg"
    class="p-1 align-self-stretch flex align-items-center justify-content-center border-round-lg"
  >
    <span>
      {{ status }}
    </span>
  </div>
</template>
