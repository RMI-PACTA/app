<script setup lang="ts">
const { fromParams } = useURLParams()
const { t } = useI18n()
const localePath = useLocalePath()

const tt = (key: string) => t(`pages/initiative.${key}`)

const id = presentOrCheckURL(fromParams('id'))
const { initiative, canManage, isMember } = await useInitiativeData(id)

const menuItems = computed((): Array<{ to: string, label: string, icon: string }> => {
  if (!canManage.value && !isMember.value) {
    return []
  }
  const result = [
    {
      to: localePath(`/initiative/${id}`),
      label: tt('Initiative Home'),
      icon: 'pi pi-home',
    },
  ]
  if (canManage.value) {
    result.push(
      {
        to: localePath(`/initiative/${id}/edit`),
        label: tt('Edit'),
        icon: 'pi pi-pencil',
      },
      {
        to: localePath(`/initiative/${id}/invitations`),
        label: tt('Invitations'),
        icon: 'pi pi-envelope',
      },
      {
        to: localePath(`/initiative/${id}/relationships`),
        label: tt('Relationships'),
        icon: 'pi pi-users',
      },
      {
        to: localePath(`/initiative/${id}/portfolios`),
        label: tt('Portfolios'),
        icon: 'pi pi-copy',
      },
    )
  }
  if (canManage.value || isMember.value) {
    result.push(
      {
        to: localePath(`/initiative/${id}/internal`),
        label: tt('Internal Information'),
        icon: 'pi pi-info-circle',
      },
    )
  }
  return result
})
</script>

<template>
  <StandardContent v-if="initiative">
    <TitleBar :title="`${tt('Initiative')}: ${initiative.name}`" />
    <div
      v-if="menuItems.length > 0"
      class="p-buttonset"
    >
      <LinkButton
        v-for="item in menuItems"
        :key="item.to"
        :to="item.to"
        :label="item.label"
        :icon="item.icon"
        active-class="border-2"
        inactive-class="p-button-outlined"
      />
    </div>
    <NuxtPage />
    <StandardDebug
      :value="initiative"
      label="Initiative"
    />
  </StandardContent>
</template>
