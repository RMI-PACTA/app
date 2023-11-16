<script setup lang="ts">
const { loading: { withLoading } } = useModal()
const { fromParams } = useURLParams()
const localePath = useLocalePath()
const router = useRouter()
const { t } = useI18n()

const [
  { getMaybeMe },
  pactaClient,
] = await Promise.all([
  useSession(),
  usePACTA(),
])
const { maybeMe } = await getMaybeMe()

const id = presentOrCheckURL(fromParams('invitationId'))
const prefix = `join[${id}]`
const [
  { data: invitation },
] = await Promise.all([
  useSimpleAsyncData(`${prefix}.getInitiativeInvitation`, () => pactaClient.getInitiativeInvitation(id)),
])
const translationPrefix = 'pages/join'
const tt = (key: string) => t(`${translationPrefix}.${key}`)

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
    <TitleBar :title="`${tt('Join Initiative:')} '${invitation.initiativeId}'`" />
    <template v-if="!invitation.usedAt">
      <p>
        {{ tt('You\'ve been invited to join an initiative') }} <b>{{ invitation.initiativeId }}</b>.
        <NuxtLink :to="localePath(`/initiative/${invitation.initiativeId}`)">
          {{ tt('You can learn more about the initiative here.') }}
        </NuxtLink>
      </p>
      <p>
        {{ tt('if-accept') }}
      </p>
      <p>
        {{ tt('what-shared') }}
      </p>
      <p>{{ tt('do-accept' ) }}</p>
      <div class="flex gap-2">
        <LinkButton
          :to="localePath(`/initiative/${invitation.initiativeId}`)"
          class="p-button-secondary p-button-outlined"
          icon="pi pi-arrow-left"
          :label="tt('Not Now')"
        />
        <PVButton
          class="p-button-success"
          :label="tt('Accept Invitation')"
          icon="pi pi-check"
          icon-pos="right"
          @click="acceptInvitation"
        />
      </div>
    </template>
    <template v-else>
      {{ tt('This invitation has already been used.') }}
    </template>

    <StandardDebug
      :value="invitation"
      label="Invitation"
    />
  </StandardContent>
</template>
