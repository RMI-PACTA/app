<script setup lang="ts">
import { type LanguageCode } from '@/lib/language'
import type OverlayPanel from 'primevue/overlaypanel'
import { useToast } from 'primevue/usetoast'

const toast = useToast()
const switchLocalePath = useSwitchLocalePath()
const { t, locale: baseLocale, locales: baseLocales } = useI18n()
const { languageWasSelectedOrDismissed } = useLocalStorage()

// We go through this rigamarole because the additional information passed into
// the i18n configuration is present here (namely, well, name), but it doesn't
// appear in the typescript configuration.
interface Option {
  code: LanguageCode
  name: string
}
const locales = computed<Option[]>(() => baseLocales.value as Option[])
const locale = computed<Option>(
  () => presentOrFileBug(locales.value.find(o => o.code === baseLocale.value)))

const prefix = 'components/locale/Selector'
const tt = (s: string) => t(`${prefix}.${s}`)
const visible = useState<boolean>(`${prefix}.visible`, () => false)

const overlayPanel = useState<OverlayPanel>(`${prefix}.overlayPanel`)
const toggleMenu = (event: Event) => {
  presentOrFileBug(overlayPanel.value).toggle(event)
  visible.value = !visible.value
}
const hideMenu = () => {
  visible.value = false
}
const showMenu = () => {
  visible.value = true
}

onMounted(() => {
  if (!languageWasSelectedOrDismissed.value) {
    toast.add({
      severity: 'success',
      group: 'language-selector',
    })
  }
})
const closeToast = () => {
  toast.removeGroup('language-selector')
  languageWasSelectedOrDismissed.value = true
}
</script>

<template>
  <div>
    <PVButton
      class="p-button-rounded px-2 py-1"
      :class="visible ? '' : 'p-button-text'"
      @click="toggleMenu"
    >
      <LanguageRepresentation
        :code="locale.code"
        :full-name="false"
      />
    </PVButton>
    <PVOverlayPanel
      ref="overlayPanel"
      @hide="hideMenu"
      @show="showMenu"
    >
      <div class="flex flex-column gap-1 align-items-start">
        <LinkButton
          v-for="option in locales"
          :key="option.code"
          class="flex gap-3 justify-content-center"
          :class="option.code === locale.code ? 'p-button-outlined' : 'p-button-text'"
          :to="switchLocalePath(option.code)"
          @click="toggleMenu"
        >
          <LanguageRepresentation :code="option.code" />
        </LinkButton>
      </div>
    </PVOverlayPanel>
    <PVToast
      position="bottom-right"
      group="language-selector"
      :dismissable="false"
      @close="closeToast"
    >
      <template #message>
        <div class="flex flex-column gap-2">
          <div>{{ tt('Availabile In') }} </div>
          <div class="flex gap-1">
            <LinkButton
              v-for="option in locales"
              :key="option.code"
              class="flex gap-1 p-2 justify-content-center p-button-rounded"
              :class="option.code === locale.code ? 'p-button-outlined' : 'p-button-text'"
              :to="switchLocalePath(option.code)"
              @click="closeToast"
            >
              <LanguageRepresentation
                :code="option.code"
                :full-name="false"
              />
            </LinkButton>
          </div>
          <div>{{ tt('Bottom Right') }}</div>
        </div>
      </template>
    </PVToast>
  </div>
</template>
