// Based heavily on [1], adapted for Nuxt 3. This was used instead of [2]
// because the Vue 3 sample doesn't account for SSR, and is generally more
// complicated than we need.
// [1] https://github.com/Azure-Samples/ms-identity-b2c-javascript-spa
// [2] https://github.com/AzureAD/microsoft-authentication-library-for-js/tree/dev/samples/msal-browser-samples/vue3-sample-app
import {
  type AccountInfo,
  type EventMessage,
  EventMessageUtils,
  EventType,
  InteractionRequiredAuthError,
  type InteractionStatus,
  type SilentRequest,
  PublicClientApplication,
} from '@azure/msal-browser'
import { type AuthenticationResult } from '@azure/msal-common'

import { computed } from 'vue'
import { type APIKey } from '~/openapi/generated/user'

export const useMSAL = async () => {
  const isAuthenticated = useState('useMSAL.isAuthenticated', () => false)

  // Don't initialize the MSAL client if we're not in the browser.
  if (process.server) {
    const jwt = useCookie('jwt')
    isAuthenticated.value = !!jwt.value
    return {
      signIn: () => Promise.reject(new Error('cannot call signIn on server')),
      signOut: () => Promise.reject(new Error('cannot call signOut on server')),
      createAPIKey: () => Promise.reject(new Error('cannot call createAPIKey on server')),
      getToken: () => Promise.reject(new Error('cannot call getToken on server')),
      isAuthenticated: computed(() => !!jwt.value),
    }
  }

  const router = useRouter()
  const { userClientWithAuth } = useAPI()

  const { $msal: { msalConfig, b2cPolicies } } = useNuxtApp()
  const scopes: string[] = ['openid', 'profile', 'offline_access', msalConfig.auth.clientId]

  const accounts = useState<AccountInfo[] | undefined>('useMSAL.accounts')
  const interactionStatus = useState<InteractionStatus | undefined>('useMSAL.interactionStatus')
  const instance = useState<PublicClientApplication | undefined>('useMSAL.instance')

  const handleResponse = async (response: AuthenticationResult, force = false): Promise<AuthenticationResult> => {
    if (!instance.value) {
      return await Promise.reject(new Error('MSAL instance was not yet initialized'))
    }

    if (!response?.account) {
      return await Promise.resolve(response)
    }

    // If this came from the cache, we don't need to refresh our token.
    if (response.fromCache && !force) {
      return response
    }

    accounts.value = [response.account]
    instance.value.setActiveAccount(response.account)
    const userClient = userClientWithAuth(response.idToken)
    try {
      await userClient.login()
      isAuthenticated.value = true
    } catch (error) {
      console.log('error at log in, signing out user', error)
      await signOut()
    }
    return response
  }

  const initializeAndAttemptLogin = async () => {
    if (instance.value) {
      console.log('instance is already initialized, returning')
      return
    }

    const inst = new PublicClientApplication(msalConfig)

    inst.addEventCallback((message: EventMessage) => {
      switch (message.eventType) {
        case EventType.ACCOUNT_ADDED:
        case EventType.ACCOUNT_REMOVED:
        case EventType.LOGIN_SUCCESS:
        case EventType.SSO_SILENT_SUCCESS:
        case EventType.HANDLE_REDIRECT_END:
        case EventType.LOGIN_FAILURE:
        case EventType.SSO_SILENT_FAILURE:
        case EventType.LOGOUT_END:
        case EventType.ACQUIRE_TOKEN_SUCCESS:
        case EventType.ACQUIRE_TOKEN_FAILURE:
          accounts.value = inst.getAllAccounts()
          break
      }

      const status = EventMessageUtils.getInteractionStatusFromEvent(message, interactionStatus.value)
      if (status !== null) {
        interactionStatus.value = status
      }
    })

    try {
      console.log('initializing MSAL client')
      await inst.initialize()
      instance.value = inst
    } catch (error) {
      console.log('failed to init instance', error)
      return
    }

    // See https://github.com/AzureAD/microsoft-authentication-library-for-js/blob/dev/lib/msal-browser/docs/initialization.md#handling-app-launch-with-0-or-more-available-accounts
    const accts = instance.value.getAllAccounts()
    if (accts.length === 0) {
      try {
        const authRes = await instance.value.ssoSilent({})
        await handleResponse(authRes)
        isAuthenticated.value = true
      } catch (error) {
        console.log('failed to init SSO silently', error)
      }
    } else if (accts.length === 1) {
      try {
        const authRes = await instance.value.acquireTokenSilent({
          scopes,
          account: accts[0],
        })
        await handleResponse(authRes)
        isAuthenticated.value = true
      } catch (error) {
        console.log('failed to acquire token silently', error)
      }
    } else {
      // When we handle this, use instance.setActiveAccount
      console.log('multiple accounts found, user needs to select one')
    }
  }

  const resolvers = useState<Array<() => void>>('useMSAL.resolvers', () => [])
  const loadMSAL = (): Promise<void> => {
    // We're already initialized
    if (instance.value) {
      return Promise.resolve()
    }

    // We're already initializing MSAL, wait with everyone else
    if (resolvers.value.length > 0) {
      return new Promise<void>((resolve) => {
        resolvers.value.push(resolve)
      })
    }

    // We're the first to request initializing MSAL, kick of the request and hop in line at the front of the queue.
    return new Promise<void>((resolve, reject) => {
      resolvers.value.push(resolve)
      initializeAndAttemptLogin()
        .then(() => {
          // Let everyone else know we've loaded the user and clear the queue.
          resolvers.value.forEach((fn) => { fn() })
          resolvers.value = []
        })
        .catch(reject)
    })
  }

  await loadMSAL()

  const account = computed(() => {
    if (!accounts.value) {
      return undefined
    }
    if (accounts.value.length < 1) {
      return undefined
    }
    if (accounts.value.length === 1) {
      return accounts.value[0]
    }
    /**
       * Due to the way MSAL caches account objects, the auth response from initiating a user-flow
       * is cached as a new account, which results in more than one account in the cache. Here we make
       * sure we are selecting the account with homeAccountId that contains the sign-up/sign-in user-flow,
       * as this is the default flow the user initially signed-in with.
       */
    const filteredAccounts = accounts.value.filter((account) => {
      return account.idTokenClaims && account.idTokenClaims.iss && account.idTokenClaims.aud &&
          account.homeAccountId.toUpperCase().includes(b2cPolicies.names.signUpSignIn.toUpperCase()) &&
          account.idTokenClaims.iss.toUpperCase().includes(b2cPolicies.authorityDomain.toUpperCase()) &&
          account.idTokenClaims.aud === msalConfig.auth.clientId
    })

    if (filteredAccounts.length === 0) {
      console.log('no accounts left after filtering', accounts.value, filteredAccounts, b2cPolicies)
      return undefined
    }

    if (filteredAccounts.length === 1) {
      return filteredAccounts[0]
    }

    // localAccountId identifies the entity for which the token asserts information.
    if (filteredAccounts.every((account) => account.localAccountId === filteredAccounts[0].localAccountId)) {
      // All filteredAccounts belong to the same user
      return filteredAccounts[0]
    }

    // Multiple users detected. Currently, we just return the first, but we should handle this more explicitly elsewhere.
    console.log('multiple users detected', filteredAccounts)
    return filteredAccounts[0]
  })

  const signIn = () => {
    if (!instance.value) {
      return Promise.reject(new Error('MSAL instance was not yet initialized'))
    }

    const req = { scopes }
    return instance.value.loginPopup(req)
      .then(handleResponse)
      .catch((err) => {
        console.log('useMSAL.loginPopup', err)
      })
  }

  const getToken = () => {
    if (!instance.value) {
      return Promise.reject(new Error('MSAL instance was not yet initialized'))
    }
    const inst = instance.value

    if (account.value === undefined) {
      // TODO: Figure out if this is a legitimate usecase.
      return Promise.reject(new Error('tried to get a token, but no account was found'))
    }

    const request: SilentRequest = {
      scopes,
      forceRefresh: false, // Set this to "true" to skip a cached token and go to the server to get a new token
      account: inst.getAccountByHomeId(account.value.homeAccountId) ?? undefined,
    }

    return inst.acquireTokenSilent(request)
      .then((response) => {
        // In case the response from B2C server has an empty idToken field
        // throw an error to initiate token acquisition
        if (response.idToken === '') {
          throw new InteractionRequiredAuthError()
        }
        return response
      })
      .then(handleResponse)
  }

  const createAPIKey = (): Promise<APIKey> => {
    return getToken()
      .then((response) => {
        const userClient = userClientWithAuth(response.idToken)
        return userClient.createApiKey()
      })
      .then((resp) => {
        if ('message' in resp) {
          throw new Error(`error creating a new API key ${resp.message}`)
        }
        return resp
      })
  }

  const signOut = (): Promise<void> => {
    if (!instance.value) {
      return Promise.reject(new Error('MSAL instance was not yet initialized'))
    }

    const logoutRequest = {
      postLogoutRedirectUri: msalConfig.auth.redirectUri,
      mainWindowRedirectUri: msalConfig.auth.logoutUri,
    }
    const userClient = userClientWithAuth('') // Logging out doesn't require auth.
    return Promise.all([
      userClient.logout(),
      instance.value.logoutPopup(logoutRequest),
    ])
      .catch((e) => { console.log('failed to log out', e) })
      .then(() => { /* cast to void */ })
      .finally(() => {
        isAuthenticated.value = false
        void router.push('/')
      })
  }

  return {
    signIn,
    signOut,
    createAPIKey,
    getToken,
    isAuthenticated: computed(() => isAuthenticated.value),
  }
}
