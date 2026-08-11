<script setup lang="ts">
import { ADMIN_COPY, ADMIN_ROUTES } from '~/constants/admin'
import type { AdminDeliveryZone } from '~/types/admin'
import { formatProductPrice, koboToNaira, nairaToKobo } from '~/utils/product'

definePageMeta({ layout: 'admin', middleware: 'admin' })

useHead({ title: `${ADMIN_COPY.deliveryHeading} · ${ADMIN_COPY.brand}` })

interface DeliveryZoneDraft {
  readonly name: string
  feeNaira: number
  active: boolean
}

const { fetchDeliveryZones, updateDeliveryZone } = useAdminDeliveryZones()
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

const sortedDrafts = computed(() =>
  zones.value.map((zone) => drafts.value[zone.name]).filter(Boolean),
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
  zones.value = zones.value.map((zone) =>
    zone.name === updatedZone.name ? updatedZone : zone,
  )
  drafts.value = {
    ...drafts.value,
    [updatedZone.name]: {
      name: updatedZone.name,
      feeNaira: koboToNaira(updatedZone.feeKobo),
      active: updatedZone.active,
    },
  }
}

const formatDraftFee = (draft: DeliveryZoneDraft): string =>
  formatProductPrice(nairaToKobo(draft.feeNaira))

const saveZone = async (draft: DeliveryZoneDraft): Promise<void> => {
  savingZoneName.value = draft.name
  saveError.value = ''

  try {
    const updatedZone = await withLoading(
      () => updateDeliveryZone(draft.name, {
        feeKobo: nairaToKobo(draft.feeNaira),
        active: draft.active,
      }),
      ADMIN_COPY.deliverySaving,
    )

    syncZone(updatedZone)
    notifySuccess(ADMIN_COPY.deliverySavedNotice, updatedZone.name)
  } catch (updateError) {
    if (!(await isExpiredSession(updateError))) {
      saveError.value = ADMIN_COPY.deliverySaveFailed
    }
  } finally {
    savingZoneName.value = null
  }
}

onMounted(() => {
  void load()
})
</script>

<template>
  <div class="admin-page">
    <header class="admin-page-head">
      <div>
        <h1>{{ ADMIN_COPY.deliveryHeading }}</h1>
        <p class="admin-page-subheading">{{ ADMIN_COPY.deliverySubheading }}</p>
      </div>
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
            <th scope="col">Preview</th>
            <th scope="col">Status</th>
            <th scope="col">Action</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="draft in sortedDrafts" :key="draft.name">
            <td>
              <strong>{{ draft.name }}</strong>
            </td>
            <td>
              <input
                v-model.number="draft.feeNaira"
                class="admin-delivery-fee-input"
                type="number"
                min="0"
                step="1"
                :aria-label="`${ADMIN_COPY.deliveryFeeLabel} for ${draft.name}`"
              >
            </td>
            <td class="admin-table-number">
              {{ formatDraftFee(draft) }}
            </td>
            <td>
              <label class="admin-delivery-toggle">
                <input v-model="draft.active" type="checkbox">
                <span>{{ ADMIN_COPY.deliveryActiveLabel }}</span>
              </label>
            </td>
            <td>
              <button
                type="button"
                class="admin-inline-action"
                :disabled="savingZoneName === draft.name"
                @click="saveZone(draft)"
              >
                {{ savingZoneName === draft.name ? ADMIN_COPY.deliverySaving : ADMIN_COPY.deliverySave }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
