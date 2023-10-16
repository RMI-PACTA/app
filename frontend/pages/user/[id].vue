<script setup lang="ts">
const pactaClient = await usePACTA()
const { fromParams } = useURLParams()

const id = presentOrCheckURL(fromParams('id'))
const prefix = `user/[${id}]`
const [
  { data: user },
] = await Promise.all([
  useSimpleAsyncData(`${prefix}.findUserById`, () => pactaClient.findUserById(id)),
])
</script>

<template>
  <StandardContent v-if="user">
    <TitleBar :title="`User: ${user.name || user.id}`" />
    <UserToolbar :user-id="id" />
    <NuxtPage />
    <StandardDebug
      :value="user"
      label="User"
    />
  </StandardContent>
</template>
