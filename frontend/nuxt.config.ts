// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: [
    '@nuxtjs/i18n',
  ],
  build: {
    // https://primevue.org/installation/#nuxtintegration
    transpile: ['primevue', 'vue-i18n'],
  },
  css: [
    '@/assets/css/overrides.css',
    '@/assets/css/theme.css',
    'primevue/resources/primevue.css',
    'primeicons/primeicons.css',
    'primeflex/primeflex.css',
  ],
  devtools: {
    enabled: true,
  },
  runtimeConfig: {
    public: {
      apiServerURL: process.env.API_SERVER_URL ?? '',
      authServerURL: process.env.AUTH_SERVER_URL ?? '',
      msalConfig: {
        userFlowName: process.env.MSAL_USER_FLOW_NAME ?? '',
        userFlowAuthority: process.env.MSAL_USER_FLOW_AUTHORITY ?? '',
        authorityDomain: process.env.MSAL_AUTHORITY_DOMAIN ?? '',
        clientID: process.env.MSAL_CLIENT_ID ?? '',
        redirectURI: process.env.MSAL_REDIRECT_URI ?? '',
        logoutURI: process.env.MSAL_LOGOUT_URI ?? '',
        minLogLevel: process.env.MSAL_MIN_LOG_LEVEL ?? 'VERBOSE',
      },
    },
  },
  typescript: {
    strict: true,
  },
  imports: {
    presets: [
      {
        from: 'vue',
        imports: ['computed', 'onMounted'],
      },
    ],
    dirs: ['globalimports'],
  },
  i18n: {
    // TODO(grady) Set a Base URL once we have it for SEO https://i18n.nuxtjs.org/guide/seo#:~:text=You%20must%20also%20set%20the%20baseUrl%20option%20to%20your%20production%20domain%20in%20order%20to%20make%20alternate%20URLs%20fully%2Dqualified%3A
    // baseUrl: 'https://my-nuxt-app.com'
    strategy: 'prefix_and_default',
    vueI18n: './i18n.config.ts',
    locales: [
      { code: 'en', iso: 'en-US', file: { path: 'en.json', cache: false }, flag: '🇬🇧', name: 'English' },
      { code: 'fr', iso: 'fr-FR', file: { path: 'fr.json', cache: false }, flag: '🇫🇷', name: 'Français' },
      { code: 'es', iso: 'es-ES', file: { path: 'es.json', cache: false }, flag: '🇩🇪', name: 'Deutsch' },
      { code: 'de', iso: 'de-DE', file: { path: 'de.json', cache: false }, flag: '🇪🇸', name: 'Español' },
    ],
    lazy: true,
    langDir: 'lang',
    defaultLocale: 'en',
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: 'i18n_redirected',
      redirectOn: 'root',
    },
  },
})
