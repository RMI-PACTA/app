<script setup lang="ts">
// const router = useRouter()
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
</script>

<template>
  <div>
    {{ initiative.internalDescription }}
  </div>
</template>
