<script setup lang="ts">
type Layout = 'mountains' | 'horizontal' | 'vertical'
type Background = 'blue' | 'white'

interface Props {
  layout?: Layout
  background?: Background
}
const props = withDefaults(defineProps<Props>(), {
  layout: 'horizontal',
  background: 'blue',
})
</script>

<template>
  <div
    :class="props.background === 'blue' ? 'bg-primary' : 'bg-white'"
    class="rmi-logo-wrapper"
  >
    <template v-if="props.background === 'blue'">
      <img
        v-if="props.layout === 'mountains'"
        src="@/assets/img/logo_mountains_bg_blue.svg"
      >
      <img
        v-else-if="props.layout === 'horizontal'"
        src="@/assets/img/logo_horizontal_bg_blue.svg"
      >
      <img
        v-else-if="props.layout === 'vertical'"
        src="@/assets/img/logo_vertical_bg_blue.svg"
      >
    </template>
    <template v-else>
      <img
        v-if="props.layout === 'mountains'"
        src="@/assets/img/logo_mountains_bg_white.svg"
      >
      <img
        v-else-if="props.layout === 'horizontal'"
        src="@/assets/img/logo_horizontal_bg_white.svg"
      >
      <img
        v-else-if="props.layout === 'vertical'"
        src="@/assets/img/logo_vertical_bg_white.svg"
      >
    </template>
  </div>
</template>

<style scoped lang="scss">
.rmi-logo-wrapper {
  img {
    height: 100%;
    width: 100%;
  }
}
</style>
