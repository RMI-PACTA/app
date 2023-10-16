<script setup lang="ts">
// const router = useRouter()
import { humanReadableDateFromStandardString } from '@/lib/time'

const pactaClient = await usePACTA()
// const { loading: { withLoading } } = useModal()
const { fromParams } = useURLParams()
// const localePath = useLocalePath()

const id = presentOrCheckURL(fromParams('id'))
const prefix = `initiative/${id}/invitations`
const [
  { data: initiative },
] = await Promise.all([
  useSimpleAsyncData(`${prefix}.getInitiative`, () => pactaClient.findInitiativeById(id)),
  useSimpleAsyncData(`${prefix}.getInvitations`, () => pactaClient.listInitiativeInvitations(id)),
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
        :code="initiative.language"
        class="inline"
      />
    </div>
    <div
      v-if="initiative.affiliation"
    >
      Created At: <b>{{ humanReadableDateFromStandardString(initiative.createdAt) }}</b>
    </div>
    {{ initiative.publicDescription }}
    <StandardDebug
      :value="initiative"
      label="Initiative"
    />
  </div>
</template>
