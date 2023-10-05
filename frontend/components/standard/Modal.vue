<script setup lang="ts">
import { dialogBreakpoints, smallDialogBreakpoints } from '@/lib/breakpoints'

const props = withDefaults(defineProps<{
  small?: boolean
  visible: boolean
  header: string
  subHeader?: string
}>(), {
  small: false,
  subHeader: '',
})

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'closed'): void
}>()

const visible = computed<boolean>({
  get: () => props.visible,
  set: (value) => { emit('update:visible', value) },
})
const breakpoints = computed(() => props.small ? smallDialogBreakpoints : dialogBreakpoints)

const onHide = () => { emit('closed') }
</script>

<template>
  <PVDialog
    v-model:visible="visible"
    :closable="true"
    :modal="true"
    :close-on-escape="true"
    :dismissable-mask="true"
    :draggable="false"
    :breakpoints="breakpoints"
    @hide="onHide"
  >
    <template #header>
      <div class="flex justify-content-between align-items-center w-full">
        <div>
          <div class="font-bold text-xl mb-2">
            {{ props.header }}
          </div>
          <div
            v-if="props.subHeader"
            class="text-md"
          >
            {{ props.subHeader }}}
          </div>
        </div>
      </div>
    </template>

    <div class="flex flex-column gap-3 reset-child-margin">
      <slot />
    </div>
  </PVDialog>
</template>

<style lang="scss">
.reset-child-margin {
  & > * {
    margin: 0;
  }
}
</style>
