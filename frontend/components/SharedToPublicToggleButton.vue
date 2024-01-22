<script setup lang="ts">
const { newModalVisibilityState } = useModal()
const { computedBooleanLocalStorageValue } = useLocalStorage()
const { t } = useI18n()

interface Props {
  value: boolean
}
const props = defineProps<Props>()
interface Emits {
  (e: 'update:value', value: boolean): void
}
const emit = defineEmits<Emits>()

const prefix = 'components/SharedToPublicToggleButton'
const tt = (s: string) => t(`${prefix}.${s}`)
const everAcked = computedBooleanLocalStorageValue(`${prefix}.everAcked`, false)

const model = computed({
  get: () => props.value,
  set: (value: boolean) => {
    if (value && !everAcked.value) {
      visible.value = true
      return
    } else {
      emit('update:value', value)
    }
  },
})

const ack = () => {
  everAcked.value = true
  model.value = true
  visible.value = false
}
const noAck = () => {
  model.value = false
}

const visible = newModalVisibilityState('SharedToPublicWarning')
</script>

<template>
  <div>
    <ExplicitInputSwitch
      v-model:value="model"
      :on-label="tt('Shared to Public')"
      :off-label="tt('Not Shared')"
    />
    <StandardModal
      v-model:visible="visible"
      :header="tt('ModalHeading')"
      :sub-header="tt('ModalSubheading')"
    >
      <p>
        {{ tt('Paragraph1') }}
      </p>
      <p>
        {{ tt('Paragraph2') }}
      </p>
      <div class="flex pt-3 gap-2 justify-content-between align-items-center flex-wrap">
        <PVButton
          :label="tt('No Ack')"
          icon="pi pi-arrow-left"
          class="p-button-secondary"
          @click="noAck"
        />
        <PVButton
          :label="tt('Ack')"
          icon="pi pi-arrow-right"
          icon-pos="right"
          @click="ack"
        />
      </div>
    </StandardModal>
  </div>
</template>
