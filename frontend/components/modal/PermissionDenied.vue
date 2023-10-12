<script setup lang="ts">
const { permissionDenied: { permissionDeniedVisible, permissionDeniedError } } = useModal()
const router = useRouter()
const { t } = useI18n()
const localePath = useLocalePath()

const prefix = 'ModalPermissionDenied'
const tt = (s: string) => t(`${prefix}.${s}`)

const navigateBack = () => {
  // If the user didn't land here directly, go back to the previous page.
  if (window.history.length > 1) {
    router.back()
    return
  }
  // Otherwise, just take them home.
  return router.push(localePath('/'))
}
</script>

<template>
  <StandardModal
    v-model:visible="permissionDeniedVisible"
    :header="tt('Heading')"
    :sub-header="tt('Subheading')"
    @closed="navigateBack"
  >
    <p>
      {{ tt('How You Got Here') }}
    </p>
    <p>
      {{ tt('Email Us') }}
      <a href="mailto:pacta-help@siliconally.org">pacta-help@siliconally.org</a>.
    </p>
    <StandardDebug
      :label="tt('Technical Error')"
      :value="permissionDeniedError"
      always
    />
  </StandardModal>
</template>
