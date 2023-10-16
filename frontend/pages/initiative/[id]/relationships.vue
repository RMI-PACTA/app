<script setup lang="ts">
import { humanReadableTimeFromStandardString } from '@/lib/time'

// const router = useRouter()
const pactaClient = await usePACTA()
// const { loading: { withLoading } } = useModal()
const { fromParams } = useURLParams()
const localePath = useLocalePath()
const { getMaybeMe } = useSession()

const { maybeMe, isAdmin } = await getMaybeMe()

const id = presentOrCheckURL(fromParams('id'))
const prefix = `initiative/${id}/relationships`
const [
  { data: relationships },
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

const makeManager = (userId: string) => {
  console.warn('make manager not yet implemented')
}
const removeFromInitiative = (userId: string) => {
  console.warn('remove from initiative not yet implemented')
}
const removeManager = (userId: string) => {
  console.warn('remove manager not yet implemented')
}
</script>

<template>
  <PVDataTable
    :value="nonEmptyRelationships"
    class="align-self-stretch"
  >
    <PVColumn
      header="User ID"
      field="userId"
      sortable
    >
      <template #body="slotProps">
        <LinkButton
          label="View User"
          class="p-button-xs p-button-outlined"
          :to="localePath(`/user/${slotProps.data.id}`)"
          icon="pi pi-external-link"
          icon-pos="right"
        />
      </template>
    </PVColumn>
    <PVColumn
      header="Updated At"
      field="updatedAt"
      sortable
    >
      <template #body="slotProps">
        {{ humanReadableTimeFromStandardString(slotProps.data.updatedAt) }}
      </template>
    </PVColumn>
    <PVColumn
      header="Manager"
      field="manager"
      sortable
    />
    <PVColumn
      v-if="canManage"
      header="Actions"
    >
      <template #body="slotProps">
        <div class="flex flex-column gap-1">
          <PVButton
            v-if="slotProps.data.manager"
            label="Remove Manager"
            class="p-button-xs p-button-success p-button-outlined"
            icon="pi pi-user-minus"
            @click="removeManager(slotProps.data.userId)"
          />
          <template v-else>
            <PVButton
              label="Make Manager"
              class="p-button-xs p-button-success p-button-outlined"
              icon="pi pi-user-plus"
              @click="makeManager(slotProps.data.userId)"
            />
            <PVButton
              label="Remove From Initiative"
              class="p-button-xs p-button-danger p-button-outlined"
              icon="pi pi-trash"
              @click="removeFromInitiative(slotProps.data.userId)"
            />
          </template>
        </div>
      </template>
    </PVColumn>
  </PVDataTable>
</template>
