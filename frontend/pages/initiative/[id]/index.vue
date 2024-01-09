<script setup lang="ts">
import { languageToOption } from '@/lib/language'

const pactaClient = usePACTA()
const { fromParams } = useURLParams()
const { humanReadableDateFromStandardString } = useTime()
const { getMaybeMe } = useSession()
const { loading: { withLoading } } = useModal()

const { maybeMe } = await getMaybeMe()

const id = presentOrCheckURL(fromParams('id'))
const prefix = `initiative/${id}`
const [
  { data: initiative },
  { data: relationships, refresh: refreshRelationships },
] = await Promise.all([
  useSimpleAsyncData(`${prefix}.getInitiative`, () => pactaClient.findInitiativeById(id)),
  useSimpleAsyncData(`${prefix}.getRelationships`, () => pactaClient.listInitiativeUserRelationshipsByInitiative(id)),
])

const status = computed(() => {
  const i = initiative.value
  if (i.isAcceptingNewMembers) {
    if (i.isAcceptingNewPortfolios) {
      return 'Open'
    }
    return 'Accepting Portfolios from Existing Members'
  } else {
    return 'Closed'
  }
})

const isAMember = computed(() => {
  const mm = maybeMe.value
  if (!mm) return false
  return relationships.value.some((r) => r.userId === mm.id && r.member)
})
const canJoin = computed(() => {
  const i = initiative.value
  return maybeMe.value && i.isAcceptingNewMembers && !i.requiresInvitationToJoin && !isAMember.value
})
const join = () => {
  void withLoading(async () => {
    await pactaClient.updateInitiativeUserRelationship(
      id,
      presentOrFileBug(maybeMe.value).id,
      { member: true },
    )
    await refreshRelationships()
  }, 'initiative/join')
}
</script>

<template>
  <div class="flex flex-column gap-3">
    <div
      v-if="initiative.affiliation"
    >
      Sponsored by: <b>{{ initiative.affiliation }}</b>
    </div>
    <div>
      Status: <b>{{ status }}</b>
    </div>
    <div class="flex gap-2">
      Language: <LanguageRepresentation
        :code="languageToOption(initiative.language).code"
        class="inline"
      />
    </div>
    <div
      v-if="initiative.affiliation"
    >
      Created At: <b>{{ humanReadableDateFromStandardString(initiative.createdAt) }}</b>
    </div>
    {{ initiative.publicDescription }}
    <PVButton
      v-if="isAMember"
      disabled
      label="You are a member of this initiative"
      icon="pi pi-check"
    />
    <PVButton
      v-if="isAMember"
      label="Leave Initiative"
      icon="pi pi-arrow-left"
      class="p-button-danger p-button-outlined"
    />
    <PVButton
      v-if="canJoin"
      label="Join Initiative"
      icon="pi pi-arrow-right"
      icon-pos="right"
      @click="join"
    />
    <StandardDebug
      :value="initiative"
      label="Initiative"
    />
  </div>
</template>
