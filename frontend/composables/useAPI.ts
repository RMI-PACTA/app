import { UserClient, type DefaultService as UserDefaultService } from '@/openapi/generated/user'
import { PACTAClient, type DefaultService as PACTADefaultService } from '@/openapi/generated/pacta'

interface API {
  userClient: UserDefaultService
  pactaClient: PACTADefaultService
  userClientWithCustomToken: (tkn: string) => UserDefaultService
}

export const useAPI = (): API => {
  const { public: { apiServerURL, authServerURL } } = useRuntimeConfig()
  const baseCfg = {
    CREDENTIALS: 'include' as const, // To satisfy typing of 'include' | 'same-origin' | etc
    WITH_CREDENTIALS: true
  }

  // If we're on the server, forward our cookie header along to the backend
  // API for auth. We don't do this for the UserClient because it uses separate
  // auth.
  let headers: Record<string, string> = {}
  if (process.server) {
    headers = useRequestHeaders(['cookie'])
  }

  const userCfg = {
    ...baseCfg,
    BASE: authServerURL
  }
  const userClient = new UserClient(userCfg)

  const pactaClient = new PACTAClient({
    ...baseCfg,
    BASE: apiServerURL,
    HEADERS: headers
  })

  return {
    userClient: userClient.default,
    pactaClient: pactaClient.default,
    userClientWithCustomToken: (tkn: string) => {
      const newCfg = {
        ...userCfg,
        TOKEN: tkn
      }
      return new UserClient(newCfg).default
    }
  }
}
