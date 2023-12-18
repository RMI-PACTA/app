<script setup lang="ts">
import type OverlayPanel from 'primevue/overlaypanel'
import { useToast } from 'primevue/usetoast'
import { type Portfolio, type PortfolioGroup } from '@/openapi/generated/pacta'

const { t } = useI18n()
const toast = useToast()
const { newPortfolioGroup: { newPortfolioGroupVisible } } = useModal()
const pactaClient = usePACTA()

interface Props {
  portfolioGroups: PortfolioGroup[]
  selectedPortfolios: Portfolio[]
}
const props = defineProps<Props>()
interface Emits {
  (e: 'changed-memberships'): void
  (e: 'changed-groups'): void
}
const emit = defineEmits<Emits>()

const prefix = 'components/portfolio/group/membership/MenuButton'
const tt = (s: string) => t(`${prefix}.${s}`)
const visible = useState<boolean>(`${prefix}.visible`, () => false)

const overlayPanel = useState<OverlayPanel>(`${prefix}.overlayPanel`)
const toggleMenu = (event: Event) => {
  presentOrFileBug(overlayPanel.value).toggle(event)
  visible.value = !visible.value
}
const hideMenu = () => {
  visible.value = false
}
const showMenu = () => {
  visible.value = true
}

const changeMemberships = (portfolioGroupId: string, add: boolean) => {
  return async (event: Event) => {
    const portfolioIds = props.selectedPortfolios.map((portfolio) => portfolio.id)
    if (add) {
      await Promise.all(portfolioIds.map((portfolioId) => pactaClient.createPortfolioGroupMembership({ portfolioId, portfolioGroupId })))
    } else {
      await Promise.all(portfolioIds.map((portfolioId) => pactaClient.deletePortfolioGroupMembership({ portfolioId, portfolioGroupId })))
    }
    emit('changed-memberships')
    toast.add({
      severity: 'success',
      group: add ? 'portfolio-group-membership-added' : 'portfolio-group-membership-removed',
    })
    toggleMenu(event)
  }
}
const changedGroups = () => {
  emit('changed-groups')
}

type Icon = 'empty' | 'partial' | 'full'

const groupOptions = computed(() => {
  const selected = props.selectedPortfolios
  const result = props.portfolioGroups.map((pg) => {
    const isMember = selected.map((portfolio) => (portfolio.groups ?? []).some((pg2) => pg.id === pg2.portfolioGroup.id))
    const anySelected = isMember.some(m => m)
    const allSelected = isMember.every(m => m)
    let icon: Icon = 'empty'
    let addIfClick = true
    let hoverText = tt('Add all portfolios to group')
    if (allSelected) {
      icon = 'full'
      addIfClick = false
      hoverText = tt('Remove all portfolios from group')
    } else if (anySelected) {
      // TODO(grady) make a pi-square-minus
      icon = 'partial'
      addIfClick = true
      hoverText = tt('Add unselected portfolios to group')
    }
    return {
      id: pg.id,
      label: pg.name,
      icon,
      cmd: changeMemberships(pg.id, addIfClick),
      hoverText,
      created: pg.createdAt,
    }
  })
  result.sort((a, b) => a.created < b.created ? 1 : -1)
  return result
})
</script>

<template>
  <div>
    <PVButton
      v-if="props.selectedPortfolios && props.selectedPortfolios.length > 0"
      :class="visible ? '' : 'p-button-outlined'"
      :label="tt('Group Memberships')"
      @click="toggleMenu"
    />
    <PVOverlayPanel
      ref="overlayPanel"
      :pt="{ content: { class: 'p-0' } }"
      @hide="hideMenu"
      @show="showMenu"
    >
      <div class="flex flex-column align-items-stretch">
        <div class="font-bold text-xl p-3 surface-200">
          {{ tt('Group Memberships') }}
        </div>
        <div
          v-for="option in groupOptions"
          :key="option.id"
          class="border-bottom-1 border-400"
        >
          <PVButton
            v-tooltip="option.hoverText"
            class="text-left p-button-text w-full"
            @click="option.cmd"
          >
            <div class="flex justify-content-start align-items-center gap-3">
              <div
                class="flex-0 border-2 border-round flex justify-content-center align-items-center"
                :class="option.icon === 'empty' ? 'bg-white' : 'bg-primary-500 text-white border-primary-500'"
                style="height: 1.25rem; width: 1.25rem;"
              >
                <i
                  v-if="option.icon === 'full'"
                  class="pi pi-check"
                  style="font-size: 1rem"
                />
                <i
                  v-if="option.icon === 'partial'"
                  class="pi pi-minus"
                  style="font-size: .8rem"
                />
              </div>
              <div class="flex-1 w-full">
                {{ option.label }}
              </div>
            </div>
          </PVButton>
        </div>
        <PVButton
          :label="tt('Create New Group')"
          icon="pi pi-plus"
          @click="() => newPortfolioGroupVisible = true"
        />
      </div>
    </PVOverlayPanel>

    <PVToast
      position="bottom-right"
      group="portfolio-group-membership"
      :dismissable="true"
    >
      <template #message>
        <div class="flex flex-column gap-2">
          <div>{{ tt('Added OK') }} </div>
        </div>
      </template>
    </PVToast>

    <PortfolioGroupNewModal
      @created="changedGroups"
    />
  </div>
</template>
