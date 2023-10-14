import { type RouteLocationNormalized } from 'vue-router'

const getLocaleFromRoute = (route: RouteLocationNormalized): string => {
  if (route.fullPath.startsWith('/en/') || route.fullPath === '/en') {
    return 'en'
  }
  if (route.fullPath.startsWith('/fr/') || route.fullPath === '/fr') {
    return 'fr'
  }
  if (route.fullPath.startsWith('/es/') || route.fullPath === '/es') {
    return 'es'
  }
  if (route.fullPath.startsWith('/de/') || route.fullPath === '/de') {
    return 'de'
  }
  return ''
}

export default defineNuxtRouteMiddleware((to: RouteLocationNormalized, from: RouteLocationNormalized) => {
  const toLocale = getLocaleFromRoute(to)
  if (toLocale === '') {
    console.error(`Navigating to ${to.fullPath} is missing locale!`)
    return
  }
  const fromLocale = getLocaleFromRoute(from)
  if (toLocale === fromLocale) {
    return
  }
  const toWithoutLocale = to.fullPath.replace(`/${toLocale}`, '')
  const fromWithoutLocale = from.fullPath.replace(`/${fromLocale}`, '')
  if (toWithoutLocale === fromWithoutLocale) {
    return
  }
  console.error(`Navigating to ${to.fullPath} is missing navigation guard!`)
})
