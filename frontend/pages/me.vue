<script setup lang="ts">
const { getMaybeMe } = useSession()

const { maybeMe, isAdmin, isSuperAdmin } = await getMaybeMe()
</script>

<template>
  <StandardContent>
    {{ maybeMe }}
    {{ isAdmin }}
    {{ isSuperAdmin }}
  </StandardContent>
</template>
