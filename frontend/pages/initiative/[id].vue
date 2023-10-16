<script setup lang="ts">
const pactaClient = await usePACTA()
const { fromParams } = useURLParams()

const id = presentOrCheckURL(fromParams('id'))
const prefix = `initiative/${id}/invitations`
const [
  { data: initiative },
  { data: relationships },
] = await Promise.all([
  useSimpleAsyncData(`${prefix}.getInitiative`, () => pactaClient.findInitiativeById(id)),
  useSimpleAsyncData(`${prefix}.getRelationships`, () => pactaClient.listInitiativeUserRelationshipsByInitiative(id)),
])
</script>

<template>
  <StandardContent v-if="initiative">
    <TitleBar :title="`Initiative: ${initiative.name}`" />
    <InitiativeToolbar
      :initiative-id="initiative.id"
      :initiative-user-relationships="relationships"
    />
    <NuxtPage />
    <StandardDebug
      :value="initiative"
      label="Initiative"
    />
    <StandardDebug
      :value="relationships"
      label="Relationships"
    />
  </StandardContent>
</template>
