<script setup lang="ts">
import { deserializeError } from 'serialize-error'

const { error: { errorModalVisible, error: modalError } } = useModal()
const error = useError()
const router = useRouter()
const { t } = useI18n()

const prefix = 'components/modal/Error'
const tt = (s: string) => t(`${prefix}.${s}`)

interface Props {
  routeBackOnClose?: boolean
}
const props = withDefaults(defineProps<Props>(), {
  routeBackOnClose: false,
})

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
  const err = error.value ?? deserializeError(modalError.value)
  if (err instanceof Error) {
    return {
      name: err.name ?? '',
      message: err.message,
      stack: err.stack?.split('\n'),
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
    :header="tt('Heading')"
    :sub-header="tt('Subheading')"
    @closed="maybeGoBack"
  >
    <StandardDebug
      :label="tt('Error Details')"
      :value="fullError"
      always
    />
    <div class="text-left text-sm">
      Some common troubleshooting steps that might be helpful:
      <ul>
        <li><b>{{ tt('Refresh' ) }}</b> - {{ tt('Refresh Explanation') }}</li>
        <li><b>{{ tt('Connection') }}</b> - {{ tt('Connection Explanation') }}</li>
        <li><b>{{ tt('Desktop') }}</b> - {{ tt('Desktop Explanation') }}</li>
      </ul>
      {{ tt('If Persists') }} <a
        href="https://github.com/RMI-pacta/app/issues/new"
        target="_blank"
      >{{ tt('File a Bug' ) }}</a>.
    </div>
  </StandardModal>
</template>
