<script setup lang="ts">
import { ZIndexUtils } from 'primevue/utils'
import { watch } from 'vue'

const { loading: { loading, loadingSet } } = useModal()

const prefix = 'ModalLoading'
const loadingModalEl = useState<HTMLElement>(`${prefix}.loadingModalEl`)

const refreshZIndex = () => {
  const mg = presentOrReload(loadingModalEl.value.parentElement)
  if (loading.value) {
    ZIndexUtils.set('ModalGroup', mg, 10000)
    ZIndexUtils.set('LoadingModal', loadingModalEl.value, 10001)
  } else {
    ZIndexUtils.clear(mg)
    ZIndexUtils.clear(loadingModalEl.value)
  }
}

const debug = useState<boolean>(`${prefix}.debug`, () => false)
const toggleDebug = () => {
  debug.value = !debug.value
}

watch(() => loading.value, refreshZIndex)
</script>

<template>
  <div
    v-show="loading"
    ref="loadingModalEl"
    data-anchor-id="loadingModal"
    class="loading-modal"
    @click="toggleDebug"
  >
    <img
      src="@/assets/img/logo_loading_animation_v1.gif"
      class="gif shadow-3"
    >
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

  .gif {
    width: 50vw;
    max-width: 10rem;
  }
}
</style>
