import { type EventMessage, EventMessageUtils, EventType, InteractionStatus, PublicClientApplication, type AccountInfo, LogLevel } from '@azure/msal-browser'

type AccountIdentifiers = Partial<Pick<AccountInfo, 'homeAccountId' | 'localAccountId' | 'username'>>

export default defineNuxtPlugin((_nuxtApp) => {
  const {
    public: {
      msalConfig: {
        userFlowName,
        userFlowAuthority,
        authorityDomain,
        clientID,
        redirectURI,
        logoutURI,
      }
    }
  } = useRuntimeConfig()

  const b2cPolicies = {
    names: {
      signUpSignIn: userFlowName
    },
    authorities: {
      signUpSignIn: {
        authority: userFlowAuthority
      }
    },
    authorityDomain
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
      cacheLocation: 'localStorage',
      storeAuthStateInCookie: true
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
              console.info(message)
              return
            case LogLevel.Verbose:
              console.debug(message)
              return
            case LogLevel.Warning:
              console.warn(message)
          }
        },
        logLevel: LogLevel.Verbose
      }
    }
  }
  const msalInstance = new PublicClientApplication(msalConfig)

  const inProgress = InteractionStatus.Startup
  const accounts = msalInstance.getAllAccounts()

  const stateInProgress = useState<InteractionStatus>('msal.inProgress', () => inProgress)
  const stateAccounts = useState('msal.accounts', () => accounts)

  msalInstance.addEventCallback((message: EventMessage) => {
    const maybeUpdateAccounts = (): void => {
      const currentAccounts = msalInstance.getAllAccounts()
      if (!accountArraysAreEqual(currentAccounts, stateAccounts.value)) {
        stateAccounts.value = currentAccounts
      }
    }
    switch (message.eventType) {
      case EventType.ACCOUNT_ADDED:
        // fallthrough
      case EventType.ACCOUNT_REMOVED:
        // fallthrough
      case EventType.LOGIN_SUCCESS:
        // fallthrough
      case EventType.SSO_SILENT_SUCCESS:
        // fallthrough
      case EventType.HANDLE_REDIRECT_END:
        // fallthrough
      case EventType.LOGIN_FAILURE:
        // fallthrough
      case EventType.SSO_SILENT_FAILURE:
        // fallthrough
      case EventType.LOGOUT_END:
        // fallthrough
      case EventType.ACQUIRE_TOKEN_SUCCESS:
        // fallthrough
      case EventType.ACQUIRE_TOKEN_FAILURE:
        maybeUpdateAccounts()
        break
    }

    const status = EventMessageUtils.getInteractionStatusFromEvent(message, stateInProgress.value)
    if (status !== null) {
      stateInProgress.value = status
    }
  })

  return {
    provide: {
      msal: {
        inProgress: stateInProgress,
        instance: msalInstance,
        accounts: stateAccounts,
        b2cPolicies,
        msalConfig
      }
    }
  }
})

/**
 * Helper function to determine whether 2 arrays are equal
 * Used to avoid unnecessary state updates
 * @param arrayA
 * @param arrayB
 */
function accountArraysAreEqual (arrayA: AccountIdentifiers[], arrayB: AccountIdentifiers[]): boolean {
  if (arrayA.length !== arrayB.length) {
    return false
  }

  const comparisonArray = [...arrayB]

  return arrayA.every((elementA) => {
    const elementB = comparisonArray.shift()
    if (elementA === undefined || elementB === undefined) {
      return false
    }

    return (elementA.homeAccountId === elementB.homeAccountId) &&
               (elementA.localAccountId === elementB.localAccountId) &&
               (elementA.username === elementB.username)
  })
}
