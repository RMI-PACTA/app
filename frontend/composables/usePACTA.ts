import { type AuthenticationResult } from '@azure/msal-common'
import type { ApiRequestOptions } from '~/openapi/generated/pacta/core/ApiRequestOptions'
import { BaseHttpRequest } from '~/openapi/generated/pacta/core/BaseHttpRequest'
import { CancelablePromise } from '~/openapi/generated/pacta/core/CancelablePromise'
import type { OpenAPIConfig } from '~/openapi/generated/pacta/core/OpenAPI'
import { request as __request } from '~/openapi/generated/pacta/core/request'
import { serializeError } from 'serialize-error'
import { createErrorWithRemediation, isSilent, Remediation } from '@/lib/error'

export const usePACTA = () => {
  const {
    pactaClient, // On the server, if there's no JWT
    pactaClientWithAuth, // On the server, forward the cookie we got from the client
    pactaClientWithHttpRequestClass, // On the client, wrap with a check for a fresh cookie.
  } = useAPI()

  const { loading: { clearLoading }, error: { errorModalVisible, error } } = useModal()

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
  const { $msal: { getToken } } = useNuxtApp()

  // We define this class as a variable so we can override the PACTA client
  // with middleware that refreshes our credentials. This matches the
  // interface of our auto-generated code, which expects a class that extends
  // BaseHttpRequest.
  const httpReqClass = class extends BaseHttpRequest {
    private readonly getToken: () => Promise<AuthenticationResult | undefined>

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
          .catch((e: unknown) => {
            // We know we're at the "bottom" of the call chain here in usePACTA
            // and so nothing will be sending a silent error, but we don't know
            // how the code will change in the future, so we check just to be
            // safe, lest we accidentally overlay multiple error modals with
            // less useful information.
            if (!isSilent(e)) {
              error.value = serializeError(e)
              errorModalVisible.value = true
              clearLoading()
            }

            const err = createErrorWithRemediation(`${options.method} ${options.url} failed`, Remediation.Silent)
            reject(err)
          })
      })
    }
  }

  return pactaClientWithHttpRequestClass(httpReqClass)
}
