<script setup lang="ts">
import { LanguageOptions, type LanguageCode } from '@/lib/language'

interface Props {
  code: LanguageCode | undefined
  fullName?: boolean
}
const props = withDefaults(defineProps<Props>(), { fullName: true })

const language = computed(() => {
  let code = props.code
  if (!code) {
    code = 'en'
  }
  return presentOrFileBug(LanguageOptions.find((l) => l.code === code), `Language ${code} not found`)
})
</script>

<template>
  <div
    class="flex align-items-center"
    :class="props.fullName ? 'gap-3' : 'gap-2'"
  >
    <div class="flag-wrapper shadow-1">
      <img
        :src="`/img/flags/${language.code}.svg`"
        :class="language.code"
      >
    </div>
    <span class="text-lg">{{ props.fullName ? language.label : language.code.toUpperCase() }}</span>
  </div>
</template>

<style lang="scss" scoped>
.flag-wrapper {
    width: 1rem;
    height: 1rem;
    overflow: hidden;
    display: flex;
    justify-content: center;
    align-items: center;
    border-radius: 50%;

    img {
      flex-shrink: 0;
      min-width: 100%;
      min-height: 100%
    }

    img.es {
      flex-shrink: initial;
    }
}
</style>
