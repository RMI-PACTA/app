export const useSignIn = () => {
  const { $msal: { msalSignIn } } = useNuxtApp()
  const pactaClient = usePACTA()
  const { refreshMaybeMe } = useSession()

  const signIn = () => msalSignIn()
    .then(() => pactaClient.userAuthenticationFollowup())
    .then(() => refreshMaybeMe())

  return { signIn }
}
