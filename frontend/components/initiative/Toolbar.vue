<script setup lang="ts">
import { type InitiativeUserRelationship } from '@/openapi/generated/pacta'

const localePath = useLocalePath()
const { getMaybeMe } = useSession()

const { isAdmin, maybeMe } = await getMaybeMe()

interface Props {
  initiativeId: string
  initiativeUserRelationships: InitiativeUserRelationship[]
}
const props = defineProps<Props>()

const isManager = computed(() => {
  const mm = maybeMe.value
  return (!!mm && props.initiativeUserRelationships.some(r => r.manager && r.userId === mm.id)) || true
})
const isMember = computed(() => {
  const mm = maybeMe.value
  return !!mm && props.initiativeUserRelationships.some(r => r.member && r.userId === mm.id)
})
const canEdit = computed<boolean>(() => isManager.value || isAdmin.value)
const canSeeInternal = computed<boolean>(() => isManager.value || isMember.value || isAdmin.value)
const showToolbar = computed<boolean>(() => canSeeInternal.value || canEdit.value)
</script>

<template>
  <div
    v-show="showToolbar"
    class="p-buttonset initiative-toolbar"
  >
    <LinkButton
      :to="localePath(`/initiative/${initiativeId}`)"
      label="Initiative Home"
      icon="pi pi-home"
      active-class="border-2"
      inactive-class="p-button-outlined"
    />
    <LinkButton
      v-if="canEdit"
      :to="localePath(`/initiative/${initiativeId}/edit`)"
      label="Edit"
      icon="pi pi-pencil"
      active-class="border-2"
      inactive-class="p-button-outlined"
    />
    <LinkButton
      v-if="canEdit"
      :to="localePath(`/initiative/${initiativeId}/invitations`)"
      label="Invitations"
      icon="pi pi-envelope"
      active-class="border-2"
      inactive-class="p-button-outlined"
    />
    <LinkButton
      v-if="canEdit"
      :to="localePath(`/initiative/${initiativeId}/relationships`)"
      label="Relationships"
      icon="pi pi-users"
      active-class="border-2"
      inactive-class="p-button-outlined"
    />
    <LinkButton
      v-if="canSeeInternal"
      :to="localePath(`/initiative/${initiativeId}/internal`)"
      label="Internal Information"
      icon="pi pi-info-circle"
      active-class="border-2"
      inactive-class="p-button-outlined"
    />
  </div>
</template>

<style scoped lang="scss">
.initiative-toolbar {
  .router-link-active {
    background: 'red';
 }

 /*
  .my-link.router-link-active:hover {
    @apply bg-green-200 font-medium;
  }
  */
}
</style>
