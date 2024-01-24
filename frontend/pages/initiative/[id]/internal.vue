<script setup lang="ts">
const { fromParams } = useURLParams()

const id = presentOrCheckURL(fromParams('id'))

const { initiative } = await useInitiativeData(id)
</script>

<template>
  <div>
    {{ initiative.internalDescription }}
  </div>
</template>
