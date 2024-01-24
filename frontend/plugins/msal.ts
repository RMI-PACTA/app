import {
  BrowserCacheLocation, LogLevel,
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

export default defineNuxtPlugin(async (_nuxtApp) => {
  const isAuthenticated = useIsAuthenticated()
  if (process.server) {
    const jwt = useCookie('jwt')
    isAuthenticated.value = !!jwt.value
    return {
      provide: {
        msal: {
          msalSignIn: () => Promise.reject(new Error('cannot call signIn on server')),
          signOut: () => Promise.reject(new Error('cannot call signOut on server')),
          createAPIKey: () => Promise.reject(new Error('cannot call createAPIKey on server')),
          getToken: () => Promise.reject(new Error('cannot call getToken on server')),
          isAuthenticated: computed(() => !!jwt.value),
        },
      },
    }
  }

  const {
    public: {
      msalConfig: {
        userFlowName,
        userFlowAuthority,
        authorityDomain,
        clientID,
        redirectURI,
        logoutURI,
        minLogLevel,
      },
    },
  } = useRuntimeConfig()
  const b2cPolicies = {
    names: {
      signUpSignIn: userFlowName,
    },
    authorities: {
      signUpSignIn: {
        authority: userFlowAuthority,
      },
    },
    authorityDomain,
  }
  const msalConfig = {
    auth: {
      clientId: clientID,
      authority: b2cPolicies.authorities.signUpSignIn.authority, // Choose sign-up/sign-in user-flow as your default.
      knownAuthorities: [b2cPolicies.authorityDomain], // You must identify your tenant's domain as a known authority.
      redirectUri: redirectURI, // Must be registered as a SPA redirectURI on your app registration
      logoutUri: logoutURI,
    },
    cache: {
      cacheLocation: BrowserCacheLocation.LocalStorage,
      claimsBasedCachingEnabled: true,
      storeAuthStateInCookie: false,
    },
    system: {
      loggerOptions: {
        loggerCallback: (level: LogLevel, message: string, containsPii: boolean) => {
          if (containsPii) {
            return
          }
          switch (level) {
            case LogLevel.Error:
              console.error(message)
              return
            case LogLevel.Info:
              // console.info(message)
              return
            case LogLevel.Verbose:
              // console.debug(message)
              return
            case LogLevel.Warning:
              console.warn(message)
              return
            case LogLevel.Trace:
              // console.trace(message)
              return
          }
        },
        logLevel: toLogLevel(minLogLevel),
      },
    },
  }

  const scopes: string[] = ['openid', 'profile', 'offline_access', msalConfig.auth.clientId]

  const router = useRouter()
  const { userClientWithAuth } = useAPI()
  const localePath = useLocalePath()

  const accounts = useState<AccountInfo[] | undefined>('useMSAL.accounts')
  const interactionStatus = useState<InteractionStatus | undefined>('useMSAL.interactionStatus')

  const instance = new PublicClientApplication(msalConfig)

  instance.addEventCallback((message: EventMessage) => {
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
        accounts.value = instance.getAllAccounts()
        break
    }

    const status = EventMessageUtils.getInteractionStatusFromEvent(message, interactionStatus.value)
    if (status !== null) {
      interactionStatus.value = status
    }
  })

  try {
    console.log('initializing MSAL client')
    await instance.initialize()
  } catch (error) {
    throw new Error('failed to init MSAL instance', { cause: error })
  }

  const signOut = (): Promise<void> => {
    const logoutRequest = {
      postLogoutRedirectUri: msalConfig.auth.redirectUri,
      mainWindowRedirectUri: msalConfig.auth.logoutUri,
    }
    const userClient = userClientWithAuth('') // Logging out doesn't require auth.
    return Promise.all([
      userClient.logout(),
      instance.logoutPopup(logoutRequest),
    ])
      .catch((e) => { console.log('failed to log out', e) })
      .then(() => { /* cast to void */ })
      .finally(() => {
        isAuthenticated.value = false
        void router.push(localePath('/'))
      })
  }

  const handleResponse = async (response: AuthenticationResult, force = false): Promise<AuthenticationResult> => {
    if (!response?.account) {
      return await Promise.resolve(response)
    }

    // If this came from the cache, we don't need to refresh our token.
    if (response.fromCache && !force) {
      return response
    }

    accounts.value = [response.account]
    instance.setActiveAccount(response.account)
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

  // See https://github.com/AzureAD/microsoft-authentication-library-for-js/blob/dev/lib/msal-browser/docs/initialization.md#handling-app-launch-with-0-or-more-available-accounts
  const accts = instance.getAllAccounts()
  if (accts.length === 0) {
    try {
      const authRes = await instance.ssoSilent({})
      await handleResponse(authRes)
      isAuthenticated.value = true
    } catch (error) {
      console.log('failed to init SSO silently', error)
    }
  } else if (accts.length === 1) {
    try {
      const authRes = await instance.acquireTokenSilent({
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

  const getToken = (): Promise<AuthenticationResult | undefined> => {
    if (account.value === undefined) {
      return new Promise<AuthenticationResult | undefined>((resolve) => { resolve(undefined) })
    }

    const request: SilentRequest = {
      scopes,
      forceRefresh: false, // Set this to "true" to skip a cached token and go to the server to get a new token
      account: instance.getAccount({ homeAccountId: account.value.homeAccountId }) ?? undefined,
    }

    return instance.acquireTokenSilent(request)
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

  const msalSignIn: () => Promise<void> = (): Promise<void> => {
    const req = { scopes }
    return instance.loginPopup(req)
      .then(handleResponse)
      .then(() => { /* cast to void */ })
      .catch((err) => {
        console.log('useMSAL.loginPopup', err)
      })
  }

  const createAPIKey = async (): Promise<APIKey | undefined> => {
    const token = await getToken()
    if (!token) {
      return undefined
    }
    const apiKey = await userClientWithAuth(token.idToken).createApiKey()
    if ('message' in apiKey) {
      throw new Error(`error creating a new API key ${apiKey.message}`)
    }
    return apiKey
  }

  return {
    provide: {
      msal: {
        msalSignIn,
        signOut,
        createAPIKey,
        getToken,
        isAuthenticated: computed(() => isAuthenticated.value),
      },
    },
  }
})

function toLogLevel (lvl: string): LogLevel {
  switch (lvl) {
    case 'ERROR':
      return LogLevel.Error
    case 'WARNING':
      return LogLevel.Warning
    case 'INFO':
      return LogLevel.Info
    case 'VERBOSE':
      return LogLevel.Verbose
    case 'TRACE':
      return LogLevel.Trace
    default:
      return LogLevel.Verbose
  }
}
