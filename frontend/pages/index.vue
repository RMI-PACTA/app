<script setup lang="ts">
const isAuthenticated = useIsAuthenticated()
const { signIn } = useSignIn()
const { t } = useI18n()
const localePath = useLocalePath()

const tt = (s: string) => t(`pages/index.${s}`)

</script>

<template>
  <StandardContent>
    <div class="border-round p-3">
      <h1 class="text-5xl mt-0 mb-2 ">
        {{ tt('Heading') }}
      </h1>
      <div class="text-2xl mt-0 mb-4">
        {{ tt('Subheading') }}
      </div>
      <p class="">
        {{ tt('Description') }}
      </p>
    </div>
    <div class="flex gap-3 flex-wrap">
      <div class="border-primary border-2 border-round p-3 flex flex-column gap-3 flex-1">
        <h2 class="m-0">
          {{ tt('Section1Heading') }}
        </h2>
        <p>{{ tt('Section1Paragraph1') }}</p>
        <p>{{ tt('Section1Paragraph2') }}</p>
      </div>
      <div
        class="w-12 md:w-6 xl:w-5 flex p-2 border-2 border-round align-items-center justify-content-center"
      >
        <PVImage
          preview
          class="img-max-w-full"
          src="/img/how-it-works.jpg"
        />
      </div>
    </div>
    <div class="flex gap-3 flex-column p-3">
      <h2 class="m-0">
        {{ tt('Section2Heading') }}
      </h2>
      <p>{{ tt('Section2Paragraph1') }}</p>
      <p>
        {{ tt('Section2Paragraph2') }}
        <LinkButton
          :to="localePath('/input-guide')"
          class="p-button-xs align-self-center py-0 ml-2 p-button-outlined"
          :label="tt('Input User Guide')"
          icon="pi pi-arrow-right"
          icon-pos="right"
          new-tab
        />
      </p>
      <p>{{ tt('Section2Paragraph3') }}</p>
      <div class="flex justify-content-center mt-2">
        <PVButton
          v-if="!isAuthenticated"
          :label="tt('Sign In')"
          icon="pi pi-sign-in"
          icon-pos="right"
          class="max-w-30rem"
          @click="signIn"
        />
        <LinkButton
          v-else
          :label="tt('Get Started')"
          icon="pi pi-arrow-right"
          icon-pos="right"
          class="max-w-30rem"
          :to="localePath('/upload')"
        />
      </div>
    </div>
    <div class="flex gap-3 flex-column p-3">
      <h2 class="m-0">
        {{ tt('Section4Heading') }}
      </h2>
      <p>
        {{ tt('Section4Paragraph1') }}
        <LinkButton
          to="https://pacta.rmi.org"
          external
          class="p-button-xs align-self-center py-0 ml-2 p-button-outlined"
          :label="tt('Section4CTA')"
          icon="pi pi-arrow-right"
          icon-pos="right"
        />
      </p>
    </div>
  </StandardContent>
</template>

<style lang="scss">
.img-max-w-full {
  max-width: 100%;

  img {
    max-width: 100%;
  }
}

.p-image-toolbar {
  z-index: 1;
  background: var(--primary-color);
}
</style>
