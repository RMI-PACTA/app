<script setup lang="ts">
import { useConfirm } from 'primevue/useconfirm'

const pactaClient = usePACTA()
const { loading: { withLoading } } = useModal()
const localePath = useLocalePath()
const { require: confirm } = useConfirm()

const prefix = 'admin/merge'

const fromUserId = useState<string>(`${prefix}.fromUserId`, () => '')
const toUserId = useState<string>(`${prefix}.toUserId`, () => '')
const done = useState<boolean>(`${prefix}.done`, () => false)

const doMerge = async () => {
  await withLoading(() => pactaClient.mergeUsers({
    fromUserId: fromUserId.value,
    toUserId: toUserId.value,
  }), `${prefix}.doMerge`)

  done.value = true
}

const clickMerge = async () => {
  await new Promise((resolve, reject) => {
    confirm({
      header: 'Are you 100% sure?',
      message: 'This will transfer all assets from the source user to the destination user, and then delete the source user. This cannot be undone. Only proceed if you have tripple checked the user IDs and are confident in this procedure.',
      icon: 'pi pi-user-minus',
      position: 'center',
      blockScroll: true,
      reject: () => { /* noop */ },
      rejectLabel: 'Cancel',
      rejectIcon: 'pi pi-times',
      rejectClass: 'p-button-secondary p-button-text',
      acceptClass: 'p-button-danger',
      acceptLabel: 'Proceed',
      accept: () => {
        doMerge().then(resolve).catch(reject)
      },
      acceptIcon: 'pi pi-check',
    })
  })
}

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
          :to="fromUserId ? localePath(`/user/${fromUserId}/edit`) : undefined"
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
          :to="toUserId ? localePath(`/user/${toUserId}/edit`) : undefined"
          new-tab
        />
      </div>
    </div>
    <PVButton
      :disabled="done || !fromUserId || !toUserId || fromUserId === toUserId"
      label="Perform Merge"
      class="p-button-danger"
      icon="pi pi-user-minus"
      @click="clickMerge"
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
