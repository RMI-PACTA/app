import { type InitiativeUserRelationship, type InitiativeInvitation } from '@/openapi/generated/pacta'

export const useInitiativeData = async (id: string) => {
  const prefix = `useInitiativeData[${id}]`

  const pactaClient = usePACTA()
  const { getMaybeMe } = useSession()

  const { maybeMe, isAdmin, isSuperAdmin } = await getMaybeMe()

  const meRelationships = useState<InitiativeUserRelationship[]>(`${prefix}.meRelationships`, () => [])
  if (maybeMe.value) {
    meRelationships.value = await pactaClient.listInitiativeUserRelationshipsByUser(maybeMe.value.id)
  }

  const isManagerByMe = computed(() => meRelationships.value.some((r) => r.initiativeId === id && r.manager))

  const canManageByMe = computed(() => isAdmin.value || isSuperAdmin.value || isManagerByMe.value)

  const maybeLookUpInvitationsByInitiative = async (): Promise<InitiativeInvitation[]> => {
    if (!canManageByMe.value) {
      return []
    }
    return await pactaClient.listInitiativeInvitations(id)
  }

  const [
    { data: initiative, refresh: refreshInitiative },
    { data: invitations, refresh: refreshInvitations },
  ] = await Promise.all([
    useSimpleAsyncData(`${prefix}.getInitiative`, () => pactaClient.findInitiativeById(id)),
    useSimpleAsyncData(`${prefix}.getInvitations`, maybeLookUpInvitationsByInitiative),
  ])

  const isMember = computed(() => initiative.value.initiativeUserRelationships.some((r) => r.member))
  const isManager = computed(() => initiative.value.initiativeUserRelationships.some((r) => r.manager))
  const canManage = computed(() => canManageByMe.value || isManager.value)

  const canJoinIfLoggedIn = computed(() => !isMember.value && !isManager.value && initiative.value.isAcceptingNewMembers && !initiative.value.requiresInvitationToJoin)
  const canDirectlyJoin = computed(() => canJoinIfLoggedIn.value && maybeMe.value)

  return {
    initiative,
    refreshInitiative,
    invitations,
    refreshInvitations,
    canManage,
    isMember,
    isManager,
    canDirectlyJoin,
    canJoinIfLoggedIn,
  }
}
