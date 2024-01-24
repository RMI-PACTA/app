<script setup lang="ts">
import type OverlayPanel from 'primevue/overlaypanel'
import { useToast } from 'primevue/usetoast'
import { type Portfolio, type PortfolioGroup } from '@/openapi/generated/pacta'
import { selectedCountSuffix } from '@/lib/selection'

const { t } = useI18n()
const toast = useToast()
const { newPortfolioGroup: { newPortfolioGroupVisible } } = useModal()
const pactaClient = usePACTA()

interface Props {
  portfolioGroups: PortfolioGroup[]
  selectedPortfolios: Portfolio[]
  btnClass?: string
}
const props = defineProps<Props>()
interface Emits {
  (e: 'changed-memberships'): void
  (e: 'changed-groups'): void
}
const emit = defineEmits<Emits>()

const prefix = 'components/portfolio/group/membership/MenuButton'
const tt = (s: string) => t(`${prefix}.${s}`)
const statePrefix = `${prefix}[${useStateIDGenerator().id()}]`
const visible = useState<boolean>(`${statePrefix}.visible`, () => false)

const overlayPanel = useState<OverlayPanel>(`${statePrefix}.overlayPanel`)
const toggleMenu = (event: Event) => {
  presentOrFileBug(overlayPanel.value).toggle(event)
  visible.value = !visible.value
}

const changeMemberships = (portfolioGroupId: string, add: boolean) => {
  return async (event: Event) => {
    const portfolioIds = props.selectedPortfolios.map((portfolio) => portfolio.id)
    if (add) {
      await Promise.all(portfolioIds.map((portfolioId) => pactaClient.createPortfolioGroupMembership({ portfolioId, portfolioGroupId })))
    } else {
      await Promise.all(portfolioIds.map((portfolioId) => pactaClient.deletePortfolioGroupMembership({ portfolioId, portfolioGroupId })))
    }
    const pg = presentOrFileBug(props.portfolioGroups.find((pg) => pg.id === portfolioGroupId))
    emit('changed-memberships')
    let summary = ''
    if (add) {
      summary = portfolioIds.length > 1 ? tt('Added OK Plural') : tt('Added OK Singular')
    } else {
      summary = portfolioIds.length > 1 ? tt('Removed OK Plural') : tt('Removed OK Singular')
    }
    summary += ` "${pg.name}"`
    toast.add({
      severity: add ? 'success' : 'warn',
      summary,
      life: 8000,
      detail: `(${props.selectedPortfolios.length}) ${tt('Portfolios')}: ${props.selectedPortfolios.map((p) => p.name).join(', ')}`,
    })
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
  // Created is an ISO date time string. This sorts by newest first, without having
  // to parse the date.
  result.sort((a, b) => a.created < b.created ? 1 : -1)
  return result
})
const classes = computed(() => `p-button-sm ${props.btnClass ?? ''} ${visible.value ? '' : 'p-button-outlined'}`)
</script>

<template>
  <PVButton
    :disabled="!props.selectedPortfolios || props.selectedPortfolios.length === 0"
    :class="classes"
    :label="tt('Group Memberships') + selectedCountSuffix(props.selectedPortfolios)"
    icon="pi pi-table"
    @click="toggleMenu"
  />
  <Teleport to=".modal-group">
    <PVOverlayPanel
      ref="overlayPanel"
      :pt="{ content: { class: 'p-0' } }"
      @hide="() => { visible = false }"
      @show="() => { visible = true }"
    >
      <div class="flex flex-column align-items-stretch">
        <div class="font-bold text-xl p-3 border-bottom-1 border-600 flex gap-2 align-items-center">
          <span>{{ tt('Group Memberships') }}</span>
          <PVButton
            icon="pi pi-times"
            class="p-button-text px-1 py-0 w-auto h-auto p-button-secondary"
            @click="toggleMenu"
          />
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
                class="pseudo-checkbox flex-0 border-2 border-round flex justify-content-center align-items-center"
                :class="option.icon === 'empty' ? 'bg-white' : 'bg-primary-500 text-white border-primary-500'"
              >
                <i
                  v-if="option.icon === 'full'"
                  class="pi pi-check text-base"
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
          class="p-button-text align-self-center p-button-secondary"
          @click="() => newPortfolioGroupVisible = true"
        />
      </div>
    </PVOverlayPanel>
    <PortfolioGroupNewModal
      @created="changedGroups"
    />
  </Teleport>
</template>

<style scoped lang="scss">
// Note, these styles are meant to match the size of the PV checkboxes, which we cannot
// use directly because of the mixed (dash) state here. They appear to be 1.25rem wide.
.pseudo-checkbox {
  width: 1.25rem;
  height: 1.25rem;
}
</style>
