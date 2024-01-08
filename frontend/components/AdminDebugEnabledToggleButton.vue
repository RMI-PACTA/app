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

const prefix = 'components/AdminDebugEnabledToggleButton'
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
}
const noAck = () => {
  everAcked.value = false
  model.value = false
}

const visible = newModalVisibilityState('AdminDebugEnabledWarning')
</script>

<template>
  <div>
    <ExplicitInputSwitch
      v-model:value="model"
      :on-label="tt('Administrator Debugging Access Enabled')"
      :off-label="tt('No Administrator Access Enabled')"
    />
    <StandardModal
      v-model:visible="visible"
      :header="tt('ModalHeading')"
      :sub-header="tt('ModalSubheading')"
    >
      <p>
        TODO - add good copy for this
      </p>
      <p>
        You're enabling administrator access to this resource. If you do so, site administrators will be able to
        access the content of this data.
      </p>
      <div class="flex gap-2 justify-content-between align-items-center flex-wrap">
        <PVButton
          :label="tt('No Ack')"
          icon="pi pi-arrow-left"
          class="p-button-secondary"
          @click="noAck"
        />
        <PVButton
          :label="tt('Ack')"
          icon="pi pi-arrow-right"
          class="p-button-secondary"
          @click="ack"
        />
      </div>
    </StandardModal>
  </div>
</template>
