// Based heavily on [1], adapted for Nuxt 3. This was used instead of [2]
// because the Vue 3 sample doesn't account for SSR, and is generally more
// complicated than we need.
// [1] https://github.com/Azure-Samples/ms-identity-b2c-javascript-spa
// [2] https://github.com/AzureAD/microsoft-authentication-library-for-js/tree/dev/samples/msal-browser-samples/vue3-sample-app
import { type AccountInfo, type AuthenticationResult, InteractionRequiredAuthError, type SilentRequest } from '@azure/msal-browser'
import type { ComputedRef } from 'vue'
import { type APIKey } from '~/openapi/generated/user'

interface MSAL {
  signIn: () => Promise<void>
  signOut: () => Promise<void>
  createAPIKey: () => Promise<APIKey>
  isAuthenticated: ComputedRef<boolean>
}

export const useMSAL = async (): Promise<MSAL> => {
  const router = useRouter()
  const isAuthenticated = useState('useMSAL.isAuthenticated', () => false)
  if (process.server) {
    const sessionCookie = useCookie('jwt')
    isAuthenticated.value = sessionCookie.value !== null && sessionCookie.value !== undefined
    return {
      signIn: async () => { await Promise.reject(new Error('cannot call signIn on server')) },
      signOut: async () => { await Promise.reject(new Error('cannot call signOut on server')) },
      createAPIKey: async () => await Promise.reject(new Error('cannot call createAPIKey on server')),
      isAuthenticated: computed(() => isAuthenticated.value),
    }
  }

  const { $msal: { /* inProgress, */ instance, accounts, msalConfig, b2cPolicies } } = useNuxtApp()
  const { userClientWithCustomToken } = useAPI()

  const clientInitialized = useState('useMSAL.clientInitialized', () => false)

  const account = useState<AccountInfo | undefined>('useMSALAuthentication.account')

  const setAccount = (acctInfo: AccountInfo): void => {
    account.value = acctInfo
    // TODO: maybe welcome new user?
  }

  const selectAccount = async (): Promise<void> => {
    if (accounts.value.length < 1) {
      // TODO: Figure out how to handle this scenario
    } else if (accounts.value.length > 1) {
      /**
         * Due to the way MSAL caches account objects, the auth response from initiating a user-flow
         * is cached as a new account, which results in more than one account in the cache. Here we make
         * sure we are selecting the account with homeAccountId that contains the sign-up/sign-in user-flow,
         * as this is the default flow the user initially signed-in with.
         */
      const filteredAccounts = accounts.value.filter((account): boolean => {
        return account.idTokenClaims?.iss !== undefined &&
          account.idTokenClaims.aud !== undefined &&
          account.homeAccountId.toUpperCase().includes(b2cPolicies.names.signUpSignIn.toUpperCase()) &&
          account.idTokenClaims.iss.toUpperCase().includes(b2cPolicies.authorityDomain.toUpperCase()) &&
          account.idTokenClaims.aud === msalConfig.auth.clientId
      })

      if (filteredAccounts.length > 1) {
        // localAccountId identifies the entity for which the token asserts information.
        if (filteredAccounts.every((account) => account.localAccountId === filteredAccounts[0].localAccountId)) {
          // All filteredAccounts belong to the same user
          setAccount(filteredAccounts[0])
        } else {
          // Multiple users detected. Logout all to be safe.
          console.log('multiple accounts found, logging out', accounts.value)
          await signOut()
        }
      } else if (filteredAccounts.length === 1) {
        setAccount(filteredAccounts[0])
      }
    } else if (accounts.value.length === 1) {
      setAccount(accounts.value[0])
    }
  }

  const handleResponse = async (response: AuthenticationResult): Promise<void> => {
    setAccount(response.account)
    const userClient = userClientWithCustomToken(response.idToken)
    try {
      await userClient.login()
      isAuthenticated.value = true
    } catch {
      await signOut()
    }
  }

  const signIn = async (): Promise<void> => {
    const req = { scopes: ['openid'] }
    try {
      const response = await instance.loginPopup(req)
      await handleResponse(response)
    } catch (err) {
      console.log('useMSAL.loginPopup', err)
    }
  }

  const getToken = async (): Promise<AuthenticationResult> => {
    if (account.value === undefined) {
      // TODO: Figure out if this is a legitimate usecase.
      return await Promise.reject(new Error('tried to get a token, but no account was found'))
    }

    const request: SilentRequest = {
      scopes: [],
      forceRefresh: false, // Set this to "true" to skip a cached token and go to the server to get a new token
      account: instance.getAccountByHomeId(account.value.homeAccountId) ?? undefined,
    }

    const response = await instance.acquireTokenSilent(request)
    // In case the response from B2C server has an empty accessToken field
    // throw an error to initiate token acquisition
    if (response.idToken === '') {
      throw new InteractionRequiredAuthError()
    }
    return response
  }

  const getTokenPopup = async (): Promise<AuthenticationResult> => {
    if (account.value === undefined) {
      // TODO: Figure out if this is a legitimate usecase.
      throw new Error('tried to get a token, but no account was found')
    }

    const request: SilentRequest = {
      scopes: [],
      forceRefresh: false, // Set this to "true" to skip a cached token and go to the server to get a new token
      account: instance.getAccountByHomeId(account.value.homeAccountId) ?? undefined,
    }

    try {
      return await getToken()
    } catch (error) {
      console.log('Silent token acquisition failed. Acquiring token using popup. \n', error)
      if (!(error instanceof InteractionRequiredAuthError)) {
        throw new Error('unexpected error while getting token', { cause: error })
      }
      // Fallback to interaction when silent call fails
      try {
        const response = await instance.acquireTokenPopup(request)
        console.log('acquireTokenPopup', response)
        return response
      } catch (error) {
        throw new Error('catch.acquireTokenPopup', { cause: error })
      }
    }
  }

  const createAPIKey = async (): Promise<APIKey> => {
    const response = await getTokenPopup()
    const userClient = userClientWithCustomToken(response.idToken)
    const apiResp = await userClient.createApiKey()
    if ('message' in apiResp) {
      throw new Error(`error creating a new API key ${apiResp.message}`)
    }
    return apiResp
  }

  const signOut = async (): Promise<void> => {
    const logoutRequest = {
      postLogoutRedirectUri: msalConfig.auth.redirectUri,
      mainWindowRedirectUri: msalConfig.auth.logoutUri,
    }
    const userClient = userClientWithCustomToken('') // Logging out doesn't require auth.
    try {
      await Promise.all([
        userClient.logout(),
        instance.logoutPopup(logoutRequest),
      ])
    } catch (e) {
      console.log('failed to log out', e)
    }
    isAuthenticated.value = false
    await router.push('/')
  }

  // Initialize the MSAL client if we're in the browser.
  if (process.client && !clientInitialized.value) {
    clientInitialized.value = true
    try {
      await instance.initialize()
      if (accounts.value.length === 0) {
        try {
          await instance.ssoSilent({})
        } catch (e) {
          console.log('failed to init SSO silently', e)
        }
      }
      const response = await instance.handleRedirectPromise()
      if (response !== null) {
        const claims = response.idTokenClaims
        if (!('tfp' in claims) || typeof claims.tfp !== 'string') {
          throw new Error('failed to find \'tfp\' claim')
        }
        const tfp = claims.tfp
        /**
             * For the purpose of setting an active account for UI update, we want to consider only the auth response resulting
             * from SUSI flow. "tfp" claim in the id token tells us the policy (NOTE: legacy policies may use "acr" instead of "tfp").
             * To learn more about B2C tokens, visit https://docs.microsoft.com/en-us/azure/active-directory-b2c/tokens-overview
             */
        if (tfp.toUpperCase() !== b2cPolicies.names.signUpSignIn.toUpperCase()) {
          throw new Error(`unexpected 'tfp' claim '${tfp.toUpperCase()}' does not match '${b2cPolicies.names.signUpSignIn.toUpperCase()}'`)
        }
        await handleResponse(response)
      } else {
        console.log('not coming from a redirect')
      }
    } catch (e) {
      console.log(e)
    }
  }

  await selectAccount()

  return {
    signIn,
    signOut,
    createAPIKey,
    isAuthenticated: computed(() => isAuthenticated.value),
  }
}
