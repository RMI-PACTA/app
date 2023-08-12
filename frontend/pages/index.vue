<script setup lang="ts">
import { type NewPet } from 'openapi/generated/pacta'

const { isAuthenticated, signIn, signOut } = await useMSAL()

const { pactaClient } = useAPI()

const pets = useState<NewPet[]>('index.pets', () => [])

const testAPI = async (): Promise<void> => {
  const resp = await pactaClient.findPets()
  if ('message' in resp) {
    console.error('error response from server', resp)
    return
  }
  pets.value = resp
}
</script>

<template>
  <StandardContent>
    <TitleBar title="RMI PACTA" />
    <p>
      This will eventually be the site for RMI's PACTA, but now it's mostly just a placeholder.
    </p>
    <p>
      This project is open source. You can view the code at <a href="https://github.com/RMI/pacta">github.com/RMI/pacta</a>.
    </p>

    <PVButton
      v-if="!isAuthenticated"
      label="Sign In"
      @click="signIn"
    />
    <div v-else>
      <PVButton
        class="inline-block"
        label="Test the API"
        @click="testAPI"
      />

      <PVButton
        class="inline-block m-2"
        label="Sign Out"
        icon="pi pi-sign-out"
        @click="signOut"
      />
    </div>
  </StandardContent>
</template>
