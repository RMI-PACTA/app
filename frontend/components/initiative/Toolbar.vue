<script setup lang="ts">
import { type InitiativeUserRelationship } from '@/openapi/generated/pacta'

const localePath = useLocalePath()
const { t } = useI18n()
const { getMaybeMe } = await useSession()
const { isAdmin, maybeMe } = await getMaybeMe()

const prefix = 'InitiativeToolbar'
const tt = (key: string) => t(`${prefix}.${key}`)

interface Props {
  initiativeId: string
  initiativeUserRelationships: InitiativeUserRelationship[]
}
const props = defineProps<Props>()

const isManager = computed(() => {
  const mm = maybeMe.value
  return (!!mm && props.initiativeUserRelationships.some(r => r.manager && r.userId === mm.id))
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
      :label="tt('Initiative Home')"
      icon="pi pi-home"
      active-class="border-2"
      inactive-class="p-button-outlined"
    />
    <LinkButton
      v-if="canEdit"
      :to="localePath(`/initiative/${initiativeId}/edit`)"
      :label="tt('Edit')"
      icon="pi pi-pencil"
      active-class="border-2"
      inactive-class="p-button-outlined"
    />
    <LinkButton
      v-if="canEdit"
      :to="localePath(`/initiative/${initiativeId}/invitations`)"
      :label="tt('Invitations')"
      icon="pi pi-envelope"
      active-class="border-2"
      inactive-class="p-button-outlined"
    />
    <LinkButton
      v-if="canEdit"
      :to="localePath(`/initiative/${initiativeId}/relationships`)"
      :label="tt('Relationships')"
      icon="pi pi-users"
      active-class="border-2"
      inactive-class="p-button-outlined"
    />
    <LinkButton
      v-if="canSeeInternal"
      :to="localePath(`/initiative/${initiativeId}/internal`)"
      :label="tt('Internal Information')"
      icon="pi pi-info-circle"
      active-class="border-2"
      inactive-class="p-button-outlined"
    />
  </div>
</template>
