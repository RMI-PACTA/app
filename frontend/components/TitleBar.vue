<script setup lang="ts">
import { onMounted } from 'vue'

const router = useRouter()

interface Props {
  title: string
}
const props = withDefaults(defineProps<Props>(), {})

const back = () => { router.go(-1) }

onMounted(() => {
  document.title = 'RMI | PACTA | ' + props.title
})
</script>

<template>
  <div class="relative flex flex-column gap-2 align-items-start">
    <PVButton
      icon="pi pi-angle-left"
      class="p-button-secondary p-button-text md:hidden p-0 text-600 pr-2 -ml-1"
      label="Back"
      @click="back"
    />
    <PVButton
      v-tooltip.right="'Go Back'"
      icon="pi pi-angle-left"
      class="p-button-secondary p-button-text hidden md:block md:absolute p-0 h-auto w-auto"
      style="left: -1.75rem; top: .5rem;"
      rounded
      @click="back"
    />
    <h1 class="inline-block text-primary title-bar mb-0">
      {{ title }}
    </h1>
  </div>
</template>

<style scoped lang="scss">
.title-bar {
  margin: 0;
}
</style>
