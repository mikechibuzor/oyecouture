<script setup lang="ts">
import { ArrowLeft, Check, Copy, MessageCircle } from 'lucide'
import {
  ADMIN_COPY,
  ADMIN_ROUTES,
  ADMIN_STATUS_ACTIONS,
  FULFILLMENT_STAGE_LABELS,
  ORDER_KIND_LABELS,
  ORDER_STATUS_LABELS,
} from '~/constants/admin'
import type { AdminOrderDetail, OrderStatus } from '~/types/admin'
import { OrderKind } from '~/types/order'
import { FulfillmentStage } from '~/types/tracking'
import type { AdminOrderTracking } from '~/types/tracking'
import { formatOrderDate } from '~/utils/admin'
import { measurementFieldLabels } from '~/utils/measurements'
import { formatProductPrice } from '~/utils/product'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const route = useRoute()
const router = useRouter()
const requestUrl = useRequestURL()
const { session } = useAdminSession()
const { fetchOrder, updateOrderStatus } = useAdminOrders()
const { advanceOrderTracking, fetchOrderTracking } = useAdminOrderTracking()
const { notifyError, notifySuccess } = useAppNotifications()

const order = ref<AdminOrderDetail | null>(null)
const isLoading = ref(true)
const hasFailed = ref(false)
const updatingStatus = ref<OrderStatus | null>(null)
const statusError = ref('')
const copiedPaymentLink = ref(false)
const tracking = ref<AdminOrderTracking | null>(null)
const isTrackingLoading = ref(false)
const trackingUnavailable = ref(false)
const trackingLoadFailed = ref(false)
const isAdvancingTracking = ref(false)
const trackingError = ref('')
const copiedTrackingLink = ref(false)
const trackingForm = reactive({
  publicNote: '',
  estimatedDeliveryDate: '',
  carrierName: '',
  carrierTrackingNumber: '',
  carrierTrackingUrl: '',
})

let copiedTimer: ReturnType<typeof setTimeout> | undefined

// Route params can be arrays, but this page only accepts one order number.
const orderNo = computed(() => {
  const rawParam = route.params.orderNo
  return Array.isArray(rawParam) ? rawParam[0] ?? '' : rawParam ?? ''
})

// Status actions are derived from the current order state.
const statusActions = computed(() =>
  order.value ? ADMIN_STATUS_ACTIONS[order.value.status] : [],
)

const nextFulfillmentStage = computed(
  () => tracking.value?.allowedNextStages[0] ?? null,
)

const nextFulfillmentLabel = computed(() =>
  nextFulfillmentStage.value
    ? FULFILLMENT_STAGE_LABELS[nextFulfillmentStage.value]
    : '',
)

const isDispatchTransition = computed(
  () => nextFulfillmentStage.value === FulfillmentStage.Dispatched,
)

// The pay link is shareable even before payment provider credentials are added.
const paymentLink = computed(() => {
  if (!order.value?.payToken) {
    return ''
  }

  return new URL(`/pay/${order.value.payToken}`, requestUrl.origin).toString()
})

// Submitted measurements are shown as labelled rows and hidden when absent.
const measurementRows = computed(() => {
  const measurements = order.value?.measurements
  if (!measurements) {
    return []
  }

  const structuredRows = Object.entries(measurementFieldLabels)
    .map(([fieldId, label]) => ({
      label,
      value: formatMeasurement(measurements.values?.[fieldId], measurements.unit ?? 'in'),
    }))
    .filter((row) => row.value !== '')

  if (structuredRows.length > 0) {
    const profileLabels = {
      suit: 'Suit / tuxedo',
      native: 'Native / kaftan / agbada',
      mixed: 'Mixed order',
    } as const

    return [
      ...(measurements.profile
        ? [{ label: 'Garment profile', value: profileLabels[measurements.profile] }]
        : []),
      ...structuredRows,
      ...(measurements.notes ? [{ label: 'Notes', value: measurements.notes }] : []),
    ]
  }

  return [
    { label: 'Neck', value: formatMeasurement(measurements.neckCm, 'cm') },
    { label: 'Chest', value: formatMeasurement(measurements.chestCm, 'cm') },
    { label: 'Waist', value: formatMeasurement(measurements.waistCm, 'cm') },
    { label: 'Hips', value: formatMeasurement(measurements.hipsCm, 'cm') },
    { label: 'Sleeve', value: formatMeasurement(measurements.sleeveCm, 'cm') },
    { label: 'Inseam', value: formatMeasurement(measurements.inseamCm, 'cm') },
    { label: 'Length', value: formatMeasurement(measurements.lengthCm, 'cm') },
    { label: 'Height', value: formatMeasurement(measurements.heightCm, 'cm') },
    { label: 'Weight', value: formatMeasurement(measurements.weightKg, 'kg') },
    { label: 'Notes', value: measurements.notes || '' },
  ].filter((row) => row.value !== '')
})

// Keeps tailoring details visible for bespoke work without adding an empty section to ready-to-wear.
const shouldShowTailoring = computed(
  () => order.value?.kind === OrderKind.Bespoke || measurementRows.value.length > 0,
)

// WhatsApp expects digits only; local Nigerian numbers are normalised to 234.
const customerWhatsAppNumber = computed(() => {
  const phone = order.value?.phone.replace(/\D/g, '') ?? ''
  if (phone.startsWith('0')) {
    return `234${phone.slice(1)}`
  }

  return phone
})

// The reply gives the customer the order number and next useful action.
const whatsAppReplyLink = computed(() => {
  if (!order.value || customerWhatsAppNumber.value.length < 8) {
    return ''
  }

  const message = [
    `Hi ${order.value.customerName}, this is Oyecouture about order ${order.value.orderNo}.`,
    `Current status: ${ORDER_STATUS_LABELS[order.value.status]}.`,
    paymentLink.value ? `Payment link: ${paymentLink.value}` : '',
  ].filter(Boolean).join('\n')

  return `https://wa.me/${customerWhatsAppNumber.value}?text=${encodeURIComponent(message)}`
})

useHead(() => ({
  title: `${order.value?.orderNo ?? ADMIN_COPY.ordersHeading} · ${ADMIN_COPY.brand}`,
}))

// Loads the order detail, sending expired sessions back to sign-in.
const loadOrder = async (): Promise<void> => {
  isLoading.value = true
  hasFailed.value = false

  try {
    order.value = await fetchOrder(orderNo.value)
    await loadTracking()
  } catch (loadError) {
    if ((loadError as { status?: number } | null)?.status === 401) {
      session.value = null
      await router.push(ADMIN_ROUTES.login)
      return
    }

    hasFailed.value = true
  } finally {
    isLoading.value = false
  }
}

const loadTracking = async (): Promise<void> => {
  isTrackingLoading.value = true
  trackingUnavailable.value = false
  trackingLoadFailed.value = false

  try {
    tracking.value = await fetchOrderTracking(orderNo.value)
    syncTrackingForm(tracking.value)
  } catch (loadError) {
    const statusCode = (loadError as { status?: number; statusCode?: number } | null)
    const status = statusCode?.status ?? statusCode?.statusCode
    tracking.value = null
    if (status === 409) {
      trackingUnavailable.value = true
    } else {
      trackingLoadFailed.value = true
    }
  } finally {
    isTrackingLoading.value = false
  }
}

// Copies the customer-facing payment URL.
const copyPaymentLink = async (): Promise<void> => {
  if (!paymentLink.value || typeof navigator === 'undefined') {
    return
  }

  await navigator.clipboard.writeText(paymentLink.value)
  copiedPaymentLink.value = true
  clearTimeout(copiedTimer)
  copiedTimer = setTimeout(() => {
    copiedPaymentLink.value = false
  }, 1600)
}

const copyTrackingLink = async (): Promise<void> => {
  if (!tracking.value?.trackingUrl || typeof navigator === 'undefined') {
    return
  }

  await navigator.clipboard.writeText(tracking.value.trackingUrl)
  copiedTrackingLink.value = true
  notifySuccess(ADMIN_COPY.trackingLinkCopied)
  clearTimeout(copiedTimer)
  copiedTimer = setTimeout(() => {
    copiedTrackingLink.value = false
  }, 1600)
}

const advanceFulfillment = async (): Promise<void> => {
  if (!tracking.value || !nextFulfillmentStage.value) {
    return
  }

  isAdvancingTracking.value = true
  trackingError.value = ''

  try {
    tracking.value = await advanceOrderTracking(orderNo.value, {
      stage: nextFulfillmentStage.value,
      expectedVersion: tracking.value.version,
      publicNote: trackingForm.publicNote,
      estimatedDeliveryDate: trackingForm.estimatedDeliveryDate,
      carrierName: isDispatchTransition.value ? trackingForm.carrierName : '',
      carrierTrackingNumber: isDispatchTransition.value
        ? trackingForm.carrierTrackingNumber
        : '',
      carrierTrackingUrl: isDispatchTransition.value ? trackingForm.carrierTrackingUrl : '',
    })
    order.value = await fetchOrder(orderNo.value)
    syncTrackingForm(tracking.value)
    notifySuccess(ADMIN_COPY.fulfillmentUpdatedNotice, tracking.value.currentStageLabel)
  } catch (updateError) {
    trackingError.value = ADMIN_COPY.fulfillmentUpdateFailed
    notifyError(ADMIN_COPY.fulfillmentUpdateFailed)
    if ((updateError as { status?: number } | null)?.status === 409) {
      await loadTracking()
    }
  } finally {
    isAdvancingTracking.value = false
  }
}

// Applies one backend-approved status transition.
const applyStatus = async (nextStatus: OrderStatus): Promise<void> => {
  if (!order.value) {
    return
  }

  updatingStatus.value = nextStatus
  statusError.value = ''

  try {
    order.value = await updateOrderStatus(order.value.orderNo, nextStatus)
  } catch {
    statusError.value = ADMIN_COPY.statusUpdateFailed
  } finally {
    updatingStatus.value = null
  }
}

// Formats a nullable measurement value without noisy decimals.
function formatMeasurement(value: number | null | undefined, unit: string): string {
  if (value === null || value === undefined) {
    return ''
  }

  return `${new Intl.NumberFormat('en-NG', { maximumFractionDigits: 1 }).format(value)} ${unit}`
}

function syncTrackingForm(currentTracking: AdminOrderTracking): void {
  trackingForm.publicNote = ''
  trackingForm.estimatedDeliveryDate = currentTracking.estimatedDeliveryDate?.slice(0, 10) ?? ''
  trackingForm.carrierName = currentTracking.shipment.carrierName ?? ''
  trackingForm.carrierTrackingNumber = currentTracking.shipment.carrierTrackingNumber ?? ''
  trackingForm.carrierTrackingUrl = currentTracking.shipment.carrierTrackingUrl ?? ''
}

onMounted(loadOrder)
onBeforeUnmount(() => clearTimeout(copiedTimer))
</script>

<template>
  <div class="admin-page">
    <NuxtLink class="admin-back-link" :to="ADMIN_ROUTES.orders">
      <StorefrontIcon :icon-node="ArrowLeft" :size="15" :stroke-width="1.75" />
      {{ ADMIN_COPY.backToOrders }}
    </NuxtLink>

    <p v-if="isLoading" class="admin-orders-message">
      {{ ADMIN_COPY.orderDetailLoading }}
    </p>

    <div v-else-if="hasFailed || !order" class="admin-orders-message" role="alert">
      <p>{{ ADMIN_COPY.orderDetailFailed }}</p>
      <button type="button" class="admin-inline-action" @click="loadOrder">
        {{ ADMIN_COPY.retry }}
      </button>
    </div>

    <template v-else>
      <header class="admin-detail-head">
        <div>
          <p class="admin-detail-eyebrow">{{ ORDER_KIND_LABELS[order.kind] }}</p>
          <h1>{{ order.orderNo }}</h1>
          <p class="admin-page-subheading">
            {{ order.customerName }} · {{ formatOrderDate(order.createdAt) }}
          </p>
        </div>

        <span class="admin-status" :data-status="order.status">
          {{ ORDER_STATUS_LABELS[order.status] }}
        </span>
      </header>

      <div class="admin-detail-sections">
        <section class="admin-detail-section" aria-labelledby="order-overview-heading">
          <header class="admin-section-heading">
            <p>{{ ADMIN_COPY.orderOverviewEyebrow }}</p>
            <h2 id="order-overview-heading">{{ ADMIN_COPY.orderOverviewHeading }}</h2>
            <span>{{ ADMIN_COPY.orderOverviewIntro }}</span>
          </header>

          <div class="admin-detail-grid">
            <section class="admin-detail-panel">
              <h2>{{ ADMIN_COPY.customerDetails }}</h2>
              <dl class="admin-detail-list">
                <div>
                  <dt>Name</dt>
                  <dd>{{ order.customerName }}</dd>
                </div>
                <div>
                  <dt>Email</dt>
                  <dd>{{ order.email }}</dd>
                </div>
                <div>
                  <dt>Phone</dt>
                  <dd>{{ order.phone }}</dd>
                </div>
                <div v-if="order.address">
                  <dt>Address</dt>
                  <dd>{{ order.address }}</dd>
                </div>
                <div v-if="order.deliveryZone">
                  <dt>Delivery area</dt>
                  <dd>{{ order.deliveryZone }}</dd>
                </div>
              </dl>
            </section>

            <section class="admin-detail-panel">
              <h2>{{ ADMIN_COPY.statusActions }}</h2>
              <div v-if="statusActions.length" class="admin-status-actions">
                <button
                  v-for="action in statusActions"
                  :key="action.status"
                  type="button"
                  class="admin-inline-action"
                  :disabled="updatingStatus !== null"
                  @click="applyStatus(action.status)"
                >
                  {{ updatingStatus === action.status ? ADMIN_COPY.statusUpdating : action.label }}
                </button>
              </div>
              <p v-else class="admin-detail-muted">{{ ADMIN_COPY.noStatusActions }}</p>
              <p v-if="statusError" class="admin-detail-error" role="alert">
                {{ statusError }}
              </p>
            </section>
          </div>
        </section>

        <section class="admin-detail-section" aria-labelledby="order-progress-heading">
          <header class="admin-section-heading">
            <p>{{ ADMIN_COPY.orderProgressEyebrow }}</p>
            <h2 id="order-progress-heading">{{ ADMIN_COPY.orderProgressHeading }}</h2>
            <span>{{ ADMIN_COPY.orderProgressIntro }}</span>
          </header>

          <section class="admin-detail-panel admin-fulfillment-panel">
            <h2>{{ ADMIN_COPY.fulfillment }}</h2>

            <p v-if="isTrackingLoading" class="admin-detail-muted">
              {{ ADMIN_COPY.fulfillmentLoading }}
            </p>
            <p v-else-if="trackingUnavailable" class="admin-detail-muted">
              {{ ADMIN_COPY.fulfillmentUnavailable }}
            </p>
            <div v-else-if="trackingLoadFailed" class="admin-detail-error" role="alert">
              <p>{{ ADMIN_COPY.fulfillmentLoadFailed }}</p>
              <button type="button" class="admin-inline-action" @click="loadTracking">
                {{ ADMIN_COPY.retry }}
              </button>
            </div>

            <template v-else-if="tracking">
              <ol class="admin-fulfillment-timeline">
                <li v-for="event in tracking.events" :key="event.id">
                  <span class="admin-fulfillment-marker" aria-hidden="true" />
                  <div>
                    <div class="admin-fulfillment-event-head">
                      <strong>{{ event.label }}</strong>
                      <time :datetime="event.occurredAt">{{ formatOrderDate(event.occurredAt) }}</time>
                    </div>
                    <p v-if="event.publicNote">{{ event.publicNote }}</p>
                    <small v-if="event.actorEmail">Updated by {{ event.actorEmail }}</small>
                  </div>
                </li>
              </ol>

              <div v-if="tracking.trackingUrl" class="admin-fulfillment-link">
                <span>
                  <small>{{ ADMIN_COPY.customerTrackingLink }}</small>
                  <strong>{{ tracking.trackingUrl }}</strong>
                </span>
                <button type="button" class="admin-inline-action" @click="copyTrackingLink">
                  <StorefrontIcon
                    :icon-node="copiedTrackingLink ? Check : Copy"
                    :size="14"
                    :stroke-width="1.75"
                  />
                  {{ copiedTrackingLink ? ADMIN_COPY.copied : ADMIN_COPY.copyTrackingLink }}
                </button>
              </div>

              <form
                v-if="nextFulfillmentStage"
                class="admin-fulfillment-form"
                @submit.prevent="advanceFulfillment"
              >
                <div class="admin-fulfillment-form-grid">
                  <label class="admin-field">
                    <span>{{ ADMIN_COPY.estimatedDeliveryLabel }}</span>
                    <input v-model="trackingForm.estimatedDeliveryDate" type="date">
                  </label>

                  <template v-if="isDispatchTransition">
                    <label class="admin-field">
                      <span>{{ ADMIN_COPY.carrierNameLabel }}</span>
                      <input v-model.trim="trackingForm.carrierName" type="text" maxlength="200">
                    </label>
                    <label class="admin-field">
                      <span>{{ ADMIN_COPY.carrierTrackingNumberLabel }}</span>
                      <input
                        v-model.trim="trackingForm.carrierTrackingNumber"
                        type="text"
                        maxlength="200"
                      >
                    </label>
                    <label class="admin-field">
                      <span>{{ ADMIN_COPY.carrierTrackingUrlLabel }}</span>
                      <input v-model.trim="trackingForm.carrierTrackingUrl" type="url">
                    </label>
                  </template>

                  <label class="admin-field admin-field-wide">
                    <span>{{ ADMIN_COPY.fulfillmentNoteLabel }}</span>
                    <textarea v-model.trim="trackingForm.publicNote" maxlength="500" rows="3" />
                    <small>{{ ADMIN_COPY.fulfillmentNoteHint }}</small>
                  </label>
                </div>

                <p v-if="trackingError" class="admin-detail-error" role="alert">
                  {{ trackingError }}
                </p>
                <div class="admin-fulfillment-submit">
                  <button
                    type="submit"
                    class="admin-inline-action is-primary"
                    :disabled="isAdvancingTracking"
                  >
                    {{
                      isAdvancingTracking
                        ? ADMIN_COPY.advancingFulfillment
                        : `${ADMIN_COPY.advanceFulfillment} ${nextFulfillmentLabel}`
                    }}
                  </button>
                </div>
              </form>
            </template>
          </section>
        </section>

        <section class="admin-detail-section" aria-labelledby="order-commercial-heading">
          <header class="admin-section-heading">
            <p>{{ ADMIN_COPY.orderCommercialEyebrow }}</p>
            <h2 id="order-commercial-heading">{{ ADMIN_COPY.orderCommercialHeading }}</h2>
            <span>{{ ADMIN_COPY.orderCommercialIntro }}</span>
          </header>

          <div class="admin-detail-grid">
            <section class="admin-detail-panel">
              <h2>{{ ADMIN_COPY.orderItems }}</h2>
              <div class="admin-detail-items">
                <div v-for="item in order.items" :key="`${item.slug}-${item.name}-${item.size}`">
                  <span>
                    <strong>{{ item.name }}</strong>
                    <small>
                      x{{ item.quantity }}<template v-if="item.size"> · {{ item.size }}</template>
                    </small>
                  </span>
                  <span>{{ formatProductPrice(item.lineKobo) }}</span>
                </div>
              </div>
              <dl class="admin-detail-totals">
                <div>
                  <dt>{{ ADMIN_COPY.subtotal }}</dt>
                  <dd>{{ formatProductPrice(order.subtotalKobo) }}</dd>
                </div>
                <div>
                  <dt>{{ ADMIN_COPY.delivery }}</dt>
                  <dd>{{ formatProductPrice(order.deliveryFeeKobo) }}</dd>
                </div>
                <div>
                  <dt>{{ ADMIN_COPY.total }}</dt>
                  <dd>{{ formatProductPrice(order.totalKobo) }}</dd>
                </div>
              </dl>
            </section>

            <section class="admin-detail-panel">
              <h2>{{ ADMIN_COPY.paymentLink }}</h2>
              <p class="admin-detail-link">{{ paymentLink }}</p>
              <div class="admin-status-actions">
                <button type="button" class="admin-inline-action" @click="copyPaymentLink">
                  <StorefrontIcon
                    :icon-node="copiedPaymentLink ? Check : Copy"
                    :size="14"
                    :stroke-width="1.75"
                  />
                  {{ copiedPaymentLink ? ADMIN_COPY.copied : ADMIN_COPY.copyPaymentLink }}
                </button>
                <a
                  v-if="whatsAppReplyLink"
                  class="admin-inline-action"
                  :href="whatsAppReplyLink"
                  target="_blank"
                  rel="noreferrer"
                >
                  <StorefrontIcon :icon-node="MessageCircle" :size="14" :stroke-width="1.75" />
                  {{ ADMIN_COPY.whatsappReply }}
                </a>
              </div>
            </section>
          </div>
        </section>

        <section
          v-if="shouldShowTailoring"
          class="admin-detail-section"
          aria-labelledby="order-tailoring-heading"
        >
          <header class="admin-section-heading">
            <p>{{ ADMIN_COPY.orderTailoringEyebrow }}</p>
            <h2 id="order-tailoring-heading">{{ ADMIN_COPY.orderTailoringHeading }}</h2>
            <span>{{ ADMIN_COPY.orderTailoringIntro }}</span>
          </header>

          <section class="admin-detail-panel">
            <h2>{{ ADMIN_COPY.measurements }}</h2>
            <dl v-if="measurementRows.length" class="admin-detail-list is-compact">
              <div v-for="row in measurementRows" :key="row.label">
                <dt>{{ row.label }}</dt>
                <dd>{{ row.value }}</dd>
              </div>
            </dl>
            <p v-else class="admin-detail-muted">{{ ADMIN_COPY.noMeasurements }}</p>
          </section>
        </section>
      </div>
    </template>
  </div>
</template>
