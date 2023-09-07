<script setup lang="ts">
import { dialogBreakpoints } from '@/lib/breakpoints'

const { error: { errorModalVisible, error } } = useModal()

const fullError = computed(() => {
  return error.value
    ? {
        name: error.value.name,
        message: error.value.message,
        stack: error.value.stack?.split('\n')
      }
    : ''
})
</script>

<template>
  <div>
    <PVDialog
      v-model:visible="errorModalVisible"
      :closable="true"
      :modal="true"
      :close-on-escape="true"
      :dismissable-mask="true"
      :draggable="false"
      :breakpoints="dialogBreakpoints"
    >
      <template #header>
        <div class="flex justify-content-between align-items-center w-full">
          <div>
            <div class="font-bold text-xl mb-2">
              An error ocurred
            </div>
            <div class="text-sm">
              Sorry about that, our team take bug reports seriously, and will try to make it right!
            </div>
          </div>
        </div>
      </template>
      <div class="flex gap-3 flex-column">
        <code class="w-full">
          {{ error }}
        </code>
        <StandardDebug
          label="Error Trace"
          :value="fullError"
          always
        />
      </div>
      <template #footer>
        <div class="text-left text-sm">
          Some common troubleshooting steps that might be helpful:
          <ul>
            <li><b>Refresh this page</b> - most of our pages save your progress as you go, so it's almost always fine to reload the page.</li>
            <li><b>Check your internet connection</b> - this site requires connection to the internet for most functionality.</li>
            <li><b>Visit this site on a desktop computer</b> - this site works best on desktop web browsers.</li>
          </ul>
          If this issue persists, please report this issue by <a
            href="https://github.com/RMI/opgee-api/issues/new"
            target="_blank"
          >filing a bug in the OPGEE repository</a>.
        </div>
      </template>
    </PVDialog>
  </div>
</template>
