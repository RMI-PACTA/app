<script setup lang="ts">
/*
  LinkButton exists to allow us to style links like PrimeVue buttons, while
  still retaining the SPA routing functionality provided by NuxtLink. It's
  very tightly based on node_modules/primevue/button/Button.vue, updated for
  TypeScript + Nuxt + the Composition API.

  More resources:
  - https://router.vuejs.org/guide/advanced/extending-router-link.html
  - https://v3.nuxtjs.org/examples/routing/nuxt-link/
  - https://v3.nuxtjs.org/api/components/nuxt-link
  - https://primefaces.org/primevue/button
  - https://github.com/primefaces/primevue/issues/3118
  - https://vuejs.org/guide/components/attrs.html#disabling-attribute-inheritance
*/
import type { RouteLocationRaw } from 'vue-router'
import { useAttrs } from 'vue'
import Ripple from 'primevue/ripple'

// See https://vuejs.org/guide/reusability/custom-directives.html#introduction
const vRipple = Ripple

interface Props {
  to?: RouteLocationRaw
  newTab?: boolean

  label?: string
  icon?: string
  iconPos?: string
  iconClass?: string
  badge?: string
  badgeClass?: string
  loading?: string
  loadingIcon?: string
  activeClass?: string
  inactiveClass?: string
}
const props = withDefaults(defineProps<Props>(), {
  to: undefined,
  newTab: false,

  label: undefined,
  icon: undefined,
  iconPos: 'left',
  iconClass: undefined,
  badge: undefined,
  badgeClass: undefined,
  loading: undefined,
  loadingIcon: 'pi pi-spinner pi-spin',
  activeClass: '',
  inactiveClass: '',
})

const attrs = useAttrs()
const router = useRouter()

const iconStyleClass = computed(() => {
  return [
    props.loading !== undefined ? 'p-button-loading-icon ' + props.loadingIcon : props.icon,
    'p-button-icon',
    props.iconClass,
    {
      'p-button-icon-left': props.iconPos === 'left' && props.label,
      'p-button-icon-right': props.iconPos === 'right' && props.label,
      'p-button-icon-top': props.iconPos === 'top' && props.label,
      'p-button-icon-bottom': props.iconPos === 'bottom' && props.label,
    },
  ]
})

const badgeStyleClass = computed(() => {
  return [
    'p-badge p-component',
    props.badgeClass,
    {
      'p-badge-no-gutter': props.badge !== undefined && props.badge.length === 1,
    },
  ]
})

const isActive = computed(() => {
  const cr = router.currentRoute.value
  if (cr.fullPath === to.value) {
    return true
  }
  return false
})
const disabled = computed(() => {
  if (props.loading !== undefined) {
    return true
  }
  if (!('disabled' in attrs)) {
    return false
  }
  if (typeof (attrs.disabled) === 'boolean') {
    return attrs.disabled
  }
  if (isActive.value) {
    return true
  }
  return true
})

const buttonClass = computed(() => {
  const result: Record<string, boolean> = {
    'p-button': true,
    'p-component': true,
    'p-button-icon-only': props.icon !== undefined && props.label === undefined,
    'p-button-vertical': (props.iconPos === 'top' || props.iconPos === 'bottom') && !!props.label,
    'p-disabled': disabled.value,
    'p-button-loading': !!props.loading,
    'p-button-loading-label-only': props.loading !== undefined && props.icon === undefined && props.label !== undefined,
    'no-underline': true,
    'click-does-nothing': disabled.value,
    [props.activeClass]: isActive.value,
    [props.inactiveClass]: !isActive.value,
  }
  return result
})

const defaultAriaLabel = computed(() => {
  if (props.label !== undefined) {
    return props.label + (props.badge !== undefined ? ` ${props.badge}` : '')
  }
  if ('aria-label' in attrs && typeof (attrs['aria-label']) === 'string') {
    return attrs['aria-label']
  }
  return ''
})

const target = computed(() => props.newTab ? '_blank' : '_self')
const to = computed(() => {
  // This should only happen if disabled is set.
  return props.to ?? '/'
})
const href = computed(() => {
  if (props.to === undefined) {
    return '#'
  }

  if (typeof (props.to) === 'string') {
    return props.to
  }

  if ('path' in props.to) {
    return props.to.path
  }

  return '#'
})
</script>

<template>
  <NuxtLink
    v-slot="{ navigate }"
    :target="target"
    :to="to"
    :aria-disabled="disabled"
    custom
  >
    <a
      v-ripple
      v-bind="attrs"
      :href="href"
      :target="target"
      rel="noopener noreferrer"
      :class="buttonClass"
      type="button"
      :aria-label="defaultAriaLabel"
      :aria-disabled="disabled"
      @click="navigate"
    >
      <span
        v-if="props.loading && !props.icon"
        :class="iconStyleClass"
      />
      <span
        v-if="props.icon"
        :class="iconStyleClass"
      />
      <span
        v-if="props.label"
        class="p-button-label"
      >{{ props.label }}</span>
      <span
        v-if="props.badge"
        :class="badgeStyleClass"
      >{{ props.badge }}</span>
      <slot />
    </a>
  </NuxtLink>
</template>

<style lang="scss">
.click-does-nothing {
  pointer-events: none;
  opacity: 0.7;
}
</style>
