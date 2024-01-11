<script setup lang="ts">

const { fromParams } = useURLParams()
const localePath = useLocalePath()
const { humanReadableTimeFromStandardString } = useTime()
const { loading: { withLoading } } = useModal()
const { t } = useI18n()

const tt = (key: string) => t(`pages/initiative/relationships.${key}`)

const [
  pactaClient,
  { getMaybeMe },
] = await Promise.all([
  usePACTA(),
  useSession(),
])
const { maybeMe, isAdmin } = await getMaybeMe()

const id = presentOrCheckURL(fromParams('id'))
const prefix = `initiative/${id}/relationships`
const [
  { data: relationships, refresh: refreshRelationships },
] = await Promise.all([
  useSimpleAsyncData(`${prefix}.getRelationships`, () => pactaClient.listInitiativeUserRelationshipsByInitiative(id)),
])
const nonEmptyRelationships = computed(() => relationships.value.filter((r) => r.manager || r.member))
const canManage = computed(() => {
  const mm = maybeMe.value
  if (!mm) return false
  if (isAdmin.value) return true
  return relationships.value.some((r) => r.userId === mm.id && r.manager)
})

const changeMembership = (userId: string, member: boolean | undefined, manager: boolean | undefined) =>
  withLoading(async () => {
    await pactaClient.updateInitiativeUserRelationship(
      id,
      userId,
      { member, manager },
    )
    await refreshRelationships()
  }, 'initiative/relationships/changeMembership')

const addManager = (userId: string) => { void changeMembership(userId, undefined, true) }
const removeManager = (userId: string) => { void changeMembership(userId, undefined, false) }
const removeMember = (userId: string) => { void changeMembership(userId, false, undefined) }
</script>

<template>
  <PVDataTable
    :value="nonEmptyRelationships"
    class="align-self-stretch"
  >
    <PVColumn
      :header="tt('User ID')"
      field="userId"
      sortable
    >
      <template #body="slotProps">
        <div class="flex flex-column gap-1">
          {{ slotProps.data.userId }}
          <LinkButton
            :label="tt('View User')"
            class="p-button-xs p-button-outlined"
            :to="localePath(`/user/${slotProps.data.userId}`)"
            icon="pi pi-external-link"
            icon-pos="right"
          />
        </div>
      </template>
    </PVColumn>
    <PVColumn
      :header="tt('Updated At')"
      field="updatedAt"
      sortable
    >
      <template #body="slotProps">
        {{ humanReadableTimeFromStandardString(slotProps.data.updatedAt).value }}
      </template>
    </PVColumn>
    <PVColumn
      :header="tt('Role')"
      field="manager"
      sortable
    >
      <template #body="slotProps">
        <span v-if="slotProps.data.manager">{{ tt('Manager') }}</span>
        <span v-else>{{ tt('Member') }}</span>
      </template>
    </PVColumn>
    <PVColumn
      v-if="canManage"
      :header="tt('Actions')"
    >
      <template #body="slotProps">
        <div class="flex flex-column gap-1">
          <PVButton
            v-if="slotProps.data.manager"
            :label="tt('Remove Manager')"
            class="p-button-xs p-button-success p-button-outlined"
            icon="pi pi-user-minus"
            @click="removeManager(slotProps.data.userId)"
          />
          <template v-else>
            <PVButton
              :label="tt('Make Manager')"
              class="p-button-xs p-button-success p-button-outlined"
              icon="pi pi-user-plus"
              @click="addManager(slotProps.data.userId)"
            />
            <PVButton
              :label="tt('Remove From Initiative')"
              class="p-button-xs p-button-danger p-button-outlined"
              icon="pi pi-trash"
              @click="removeMember(slotProps.data.userId)"
            />
          </template>
        </div>
      </template>
    </PVColumn>
  </PVDataTable>
</template>
