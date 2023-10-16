<script setup lang="ts">
const { loading: { withLoading } } = useModal()
const { fromParams } = useURLParams()
const localePath = useLocalePath()
const router = useRouter()
const { getMaybeMe } = useSession()
const pactaClient = await usePACTA()

const { maybeMe } = await getMaybeMe()

const id = presentOrCheckURL(fromParams('invitationId'))
const prefix = `join[${id}]`
const [
  { data: invitation },
] = await Promise.all([
  useSimpleAsyncData(`${prefix}.getInitiativeInvitation`, () => pactaClient.getInitiativeInvitation(id)),
  // useSimpleAsyncData(`${prefix}.getMe`, () => pactaClient.findUserByMe()),
])

if (invitation.value && maybeMe.value && invitation.value.usedByUserId === maybeMe.value.id) {
  void router.push(localePath(`/initiative/${invitation.value.initiativeId}/internal`))
}

const acceptInvitation = () => withLoading(
  () => pactaClient.claimInitiativeInvitation(id).then(() => router.push(localePath(`/initiative/${invitation.value.initiativeId}/internal`))),
  'Join/AcceptingInvitation',
)
</script>

<template>
  <StandardContent v-if="invitation">
    <TitleBar :title="`Join Initiative '${invitation.initiativeId}'`" />
    <template v-if="!invitation.usedAt">
      <p>
        You've been invited to join the initiative <b>{{ invitation.initiativeId }}</b>.
        <NuxtLink :to="localePath(`/initiative/${invitation.initiativeId}`)">
          You can learn more about that initiative here.
        </NuxtLink>
      </p>
      <p>
        If you accept this invitation, you'll be able to add your portfolios to the
        project, where they will be included in aggregated analysis.
        You aren't required to accept this invitation if you just want to run the
        analysis on your own, but joining the initiative may enable you to
        access internal benefits, like the ability to see the analysis results
        from the aggregated portfolios.
      </p>
      <p>
        If you accept X, Y, Z will be whared with the initiative administrators,
        but A, B, C will not be.
      </p>
      <p>Do you want to accept this invitation?</p>
      <div class="flex gap-2">
        <LinkButton
          :to="localePath(`/initiative/${invitation.initiativeId}`)"
          class="p-button-secondary p-button-outlined"
          icon="pi pi-arrow-left"
          label="Not Now"
        />
        <PVButton
          class="p-button-success"
          label="Accept Invitation"
          icon="pi pi-check"
          icon-pos="right"
          @click="acceptInvitation"
        />
      </div>
    </template>
    <template v-else>
      This invitation has already been used.
    </template>

    <StandardDebug
      :value="invitation"
      label="Invitation"
    />
  </StandardContent>
</template>
