<script setup lang="ts">
import { X } from 'lucide'
import { ADMIN_COPY, ADMIN_ROUTES } from '~/constants/admin'
import type { AdminDeliveryZone } from '~/types/admin'
import { formatProductPrice, koboToNaira, nairaToKobo } from '~/utils/product'

definePageMeta({ layout: 'admin', middleware: 'admin' })

useHead({ title: `${ADMIN_COPY.deliveryHeading} · ${ADMIN_COPY.brand}` })

interface DeliveryZoneDraft {
  name: string
  feeNaira: number
  active: boolean
}

type DeliveryZoneEditorMode = 'create' | 'edit'

const { createDeliveryZone, fetchDeliveryZones, updateDeliveryZone } = useAdminDeliveryZones()
const { session } = useAdminSession()
const router = useRouter()
const { notifySuccess } = useAppNotifications()
const { withLoading } = useAppLoading()

const zones = ref<readonly AdminDeliveryZone[]>([])
const drafts = ref<Record<string, DeliveryZoneDraft>>({})
const isLoading = ref(true)
const hasFailed = ref(false)
const saveError = ref('')
const savingZoneName = ref<string | null>(null)
const deliveryEditorMode = ref<DeliveryZoneEditorMode | null>(null)
const selectedZone = ref<AdminDeliveryZone | null>(null)
const selectedDraft = ref<DeliveryZoneDraft | null>(null)
const isEditorSaving = ref(false)
const editorErrorMessage = ref('')
const editorDrawer = ref<HTMLElement | null>(null)
const editorCloseButton = ref<HTMLButtonElement | null>(null)

const sortedDrafts = computed(() =>
  [...zones.value]
    .map((zone) => drafts.value[zone.name])
    .filter(Boolean)
    .sort((firstZone, secondZone) => {
      if (firstZone.active !== secondZone.active) {
        return firstZone.active ? -1 : 1
      }

      if (firstZone.feeNaira !== secondZone.feeNaira) {
        return firstZone.feeNaira - secondZone.feeNaira
      }

      return firstZone.name.localeCompare(secondZone.name)
    }),
)

const isEditorOpen = computed(() => deliveryEditorMode.value !== null)

const selectedZoneTitle = computed(() =>
  deliveryEditorMode.value === 'create'
    ? ADMIN_COPY.deliveryNewHeading
    : selectedZone.value?.name ?? ADMIN_COPY.deliveryHeading,
)

const canSubmitSelectedZone = computed(() =>
  selectedDraft.value !== null &&
  selectedDraft.value.name.trim().length > 0 &&
  selectedDraft.value.feeNaira >= 0 &&
  !isEditorSaving.value,
)

const isExpiredSession = async (requestError: unknown): Promise<boolean> => {
  if ((requestError as { status?: number } | null)?.status !== 401) {
    return false
  }

  session.value = null
  await router.push(ADMIN_ROUTES.login)

  return true
}

const toDrafts = (deliveryZones: readonly AdminDeliveryZone[]): Record<string, DeliveryZoneDraft> =>
  Object.fromEntries(
    deliveryZones.map((zone) => [
      zone.name,
      { name: zone.name, feeNaira: koboToNaira(zone.feeKobo), active: zone.active },
    ]),
  )

const load = async (): Promise<void> => {
  isLoading.value = true
  hasFailed.value = false
  saveError.value = ''

  try {
    const response = await withLoading(fetchDeliveryZones, ADMIN_COPY.deliveryLoading)
    zones.value = response.zones
    drafts.value = toDrafts(response.zones)
  } catch (loadError) {
    if (await isExpiredSession(loadError)) {
      return
    }

    hasFailed.value = true
    zones.value = []
    drafts.value = {}
  } finally {
    isLoading.value = false
  }
}

const syncZone = (updatedZone: AdminDeliveryZone): void => {
  const zoneExists = zones.value.some((zone) => zone.name === updatedZone.name)
  zones.value = zoneExists
    ? zones.value.map((zone) => (zone.name === updatedZone.name ? updatedZone : zone))
    : [...zones.value, updatedZone]

  drafts.value = {
    ...drafts.value,
    [updatedZone.name]: {
      name: updatedZone.name,
      feeNaira: koboToNaira(updatedZone.feeKobo),
      active: updatedZone.active,
    },
  }

  if (selectedZone.value?.name === updatedZone.name) {
    selectedZone.value = updatedZone
    selectedDraft.value = toDrafts([updatedZone])[updatedZone.name] ?? null
  }
}

const formatDraftFee = (draft: DeliveryZoneDraft): string =>
  formatProductPrice(nairaToKobo(draft.feeNaira))

const zoneStatusLabel = (zone: AdminDeliveryZone | DeliveryZoneDraft): string =>
  zone.active ? ADMIN_COPY.deliveryVisibleStatus : ADMIN_COPY.deliveryHiddenStatus

const zoneStatusTone = (zone: AdminDeliveryZone | DeliveryZoneDraft): string =>
  zone.active ? 'paid' : 'cancelled'

const createEmptyDraft = (): DeliveryZoneDraft => ({
  name: '',
  feeNaira: 0,
  active: true,
})

const openNewZoneEditor = (): void => {
  deliveryEditorMode.value = 'create'
  selectedZone.value = null
  selectedDraft.value = createEmptyDraft()
  editorErrorMessage.value = ''
}

const openZoneEditor = (zoneName: string): void => {
  const zone = zones.value.find((deliveryZone) => deliveryZone.name === zoneName) ?? null
  deliveryEditorMode.value = 'edit'
  selectedZone.value = zone
  selectedDraft.value = zone ? toDrafts([zone])[zone.name] ?? null : null
  editorErrorMessage.value = zone ? '' : ADMIN_COPY.deliveryLoadFailed
}

const closeZoneEditor = (): void => {
  deliveryEditorMode.value = null
  selectedZone.value = null
  selectedDraft.value = null
  editorErrorMessage.value = ''
}

const persistSelectedZone = async (
  draft: DeliveryZoneDraft,
  activeOverride?: boolean,
): Promise<AdminDeliveryZone | null> => {
  const isCreating = deliveryEditorMode.value === 'create'
  const nextActive = activeOverride ?? draft.active
  isEditorSaving.value = true
  editorErrorMessage.value = ''
  saveError.value = ''

  try {
    const savedZone = await withLoading(
      () => isCreating
        ? createDeliveryZone({
            name: draft.name.trim(),
            feeKobo: nairaToKobo(draft.feeNaira),
            active: nextActive,
          })
        : updateDeliveryZone(draft.name, {
            feeKobo: nairaToKobo(draft.feeNaira),
            active: nextActive,
          }),
      isCreating ? ADMIN_COPY.deliveryCreating : ADMIN_COPY.deliverySaving,
    )

    syncZone(savedZone)
    closeZoneEditor()
    notifySuccess(
      isCreating ? ADMIN_COPY.deliveryCreatedNotice : ADMIN_COPY.deliverySavedNotice,
      savedZone.name,
    )

    return savedZone
  } catch (persistError) {
    if (!(await isExpiredSession(persistError))) {
      editorErrorMessage.value = isCreating
        ? ADMIN_COPY.deliveryCreateFailed
        : ADMIN_COPY.deliverySaveFailed
    }

    return null
  } finally {
    isEditorSaving.value = false
  }
}

const saveSelectedZone = async (): Promise<void> => {
  if (!selectedDraft.value || !canSubmitSelectedZone.value) {
    return
  }

  await persistSelectedZone(selectedDraft.value)
}

const toggleSelectedZoneVisibility = async (): Promise<void> => {
  if (!selectedDraft.value || deliveryEditorMode.value !== 'edit') {
    return
  }

  savingZoneName.value = selectedDraft.value.name
  await persistSelectedZone(selectedDraft.value, !selectedDraft.value.active)
  savingZoneName.value = null
}

// Nested controls keep their native behavior; the rest of the row opens edit.
const isInteractiveRowTarget = (eventTarget: EventTarget | null): boolean =>
  eventTarget instanceof Element &&
  Boolean(eventTarget.closest('a, button, input, select, textarea, summary'))

const handleZoneRowClick = (mouseEvent: MouseEvent, zoneName: string): void => {
  if (!isInteractiveRowTarget(mouseEvent.target)) {
    openZoneEditor(zoneName)
  }
}

const handleZoneRowKeydown = (keyboardEvent: KeyboardEvent, zoneName: string): void => {
  if (
    isInteractiveRowTarget(keyboardEvent.target) ||
    !['Enter', ' '].includes(keyboardEvent.key)
  ) {
    return
  }

  keyboardEvent.preventDefault()
  openZoneEditor(zoneName)
}

useModalFocus({
  isOpen: isEditorOpen,
  container: editorDrawer,
  initialFocus: editorCloseButton,
  onEscape: closeZoneEditor,
  lockScroll: false,
})

watch(isEditorOpen, (isOpen) => {
  if (!import.meta.client) {
    return
  }

  document.body.classList.toggle('has-admin-drawer', isOpen)
})

onMounted(() => {
  void load()
})

onBeforeUnmount(() => {
  document.body.classList.remove('has-admin-drawer')
})
</script>

<template>
  <div class="admin-page">
    <header class="admin-page-head">
      <div>
        <h1>{{ ADMIN_COPY.deliveryHeading }}</h1>
        <p class="admin-page-subheading">{{ ADMIN_COPY.deliverySubheading }}</p>
      </div>

      <button type="button" class="admin-inline-action" @click="openNewZoneEditor">
        {{ ADMIN_COPY.deliveryCreate }}
      </button>
    </header>

    <p v-if="saveError" class="admin-detail-error" role="alert">
      {{ saveError }}
    </p>

    <p v-if="isLoading" class="admin-orders-message">{{ ADMIN_COPY.deliveryLoading }}</p>

    <div v-else-if="hasFailed" class="admin-orders-message" role="alert">
      <p>{{ ADMIN_COPY.deliveryLoadFailed }}</p>
      <button type="button" class="admin-inline-action" @click="load">
        {{ ADMIN_COPY.retry }}
      </button>
    </div>

    <p v-else-if="sortedDrafts.length === 0" class="admin-orders-message">
      {{ ADMIN_COPY.deliveryEmpty }}
    </p>

    <div v-else class="admin-table-scroll">
      <table class="admin-table admin-delivery-table">
        <thead>
          <tr>
            <th scope="col">{{ ADMIN_COPY.deliveryZoneLabel }}</th>
            <th scope="col">{{ ADMIN_COPY.deliveryFeeLabel }}</th>
            <th scope="col">Status</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="draft in sortedDrafts"
            :key="draft.name"
            class="admin-table-row"
            tabindex="0"
            role="button"
            :aria-label="`Edit ${draft.name}`"
            @click="handleZoneRowClick($event, draft.name)"
            @keydown="handleZoneRowKeydown($event, draft.name)"
          >
            <td :data-label="ADMIN_COPY.deliveryZoneLabel">
              <strong>{{ draft.name }}</strong>
            </td>
            <td :data-label="ADMIN_COPY.deliveryFeeLabel" class="admin-table-number">
              {{ formatDraftFee(draft) }}
            </td>
            <td data-label="Status">
              <span class="admin-status" :data-status="zoneStatusTone(draft)">
                {{ zoneStatusLabel(draft) }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <Teleport to="body">
      <div
        v-if="isEditorOpen"
        class="admin-drawer-backdrop"
        role="presentation"
        @click.self="closeZoneEditor"
      >
        <aside
          ref="editorDrawer"
          class="admin-drawer"
          role="dialog"
          aria-modal="true"
          aria-labelledby="admin-delivery-editor-title"
          tabindex="-1"
        >
          <header class="admin-drawer-head">
            <div>
              <p class="admin-detail-eyebrow">{{ ADMIN_COPY.deliveryHeading }}</p>
              <h2 id="admin-delivery-editor-title">{{ selectedZoneTitle }}</h2>
            </div>

            <div class="admin-drawer-actions">
              <button
                v-if="deliveryEditorMode === 'edit' && selectedDraft"
                type="button"
                class="admin-inline-action"
                :disabled="isEditorSaving || savingZoneName === selectedDraft.name"
                @click="toggleSelectedZoneVisibility"
              >
                {{ selectedDraft.active ? ADMIN_COPY.deliveryHide : ADMIN_COPY.deliveryShow }}
              </button>
              <button
                ref="editorCloseButton"
                type="button"
                class="admin-drawer-close"
                :aria-label="ADMIN_COPY.closeDeliveryEditor"
                @click="closeZoneEditor"
              >
                <StorefrontIcon :icon-node="X" :size="18" :stroke-width="1.75" />
              </button>
            </div>
          </header>

          <form
            v-if="selectedDraft"
            class="admin-product-form"
            @submit.prevent="saveSelectedZone"
          >
            <div class="admin-detail-grid">
              <label class="admin-field">
                <span>{{ ADMIN_COPY.deliveryNewNameLabel }}</span>
                <input
                  v-model.trim="selectedDraft.name"
                  type="text"
                  required
                  autocomplete="off"
                  :readonly="deliveryEditorMode === 'edit'"
                >
              </label>

              <label class="admin-field">
                <span>{{ ADMIN_COPY.deliveryNewFeeLabel }}</span>
                <input
                  v-model.number="selectedDraft.feeNaira"
                  type="number"
                  min="0"
                  step="1"
                  required
                >
              </label>

              <div class="admin-detail-panel">
                <h2>Checkout visibility</h2>
                <p class="admin-detail-muted">
                  Hidden areas remain in admin but stop appearing on checkout and bespoke confirmation.
                </p>
                <label class="admin-delivery-toggle">
                  <input v-model="selectedDraft.active" type="checkbox">
                  <span>{{ ADMIN_COPY.deliveryActiveLabel }}</span>
                </label>
              </div>

              <div class="admin-detail-panel">
                <h2>Preview</h2>
                <p class="admin-detail-value">{{ formatDraftFee(selectedDraft) }}</p>
                <span class="admin-status" :data-status="zoneStatusTone(selectedDraft)">
                  {{ zoneStatusLabel(selectedDraft) }}
                </span>
              </div>
            </div>

            <p v-if="editorErrorMessage" class="admin-detail-error" role="alert">
              {{ editorErrorMessage }}
            </p>

            <div class="admin-product-form-actions">
              <button
                type="submit"
                class="admin-inline-action"
                :disabled="!canSubmitSelectedZone"
              >
                {{ isEditorSaving ? ADMIN_COPY.deliverySaving : ADMIN_COPY.deliverySave }}
              </button>
            </div>
          </form>
        </aside>
      </div>
    </Teleport>
  </div>
</template>
