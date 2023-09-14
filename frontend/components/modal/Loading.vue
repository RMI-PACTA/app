<script setup lang="ts">
const { loading: { loading, loadingSet } } = useModal()

const prefix = 'ModalLoading'
const loadingModalEl = useState<HTMLElement>(`${prefix}.loadingModalEl`)

const debug = useState<boolean>(`${prefix}.debug`, () => false)
const toggleDebug = () => {
  debug.value = !debug.value
}

/*
See https://github.com/RMI-PACTA/app/pull/7#r1319209535 - we may be able to remove this if we don't use modals with nested auto-z-indexed elements

import { ZIndexUtils } from 'primevue/utils'
import { watch } from 'vue'

const refreshZIndex = () => {
  const mg = presentOrSuggestReload(loadingModalEl.value.parentElement)
  if (loading.value) {
    ZIndexUtils.set('ModalGroup', mg, 10000)
    ZIndexUtils.set('LoadingModal', loadingModalEl.value, 10001)
  } else {
    ZIndexUtils.clear(mg)
    ZIndexUtils.clear(loadingModalEl.value)
  }
}
watch(() => loading.value, refreshZIndex)
*/
</script>

<template>
  <div
    v-show="loading"
    ref="loadingModalEl"
    data-anchor-id="loadingModal"
    class="loading-modal"
  >
    <div @click="() => { toggleDebug() } ">
      <video
        class="loading-animation shadow-3"
        autoplay
        loop
        muted
        playsinline
      >
        <source
          src="@/assets/img/logo_loading_animation_v1.webm"
          type="video/webm"
        >
        <source
          src="@/assets/img/logo_loading_animation_v1.mp4"
          type="video/mp4"
        >
      </video>
    </div>
    <ul
      v-if="debug"
      class="demo-mode border-1 pr-3 py-2"
    >
      <li
        v-for="l in loadingSet"
        :key="l"
      >
        {{ l }}
      </li>
    </ul>
  </div>
</template>

<style lang="scss">
.loading-modal {
  box-shadow: none;
  overflow: visible;
  position: fixed;
  top: 0;
  left: 0;
  height: 100vh;
  width: 100vw;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background: rgb(255 255 255 / 60%);

  .loading-animation {
    width: 50vw;
    max-width: 10rem;
  }
}
</style>
