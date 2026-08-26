<script setup lang="ts">
import { ArrowLeft, Check, Copy, MessageCircle } from 'lucide'
import {
  ADMIN_COPY,
  ADMIN_ROUTES,
  ADMIN_STATUS_ACTIONS,
  ORDER_KIND_LABELS,
  ORDER_STATUS_LABELS,
} from '~/constants/admin'
import type { AdminOrderDetail, OrderStatus } from '~/types/admin'
import { formatOrderDate } from '~/utils/admin'
import { measurementFieldLabels } from '~/utils/measurements'
import { formatProductPrice } from '~/utils/product'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const route = useRoute()
const router = useRouter()
const requestUrl = useRequestURL()
const { session } = useAdminSession()
const { fetchOrder, updateOrderStatus } = useAdminOrders()

const order = ref<AdminOrderDetail | null>(null)
const isLoading = ref(true)
const hasFailed = ref(false)
const updatingStatus = ref<OrderStatus | null>(null)
const statusError = ref('')
const copiedPaymentLink = ref(false)

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

        <section class="admin-detail-panel admin-detail-panel-wide">
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
          <h2>{{ ADMIN_COPY.measurements }}</h2>
          <dl v-if="measurementRows.length" class="admin-detail-list is-compact">
            <div v-for="row in measurementRows" :key="row.label">
              <dt>{{ row.label }}</dt>
              <dd>{{ row.value }}</dd>
            </div>
          </dl>
          <p v-else class="admin-detail-muted">{{ ADMIN_COPY.noMeasurements }}</p>
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
    </template>
  </div>
</template>
