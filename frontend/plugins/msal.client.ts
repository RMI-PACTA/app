import { BrowserCacheLocation, LogLevel } from '@azure/msal-browser'

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
              console.info(message)
              return
            case LogLevel.Verbose:
              console.debug(message)
              return
            case LogLevel.Warning:
              console.warn(message)
              return
            case LogLevel.Trace:
              console.trace(message)
              return
          }
        },
        logLevel: toLogLevel(minLogLevel),
      },
    },
  }

  return {
    provide: {
      msal: {
        b2cPolicies,
        msalConfig,
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
