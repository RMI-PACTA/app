<script setup lang="ts">
const router = useRouter()
const localePath = useLocalePath()
const { permissionDenied: { setPermissionDenied } } = useModal()
const { getMaybeMe } = useSession()
const { maybeMe } = await getMaybeMe()

onMounted(() => {
  const mm = maybeMe.value
  if (!mm) {
    setPermissionDenied(new Error("You aren't logged in, so we can't find your profile."))
    return
  }
  void router.push(localePath(`/user/${mm.id}`))
})
</script>

<template>
  <StandardContent>
    <TitleBar title="Redirecting to your profile..." />
  </StandardContent>
</template>
