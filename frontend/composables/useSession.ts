import { type User } from '@/openapi/generated/pacta'

export const useSession = () => {
  const isAuthenticated = useIsAuthenticated()
  const pactaClient = usePACTA()

  const prefix = 'useSession'
  const currentUser = useState<User | undefined>(`${prefix}.currentUser`, () => undefined)
  const isAdmin = computed<boolean>(() => !!currentUser.value && currentUser.value.admin)
  const isSuperAdmin = computed<boolean>(() => !!currentUser.value && currentUser.value.superAdmin)

  const resolvers = useState<Array<() => void>>(`${prefix}.resolvers`, () => [])
  const loadCurrentUser = (hardRefresh = false): Promise<void> => {
    // Return the cached user if we don't explicitly want a new one
    if (currentUser.value && !hardRefresh) {
      return Promise.resolve()
    }

    // We're already loading a user, wait with everyone else
    if (resolvers.value.length > 0) {
      return new Promise((resolve) => {
        resolvers.value.push(resolve)
      })
    }

    // We're the first to request a user, kick of the request and hop in line at the front of the queue.
    return new Promise((resolve) => {
      resolvers.value.push(resolve)
      void pactaClient.findUserByMe()
        .then(m => {
          currentUser.value = m

          // Let everyone else know we've loaded the user and clear the queue.
          resolvers.value.forEach((fn) => { fn() })
          resolvers.value = []
        })
    })
  }
  const getMe = async () => {
    await loadCurrentUser()
    return {
      // LoadCurrentUser's return is only undefined as a technicality to support
      // the single-lookup behavior above. This cast is safe.
      me: currentUser as Ref<User>,
      isAdmin,
      isSuperAdmin,
    }
  }
  const getMaybeMe = async () => {
    if (isAuthenticated.value) {
      await loadCurrentUser()
    }
    return {
      // Will be a Ref with a value of undefined if the user isn't logged in.
      maybeMe: currentUser,
      isAdmin,
      isSuperAdmin,
    }
  }

  const refreshMaybeMe = async () => {
    if (isAuthenticated.value) {
      await loadCurrentUser(true)
    }
  }

  return {
    isAuthenticated,
    getMe,
    getMaybeMe,
    refreshMaybeMe,
    currentUser,
  }
}
