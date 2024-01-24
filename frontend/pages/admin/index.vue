<script setup lang="ts">
const localePath = useLocalePath()

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
    href: '/admin/pacta-version',
  },
  {
    title: 'Initiatives',
    icon: 'pi pi-map',
    desc: 'Create, update + manage Initiatives',
    href: '/admin/initiative',
  },
  {
    title: 'User List',
    icon: 'pi pi-users',
    desc: 'View a list of users, and search for users based on name or email, then edit or delete them',
    href: '/admin/users',
  },
  {
    title: 'User Merge',
    icon: 'pi pi-user-minus',
    desc: 'Merge two user accounts into one, retaining all data',
    href: '/admin/merge',
  },
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
            <div class="flex align-items-center justify-content-between">
              <div class="flex align-items-center">
                <i
                  :class="item.icon"
                  class="text-2xl mr-3"
                />
                <span>{{ item.title }}</span>
              </div>
              <LinkButton
                class="p-button-sm"
                :to="localePath(item.href)"
                icon="pi pi-arrow-right"
                icon-pos="right"
                label="Go"
              />
            </div>
          </template>
          <template #content>
            {{ item.desc }}
          </template>
        </PVCard>
      </div>
    </div>
  </StandardContent>
</template>
