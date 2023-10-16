<script setup lang="ts">
const { getMaybeMe } = useSession()

const localePath = useLocalePath()
const { isAdmin, maybeMe } = await getMaybeMe()

interface Props {
  userId: string
}
const props = defineProps<Props>()

const isMe = computed<boolean>(() => {
  const mm = maybeMe.value
  return !!mm && props.userId === mm.id
})
const canEdit = computed<boolean>(() => isMe.value || isAdmin.value)
const showToolbar = computed<boolean>(() => canEdit.value)
</script>

<template>
  <div
    v-show="showToolbar"
    class="p-buttonset"
  >
    <LinkButton
      :to="localePath(`/user/${props.userId}`)"
      label="Profile"
      icon="pi pi-home"
      active-class="border-2"
      inactive-class="p-button-outlined"
    />
    <LinkButton
      v-if="canEdit"
      :to="localePath(`/user/${props.userId}/edit`)"
      label="Edit"
      icon="pi pi-pencil"
      active-class="border-2"
      inactive-class="p-button-outlined"
    />
  </div>
</template>
