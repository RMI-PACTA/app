import { type AuthenticationResult } from '@azure/msal-common'
import type { ApiRequestOptions } from '~/openapi/generated/pacta/core/ApiRequestOptions'
import { BaseHttpRequest } from '~/openapi/generated/pacta/core/BaseHttpRequest'
import { CancelablePromise } from '~/openapi/generated/pacta/core/CancelablePromise'
import type { OpenAPIConfig } from '~/openapi/generated/pacta/core/OpenAPI'
import { request as __request } from '~/openapi/generated/pacta/core/request'

export const usePACTA = async () => {
  const {
    pactaClient, // On the server, if there's no JWT
    pactaClientWithAuth, // On the server, forward the cookie we got from the client
    pactaClientWithHttpRequestClass, // On the client, wrap with a check for a fresh cookie.
  } = useAPI()

  if (process.server) {
    const jwt = useCookie('jwt')
    if (jwt.value) {
      return pactaClientWithAuth(jwt.value)
    }
    return pactaClient
  }

  // If we're on the client, we can see if Azure has cached credentials and use
  // those, or refresh them if not. Our cookies have the same lifetime as our
  // access tokens, so we refresh them together.
  const { getToken } = await useMSAL()

  // We define this class as a variable so we can override the PACTA client
  // with middleware that refreshes our credentials. This matches the
  // interface of our auto-generated code, which expects a class that extends
  // BaseHttpRequest.
  const httpReqClass = class extends BaseHttpRequest {
    private readonly getToken: () => Promise<AuthenticationResult>

    constructor (config: OpenAPIConfig) {
      super(config)
      this.getToken = getToken
    }

    /**
     * Request method
     * @param options The request options from the service
     * @returns CancelablePromise<T>
     * @throws ApiError
     */
    public override request<T>(options: ApiRequestOptions): CancelablePromise<T> {
      return new CancelablePromise((resolve, reject, onCancel) => {
        this.getToken()
          .then(() => {
            const cancelablePromise = __request<T>(this.config, options)
            onCancel(() => {
              cancelablePromise.cancel()
            })
            return cancelablePromise
          })
          .then(resolve)
          .catch(reject)
      })
    }
  }

  return pactaClientWithHttpRequestClass(httpReqClass)
}
