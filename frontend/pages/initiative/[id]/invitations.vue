<script setup lang="ts">
import type DataTable from 'primevue/datatable'

const { humanReadableTimeFromStandardString } = useTime()
const { public: { baseURL } } = useRuntimeConfig()
const { t } = useI18n()
const { loading: { withLoading } } = useModal()
const { fromParams } = useURLParams()
const localePath = useLocalePath()
const pactaClient = usePACTA()

const tt = (key: string) => t(`pages/initiative/invitations.${key}`)

const id = presentOrCheckURL(fromParams('id'))
const prefix = `initiative/${id}/invitations`
const createModalVisible = useState<boolean>(`${prefix}.createModalVisible`, () => false)
const rawNewInvitations = useState<string>(`${prefix}.rawNewInvitations`, () => '')
const numToRandomize = useState<number>(`${prefix}.numToRandomize`, () => 5)
const dataTable = useState<DataTable>(`${prefix}.dataTable`)

const { initiative, invitations, refreshInitiative, refreshInvitations } = await useInitiativeData(id)

const newInvitations = computed(() => {
  const raw = rawNewInvitations.value
  if (!raw) return []
  return raw.split(/[\s,]+/).map((line) => line.trim()).filter((line) => line)
})
const initiativeIDStructure = /^[a-zA-Z0-9_-]+$/
const newInvitationsErrorMessage = computed(() => {
  const errors: string[] = []
  const ni = newInvitations.value
  const exists = new Set(invitations.value.map((i) => i.id))
  const novel = new Set()
  for (const id of ni) {
    if (exists.has(id)) {
      errors.push(tt('Already Exists') + ': ' + id)
    }
    if (novel.has(id)) {
      errors.push(tt('Duplicate') + ': ' + id)
    }
    if (!initiativeIDStructure.test(id)) {
      errors.push(tt('Includes Illegal Characters') + ': ' + id)
    }
    novel.add(id)
  }
  return errors.join('\n')
})

const generateRandom = () => {
  const newValues = []
  const n = numToRandomize.value
  for (let i = 0; i < n; i++) {
    const rv = Math.round(Math.random() * 1000000)
    newValues.push(`${initiative.value.id}-${rv}`)
  }
  rawNewInvitations.value = newValues.join('\n')
}

const createInvitations = () => {
  const initiativeId = initiative.value.id
  const requests = newInvitations.value.map((id) => ({
    id,
    initiativeId,
  })).map((request) => pactaClient.createInitiativeInvitation(request))
  return withLoading(
    () => Promise.all(requests)
      .then(refreshInitiative)
      .then(() => { createModalVisible.value = false }),
    `${prefix}.createInvitations`,
  )
}
const deleteInvitation = (id: string) => withLoading(
  () => pactaClient.deleteInitiativeInvitation(id).then(refreshInitiative),
  `${prefix}.deleteInvitation`,
)
const deleteAll = () => withLoading(
  () => Promise.all(invitations.value.filter(i => !i.usedAt).map((i) => deleteInvitation(i.id))).then(refreshInvitations),
  `${prefix}.deleteAll`,
)
const invitationURL = (id: string) => {
  return `${baseURL}${localePath(`/join/${id}`)}`
}
const getData = (e: { data: any }): any => {
  return e.data
}
const doExport = () => {
  dataTable.value.exportCSV()
}
</script>

<template>
  <div class="flex flex-column gap-3">
    <p>
      {{ initiative.requiresInvitationToJoin ? tt('Yes Invitations') : tt('No Invitations') }}
      {{ tt('You Can Change') }}
    </p>
    <template v-if="initiative.requiresInvitationToJoin">
      <PVButton
        icon="pi pi-plus"
        :label="tt('Create Initiative Invitations')"
        class="align-self-start"
        @click="() => createModalVisible = true"
      />
      <PVDataTable
        ref="dataTable"
        :value="invitations"
        :export-filename="`initiative-${initiative.id}-invitations`"
        :export-function="getData"
        size="small"
      >
        <PVColumn
          sortable
          :header="tt('Invitation Code')"
          field="id"
        >
          <template #body="slotProps">
            <div class="flex flex-column gap-1 align-items-start">
              <span>{{ slotProps.data.id }}</span>
              <div class="p-buttonset">
                <CopyToClipboardButton
                  :cta="tt('Copy Code')"
                  class="p-button-outlined p-button-xs"
                  :value="slotProps.data.id"
                />
                <CopyToClipboardButton
                  :cta="tt('Copy Share URL')"
                  class="p-button-outlined p-button-xs"
                  icon="pi pi-link"
                  :value="invitationURL(slotProps.data.id)"
                />
              </div>
            </div>
          </template>
        </PVColumn>
        <PVColumn
          sortable
          :header="tt('Used At')"
          field="usedAt"
        >
          <template #body="slotProps">
            <div
              v-if="slotProps.data.usedAt"
              class="flex flex-column gap-1 align-items-start"
            >
              <span>{{ tt('Used At') }} {{ humanReadableTimeFromStandardString(slotProps.data.usedAt) }}</span>
              <LinkButton
                :to="localePath(`/user/${slotProps.data.usedByUserId}`)"
                :label="tt('User Profile')"
                class="p-button-outlined p-button-xs"
                icon="pi pi-external-link"
                icon-pos="right"
              />
            </div>
            <template v-else>
              <div class="flex flex-column gap-1 align-items-start">
                <span>{{ tt('Unused') }}</span>
                <PVButton
                  icon="pi pi-trash"
                  class="p-button-danger p-button-outlined p-button-xs"
                  :label="tt('Revoke')"
                  @click="() => deleteInvitation(slotProps.data.id)"
                />
              </div>
            </template>
          </template>
        </PVColumn>
      </PVDataTable>
      <div class="flex gap-2">
        <PVButton
          :label="tt('Export')"
          icon="pi pi-download"
          class="p-button-secondary p-button-outlined"
          @click="doExport"
        />
        <PVButton
          :label="tt('Revoke All Unused')"
          icon="pi pi-trash"
          class="p-button-danger p-button-outlined"
          @click="deleteAll"
        />
      </div>
    </template>
    <StandardDebug
      :value="invitations"
      label="Invitations"
    />
    <StandardModal
      v-model:visible="createModalVisible"
      :header="tt('Create Header')"
      :sub-header="tt('Create Subheader')"
    >
      <FormField
        :label="tt('New Invitations')"
        :help-text="tt('New Invitations Help')"
      >
        <div class="flex gap-2 w-full">
          <PVTextarea
            v-model="rawNewInvitations"
            class="flex-1"
            auto-resize
          />
          <PVButton
            class="p-button-outlined"
            @click="generateRandom"
          >
            <div class="flex flex-column gap-2 align-items-center">
              <span>{{ tt('Generate Random') }}</span>
              <PVInputNumber
                v-model="numToRandomize"
                class="w-5rem text-align-center"
                input-class="w-5rem"
                @click="$event.stopPropagation()"
              />
            </div>
          </PVButton>
        </div>
        <p
          v-if="newInvitationsErrorMessage"
          class="text-red-500 mb-0 mt-1"
        >
          {{ newInvitationsErrorMessage }}
        </p>
      </FormField>

      <div class="flex gap-2">
        <PVButton
          :label="tt('Cancel')"
          icon="pi pi-arrow-left"
          class="p-button-secondary p-button-outlined"
          @click="() => createModalVisible = false"
        />
        <PVButton
          :disabled="newInvitations.length === 0 || !!newInvitationsErrorMessage"
          :label="`${tt('Create Initiative Invitations')} (${newInvitations.length})`"
          icon="pi pi-arrow-right"
          icon-pos="right"
          @click="createInvitations"
        />
      </div>
    </StandardModal>
  </div>
</template>
