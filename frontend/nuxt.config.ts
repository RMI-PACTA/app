// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: [
    '@nuxtjs/i18n',
  ],
  build: {
    transpile: [
      'primevue', // https://primevue.org/installation/#nuxtintegration
      'vue-i18n',
    ],
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

    timeline: {
      enabled: true,
    },
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
      baseURL: process.env.BASE_URL ?? '',
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
    baseUrl: process.env.BASE_URL,
    strategy: (process.env.I18N_STRATEGY ?? 'prefix_except_default') as 'no_prefix' | 'prefix_except_default' | 'prefix' | 'prefix_and_default',
    vueI18n: './i18n.config.ts',
    locales: [
      { code: 'en', language: 'en-US', file: { path: 'en.json', cache: false }, name: 'English' },
      { code: 'fr', language: 'fr-FR', file: { path: 'fr.json', cache: false }, name: 'Français' },
      { code: 'es', language: 'es-ES', file: { path: 'es.json', cache: false }, name: 'Español' },
      { code: 'de', language: 'de-DE', file: { path: 'de.json', cache: false }, name: 'Deutsch' },
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
