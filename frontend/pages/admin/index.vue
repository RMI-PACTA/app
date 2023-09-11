<script setup lang="ts">
interface AdminItem {
  title: string
  icon: string
  desc: string
  href: string
}

const adminItems: AdminItem[] = [
  {
    title: 'PACTA Versions',
    icon: 'pi pi-check-square',
    desc: 'Create, update + manage versions of the PACTA algorithm + docker image',
    href: '/admin/pacta-version'
  }
]
</script>

<template>
  <StandardContent>
    <TitleBar title="Admin Console" />
    <p>
      This is the home for administrators of the PACTA platform.
    </p>
    <div class="grid align-items-stretch pb-3">
      <div
        v-for="item in adminItems"
        :key="item.title"
        class="col-12 md:col-6"
      >
        <PVCard class="h-full">
          <template #title>
            <div class="flex align-items-center">
              <i
                :class="item.icon"
                class="text-2xl mr-3"
              />
              <span>{{ item.title }}</span>
            </div>
          </template>
          <template #content>
            {{ item.desc }}
          </template>
          <template #footer>
            <div class="flex justify-content-end w-full">
              <LinkButton
                class="p-button-sm"
                :to="item.href"
                icon="pi pi-arrow-right"
                icon-pos="right"
                label="Go"
              />
            </div>
          </template>
        </PVCard>
      </div>
    </div>
  </standardcontent>
</template>
