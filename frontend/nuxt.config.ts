// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  build: {
    // https://primevue.org/installation/#nuxtintegration
    transpile: ['primevue'],
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
})
