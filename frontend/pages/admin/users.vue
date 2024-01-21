<script setup lang="ts">
import { type User } from '@/openapi/generated/pacta'

const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()
const localePath = useLocalePath()

const prefix = 'admin/users'

const textQuery = useState<string>(`${prefix}.textQuery`, () => '')
const canEditQuery = useState<boolean>(`${prefix}.canEditQuery`, () => true)
const users = useState<User[]>(`${prefix}.users`, () => [])
const hasNextPage = useState<boolean>(`${prefix}.hasNextPage`, () => false)
const cursor = useState<string>(`${prefix}.cursor`, () => '')

const executeSearch = () => pactaClient.userQuery({
  wheres: [
    {
      nameOrEmailLike: textQuery.value,
    },
  ],
  cursor: cursor.value,
}).then((resp) => {
  users.value = resp.users
  cursor.value = resp.cursor
  hasNextPage.value = resp.hasNextPage
  canEditQuery.value = false
})

const performSearch = () => withLoading(executeSearch, `${prefix}.executeSearch`)
const more = () => withLoading(executeSearch, `${prefix}.more`)
const moreAll = async () => {
  await withLoading(async () => {
    while (hasNextPage.value) {
      await more()
    }
  }, `${prefix}.moreAll`)
}
const reset = () => {
  textQuery.value = ''
  canEditQuery.value = true
  users.value = []
  cursor.value = ''
  hasNextPage.value = false
}
</script>

<template>
  <StandardContent>
    <TitleBar title="User Search" />
    <p>To list all users (by most recent first) just perform an empty search.</p>
    <div class="p-inputgroup">
      <PVInputText
        v-model="textQuery"
        :disabled="!canEditQuery"
        placeholder="Search By Name or Email"
      />
      <PVButton
        label="Search"
        :disabled="!canEditQuery"
        @click="performSearch"
      />
      <PVButton
        v-if="!canEditQuery"
        class="p-button-secondary p-button-text"
        icon="pi pi-refresh"
        @click="reset"
      />
    </div>
    <PVDataTable
      v-show="!canEditQuery"
      data-key="id"
      class="w-full"
      :value="users"
      empty-message="No Results"
    >
      <PVColumn
        header="ID"
      >
        <template #body="slotProps">
          <div class="flex flex-column gap-1 align-items-start">
            <span>{{ slotProps.data.id }}</span>
            <CopyToClipboardButton
              :value="slotProps.data.id"
              class="p-button-secondary p-button-outlined p-button-xs"
            />
          </div>
        </template>
      </PVColumn>
      <PVColumn
        header="Name"
        field="name"
      />
      <PVColumn
        header="Email"
      >
        <template #body="slotProps">
          <div class="flex flex-column align-items-start">
            <span>Canonical:</span>
            <div class="flex gap-1 align-items-center pb-2">
              <b>{{ slotProps.data.canonicalEmail }}</b>
              <CopyToClipboardButton
                :value="slotProps.data.canonicalEmail"
                class="p-button-secondary p-button-outlined p-button-xs"
              />
            </div>
            <span>Entered:</span>
            <div class="flex gap-1 align-items-center">
              <b>{{ slotProps.data.canonicalEmail }}</b>
              <CopyToClipboardButton
                :value="slotProps.data.canonicalEmail"
                class="p-button-secondary p-button-outlined p-button-xs"
              />
            </div>
          </div>
        </template>
      </PVColumn>
      <PVColumn header="More">
        <template #body="slotProps">
          <div class="flex flex-column gap-1">
            <StandardDebug
              label="Metadata"
              :value="slotProps.data"
              class="w-full"
              always
            />
            <LinkButton
              :to="localePath(`/user/${slotProps.data.id}/edit`)"
              icon="pi pi-arrow-right"
              icon-pos="right"
              class="p-button-xs p-button-outlined"
              label="Edit"
            />
          </div>
        </template>
      </PVColumn>
      <template
        v-if="hasNextPage"
        #footer
      >
        <div class="flex gap-3 flex-wrap">
          <PVButton
            label="More..."
            @click="more"
          />
          <PVButton
            label="Get All (slow)"
            class="p-button-oulined p-button-secondary"
            @click="moreAll"
          />
        </div>
      </template>
    </PVDataTable>
  </StandardContent>
</template>
