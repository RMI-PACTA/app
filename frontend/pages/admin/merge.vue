<script setup lang="ts">
const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()
const localePath = useLocalePath()

const prefix = 'admin/merge'

const fromUserId = useState<string>(`${prefix}.fromUserId`, () => '')
const toUserId = useState<string>(`${prefix}.toUserId`, () => '')
const done = useState<boolean>(`${prefix}.done`, () => false)

const doMerge = () => withLoading(() => pactaClient.mergeUsers({
  fromUserId: fromUserId.value,
  toUserId: toUserId.value,
}).then((resp) => {
  done.value = true
}), `${prefix}.doMerge`)

const reset = () => {
  fromUserId.value = ''
  toUserId.value = ''
  done.value = false
}
</script>

<template>
  <StandardContent>
    <TitleBar title="User Merge" />
    <p>The first user entered (Source) will have all of their assets transferred to the second user (Destination), and then their account will be deleted. Use with extreme caution!</p>
    <div class="flex gap-2">
      <span class="font-bold text-lg">
        Source User ID:
      </span>
      <div class="p-inputgroup">
        <PVInputText
          v-model="fromUserId"
          :disabled="done"
          placeholder="Source UserID"
        />
        <LinkButton
          class="p-button-secondary p-button-text"
          icon="pi pi-external-link"
          :to="localePath(`/user/${fromUserId}/edit`)"
          new-tab
        />
      </div>
    </div>
    <div class="flex gap-2">
      <span class="font-bold text-lg">
        Destination User ID:
      </span>
      <div class="p-inputgroup">
        <PVInputText
          v-model="toUserId"
          :disabled="done"
          placeholder="Destination UserID"
        />
        <LinkButton
          class="p-button-secondary p-button-text"
          icon="pi pi-external-link"
          :to="localePath(`/user/${toUserId}/edit`)"
          new-tab
        />
      </div>
    </div>
    <PVButton
      :disabled="done || !fromUserId || !toUserId || fromUserId === toUserId"
      label="Perform Merge"
      class="p-button-danger"
      icon="pi pi-user-minus"
      @click="doMerge"
    />
    <PVMessage
      v-if="done"
      severity="success"
    >
      The merge has been completed.
    </PVMessage>
    <PVButton
      v-if="done"
      label="Reset"
      class="p-button-secondary"
      icon="pi pi-refresh"
      @click="reset"
    />
  </StandardContent>
</template>
