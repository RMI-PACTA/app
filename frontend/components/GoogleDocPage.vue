<script setup lang="ts">
const router = useRouter()
interface Props {
  src: string
}
const props = defineProps<Props>()
</script>

<template>
  <StandardFullWidth>
    <div
      class="overflow-hidden w-screen h-screen flex justify-content-center"
    >
      <div class="relative">
        <div class="absolute top-0 left-0 flex align-items-center py-2 pr-4 pl-3 gap-1 bg-primary">
          <PVButton
            icon="pi pi-arrow-left"
            class="p-button-text text-white py-2 w-auto px-2"
            @click="() => router.go(-1)"
          />
          <LinkButton
            to="/"
            class="p-button-text text-2xl text-white py-0 px-0"
          >
            PACTA Documents
          </LinkButton>
        </div>
        <iframe
          :src="props.src"
          frameBorder="0"
          style="width: calc(100vw + 144pt); height: calc(100vh); margin: 0 -72pt; max-width: calc(612pt + 140pt);"
        />
      </div>
    </div>
  </StandardFullWidth>
</template>
