<script setup lang="ts">
interface Props {
  routeBackOnClose?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  routeBackOnClose: false
})

const { error: { errorModalVisible, error: modalError } } = useModal()
const error = useError()
const router = useRouter()

const maybeGoBack = async () => {
  if (!props.routeBackOnClose) {
    return
  }

  if (window.history.length > 1) {
    await clearError().then(router.back)
  } else {
    await clearError({ redirect: '/' })
  }
}

const fullError = computed(() => {
  const err = error.value ?? modalError.value
  if (err instanceof Error) {
    return {
      name: err.name ?? '',
      message: err.message,
      stack: err.stack?.split('\n')
    }
  } else if (err) {
    return err
  } else {
    return ''
  }
})
</script>

<template>
  <StandardModal
    v-model:visible="errorModalVisible"
    header="An error ocurred"
    sub-header="Sorry about that, our team take bug reports seriously, and will try to make it right!"
    @closed="maybeGoBack"
  >
    <StandardDebug
      label="Error Details"
      :value="fullError"
      always
    />
    <div class="text-left text-sm">
      Some common troubleshooting steps that might be helpful:
      <ul>
        <li><b>Refresh this page</b> - most of our pages save your progress as you go, so it's almost always fine to reload the page.</li>
        <li><b>Check your internet connection</b> - this site requires connection to the internet for most functionality.</li>
        <li><b>Visit this site on a desktop computer</b> - this site works best on desktop web browsers.</li>
      </ul>
      If this issue persists, please report this issue by <a
        href="https://github.com/RMI-pacta/app/issues/new"
        target="_blank"
      >filing a bug in the PACTA repository</a>.
    </div>
  </StandardModal>
</template>
