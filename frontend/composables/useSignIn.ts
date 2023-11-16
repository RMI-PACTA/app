export const useSignIn = async () => {
  const [
    { msalSignIn },
    pactaClient,
    { refreshMaybeMe },
  ] = await Promise.all([
    useMSAL(),
    usePACTA(),
    useSession(),
  ])

  const signIn = () => msalSignIn()
    .then(() => pactaClient.userAuthenticationFollowup())
    .then(() => refreshMaybeMe())

  return { signIn }
}
