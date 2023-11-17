import { UserClient } from '@/openapi/generated/user'
import { PACTAClient } from '@/openapi/generated/pacta'

import type { BaseHttpRequest } from '@/openapi/generated/pacta/core/BaseHttpRequest'
import type { OpenAPIConfig } from '@/openapi/generated/pacta/core/OpenAPI'

type HttpRequestConstructor = new (config: OpenAPIConfig) => BaseHttpRequest

// Note: This is a low-level composable intended to be used by other composables
// like usePACTA or the $msal plugin, it probably shouldn't be used by end
// clients.
export const useAPI = () => {
  const { public: { apiServerURL, authServerURL } } = useRuntimeConfig()

  const baseCfg = {
    CREDENTIALS: 'include' as const, // To satisfy typing of 'include' | 'same-origin' | etc
    WITH_CREDENTIALS: true,
  }

  const pactaCfg = {
    ...baseCfg,
    BASE: apiServerURL,
  }

  return {
    // The three different PACTA clients are for authentication in different
    // cases (client/server, cookies/no cookies, etc).
    pactaClient: new PACTAClient(pactaCfg).default,
    pactaClientWithHttpRequestClass: (req: HttpRequestConstructor) => {
      return new PACTAClient(pactaCfg, req).default
    },
    pactaClientWithAuth: (tkn: string) => {
      return new PACTAClient({
        ...pactaCfg,
        TOKEN: tkn,
      }).default
    },
    // Auth for the user service comes from Azure and needs to be manually
    // appended to each UserService request.
    userClientWithAuth: (tkn: string) => {
      const newCfg = {
        ...baseCfg,
        BASE: authServerURL,
        TOKEN: tkn,
      }
      return new UserClient(newCfg).default
    },
  }
}
