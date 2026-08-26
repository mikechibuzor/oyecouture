<script setup lang="ts">
import type { IconNode } from 'lucide'
import { ClipboardPaste, MessageCircle, Ruler } from 'lucide'
import { CONFIRM_COPY } from '~/constants/confirm'
import { CART_COPY } from '~/constants/cart'
import { PRODUCT_COPY } from '~/constants/products'
import type { DeliveryZone, MeasurementsRequest, Order } from '~/types/order'
import { StorefrontRoute } from '~/types/storefront'
import type { SelectOption } from '~/types/ui'
import { buildWhatsAppOrderUrl, type WhatsAppOrderDetails } from '~/utils/cart'
import { getMeasurementRequirements, parseMeasurementClipboard } from '~/utils/measurements'
import { formatProductPrice } from '~/utils/product'

const { clear, isEmpty, lines, subtotalKobo } = useCart()
const { submitOrder } = useOrders()
const { data: zoneData } = await useDeliveryZones()
const { data: siteSettings } = await useSiteSettings()
const { withLoading } = useAppLoading()

const WHATSAPP_ICON: IconNode = MessageCircle
const MEASUREMENT_GUIDE_ICON: IconNode = Ruler
const MEASUREMENT_PASTE_ICON: IconNode = ClipboardPaste

const customerName = ref('')
const email = ref('')
const phone = ref('')
const address = ref('')
const deliveryZoneName = ref('')
const showMeasurements = ref(false)
const notes = ref('')
const measurementValues = ref<Record<string, string>>({})
const measurementFeedback = ref<{ message: string, status: 'success' | 'error' } | null>(null)

const isSubmitting = ref(false)
const errorMessage = ref('')
const placedOrder = ref<Order | null>(null)
// Kept because the bag is cleared once the order is placed, and the success
// screen still needs the lines to build the WhatsApp message.
const placedLines = ref<typeof lines.value>([])
const placedSubtotalKobo = ref(0)
const placedOrderDetails = ref<WhatsAppOrderDetails | null>(null)

useSeoMeta({ title: `${CONFIRM_COPY.heading} | Oyecouture Bespoke` })

const formattedSubtotal = computed(() => formatProductPrice(subtotalKobo.value))
const zones = computed<readonly DeliveryZone[]>(() => zoneData.value?.zones ?? [])
const zoneSelectOptions = computed<SelectOption[]>(() =>
  zones.value.map((zone) => ({
    value: zone.name,
    label: `${zone.name} · ${formatProductPrice(zone.feeKobo)}`,
  })),
)
const selectedZone = computed(() =>
  zones.value.find((zone) => zone.name === deliveryZoneName.value) ?? null,
)
const deliveryFeeKobo = computed(() => selectedZone.value?.feeKobo ?? 0)
const totalKobo = computed(() => subtotalKobo.value + deliveryFeeKobo.value)
const configuredWhatsAppNumber = computed(() =>
  String(siteSettings.value?.whatsappNumber ?? ''),
)
const whatsAppUrl = computed(() =>
  placedOrder.value
    ? buildWhatsAppOrderUrl(
        configuredWhatsAppNumber.value,
        placedLines.value,
        placedSubtotalKobo.value,
        placedOrder.value.orderNo,
        placedOrderDetails.value ?? undefined,
      )
    : '',
)
const isWhatsAppConfigured = computed(
  () => configuredWhatsAppNumber.value.replace(/\D/g, '').length > 0,
)
const measurementRequirements = computed(() => getMeasurementRequirements(lines.value))
const measurementFields = computed(() => measurementRequirements.value?.fields ?? [])
const measurementGuideUrl = computed(() => {
  const profile = measurementRequirements.value?.profile ?? 'native'
  const query = new URLSearchParams({ profile, returnTo: StorefrontRoute.Confirm })

  return `${StorefrontRoute.Measurements}?${query.toString()}#worksheet`
})

const measurementNumber = (fieldId: string): number | null => {
  const rawValue = measurementValues.value[fieldId]?.trim()
  const value = rawValue ? Number(rawValue) : Number.NaN

  return Number.isFinite(value) && value > 0 ? value : null
}

const inchesToCentimetres = (value: number | null): number | null =>
  value === null ? null : Math.round(value * 25.4) / 10

// Converts the free-text measurement inputs into numbers, dropping blanks so an
// untouched field is sent as null rather than 0.
const buildMeasurements = (): MeasurementsRequest | null => {
  if (!showMeasurements.value || !measurementRequirements.value) {
    return null
  }

  const values = Object.fromEntries(
    measurementFields.value
      .map((field) => [field.id, measurementNumber(field.id)] as const)
      .filter((entry): entry is readonly [string, number] => entry[1] !== null),
  )
  const preferredLength = measurementNumber('jacket-length') ?? measurementNumber('top-length')

  return {
    profile: measurementRequirements.value.profile,
    unit: 'in',
    values,
    neckCm: inchesToCentimetres(measurementNumber('neck')),
    chestCm: inchesToCentimetres(measurementNumber('chest')),
    waistCm: inchesToCentimetres(measurementNumber('waist')),
    hipsCm: null,
    sleeveCm: inchesToCentimetres(measurementNumber('sleeve')),
    inseamCm: null,
    lengthCm: inchesToCentimetres(preferredLength),
    heightCm: null,
    weightKg: null,
    notes: notes.value.trim(),
  }
}

const pasteMeasurements = async (): Promise<void> => {
  measurementFeedback.value = null
  if (!navigator.clipboard?.readText) {
    measurementFeedback.value = {
      message: CONFIRM_COPY.measurementsPasteUnavailable,
      status: 'error',
    }
    return
  }

  try {
    const clipboardText = await navigator.clipboard.readText()
    if (!clipboardText.trim()) {
      measurementFeedback.value = {
        message: CONFIRM_COPY.measurementsPasteEmpty,
        status: 'error',
      }
      return
    }

    const payload = parseMeasurementClipboard(clipboardText)
    if (!payload) {
      measurementFeedback.value = {
        message: CONFIRM_COPY.measurementsPasteInvalid,
        status: 'error',
      }
      return
    }

    const acceptedFieldIds = new Set(measurementFields.value.map((field) => field.id))
    const matchingValues = Object.entries(payload.values)
      .filter(([fieldId]) => acceptedFieldIds.has(fieldId))

    if (matchingValues.length === 0) {
      measurementFeedback.value = {
        message: CONFIRM_COPY.measurementsPasteInvalid,
        status: 'error',
      }
      return
    }

    for (const [fieldId, value] of matchingValues) {
      measurementValues.value[fieldId] = String(value)
    }
    showMeasurements.value = true
    measurementFeedback.value = {
      message: CONFIRM_COPY.measurementsPasteSuccess,
      status: 'success',
    }
  } catch {
    measurementFeedback.value = {
      message: CONFIRM_COPY.measurementsPasteUnavailable,
      status: 'error',
    }
  }
}

const handleSubmit = async (): Promise<void> => {
  if (isSubmitting.value || isEmpty.value) {
    return
  }

  if (!deliveryZoneName.value) {
    errorMessage.value = CONFIRM_COPY.zoneRequired
    return
  }

  isSubmitting.value = true
  errorMessage.value = ''

  // Snapshot the bag before clearing it, so the success screen can still show it.
  const submittedLines = [...lines.value]
  const submittedSubtotal = subtotalKobo.value
  const submittedMeasurements = buildMeasurements()
  const submittedOrderDetails: WhatsAppOrderDetails = {
    customerName: customerName.value.trim(),
    email: email.value.trim(),
    phone: phone.value.trim(),
    address: address.value.trim(),
    deliveryZone: deliveryZoneName.value,
    measurements: submittedMeasurements,
  }

  try {
    const { order } = await withLoading(
      () => submitOrder({
        customerName: submittedOrderDetails.customerName ?? '',
        email: submittedOrderDetails.email ?? '',
        phone: submittedOrderDetails.phone ?? '',
        items: submittedLines.map((line) => ({
          slug: line.slug,
          size: line.size ?? '',
          quantity: line.quantity,
        })),
        measurements: submittedMeasurements,
        address: submittedOrderDetails.address,
        deliveryZone: submittedOrderDetails.deliveryZone,
      }),
      CONFIRM_COPY.submitting,
    )

    placedOrder.value = order
    placedLines.value = submittedLines
    placedSubtotalKobo.value = submittedSubtotal
    placedOrderDetails.value = {
      ...submittedOrderDetails,
      deliveryFeeKobo: order.deliveryFeeKobo,
      totalKobo: order.totalKobo,
    }
    clear()
  } catch (submitError) {
    const apiError = submitError as { data?: { error?: { message?: string } } }
    errorMessage.value = apiError?.data?.error?.message ?? CONFIRM_COPY.errorFallback
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="confirm-page">
    <!-- Success: the order exists and carries a reference. -->
    <section v-if="placedOrder" class="confirm-success">
      <p class="confirm-eyebrow">{{ CONFIRM_COPY.successEyebrow }}</p>
      <h1>{{ CONFIRM_COPY.successHeading }}</h1>
      <p class="confirm-order-no">{{ placedOrder.orderNo }}</p>
      <p class="confirm-success-body">{{ CONFIRM_COPY.successBody }}</p>

      <a
        v-if="isWhatsAppConfigured"
        class="confirm-submit"
        :href="whatsAppUrl"
        target="_blank"
        rel="noopener noreferrer"
      >
        <StorefrontIcon
          :icon-node="WHATSAPP_ICON"
          :size="18"
          :stroke-width="1.5"
          aria-hidden="true"
        />
        <span>{{ CONFIRM_COPY.successAction }}</span>
      </a>
    </section>

    <template v-else>
      <header class="confirm-header">
        <nav class="category-breadcrumbs" :aria-label="PRODUCT_COPY.breadcrumbs.label">
          <NuxtLink :to="StorefrontRoute.Home">{{ PRODUCT_COPY.breadcrumbs.home }}</NuxtLink>
          <span aria-hidden="true">/</span>
          <span>{{ CONFIRM_COPY.breadcrumb }}</span>
        </nav>
        <h1>{{ CONFIRM_COPY.heading }}</h1>
        <p>{{ CONFIRM_COPY.description }}</p>
      </header>

      <p v-if="isEmpty" class="confirm-empty">
        {{ CONFIRM_COPY.emptyBag }}
        <NuxtLink :to="StorefrontRoute.Collection">{{ CONFIRM_COPY.emptyBagAction }}</NuxtLink>
      </p>

      <form v-else class="confirm-layout" @submit.prevent="handleSubmit">
        <div class="confirm-fields">
          <fieldset>
            <legend>{{ CONFIRM_COPY.contactHeading }}</legend>

            <label>
              <span>{{ CONFIRM_COPY.name }}</span>
              <input v-model="customerName" type="text" name="name" autocomplete="name" required>
            </label>

            <label>
              <span>{{ CONFIRM_COPY.email }}</span>
              <input v-model="email" type="email" name="email" autocomplete="email" required>
            </label>

            <label>
              <span>{{ CONFIRM_COPY.phone }}</span>
              <input v-model="phone" type="tel" name="phone" autocomplete="tel" required>
              <small>{{ CONFIRM_COPY.phoneHint }}</small>
            </label>
          </fieldset>

          <fieldset>
            <legend>{{ CONFIRM_COPY.deliveryHeading }}</legend>

            <label>
              <span>{{ CONFIRM_COPY.address }}</span>
              <textarea
                v-model="address"
                name="address"
                autocomplete="street-address"
                rows="3"
                required
              />
            </label>

            <div class="checkout-zone-field">
              <span>{{ CONFIRM_COPY.zone }}</span>
              <UiSelect
                v-model="deliveryZoneName"
                :options="zoneSelectOptions"
                :label="CONFIRM_COPY.zone"
                :placeholder="CONFIRM_COPY.zonePlaceholder"
              />
              <small>{{ CONFIRM_COPY.zoneHint }}</small>
            </div>
          </fieldset>

          <fieldset v-if="measurementRequirements" id="measurements" class="confirm-measurements">
            <legend>{{ CONFIRM_COPY.measurementsHeading }}</legend>
            <p class="confirm-measurements-hint">{{ CONFIRM_COPY.measurementsHint }}</p>

            <div class="confirm-measurement-actions">
              <NuxtLink class="confirm-measurement-action" :to="measurementGuideUrl">
                <StorefrontIcon :icon-node="MEASUREMENT_GUIDE_ICON" :size="16" :stroke-width="1.6" />
                {{ CONFIRM_COPY.measurementsGuide }}
              </NuxtLink>
              <button class="confirm-measurement-action" type="button" @click="pasteMeasurements">
                <StorefrontIcon :icon-node="MEASUREMENT_PASTE_ICON" :size="16" :stroke-width="1.6" />
                {{ CONFIRM_COPY.measurementsPaste }}
              </button>
            </div>

            <p
              v-if="measurementFeedback"
              class="confirm-measurement-feedback"
              :data-status="measurementFeedback.status"
              role="status"
            >
              {{ measurementFeedback.message }}
            </p>

            <label class="confirm-toggle">
              <input v-model="showMeasurements" type="checkbox">
              <span>{{ CONFIRM_COPY.measurementsToggle }}</span>
            </label>

            <div v-if="showMeasurements" class="confirm-measurement-grid">
              <label v-for="field in measurementFields" :key="field.id">
                <span>{{ field.label }} (in)</span>
                <input
                  v-model="measurementValues[field.id]"
                  type="number"
                  inputmode="decimal"
                  step="0.25"
                  min="0"
                  :name="`measurement-${field.id}`"
                >
              </label>
            </div>

            <label v-if="showMeasurements">
              <span>{{ CONFIRM_COPY.notes }}</span>
              <textarea v-model="notes" rows="3" />
            </label>
          </fieldset>
        </div>

        <aside class="confirm-summary">
          <h2>{{ CONFIRM_COPY.summaryHeading }}</h2>

          <ul>
            <li v-for="line in lines" :key="line.key">
              <span>
                {{ line.name }}
                <template v-if="line.size">({{ CART_COPY.size }} {{ line.size }})</template>
                × {{ line.quantity }}
              </span>
              <strong>{{ formatProductPrice(line.priceKobo * line.quantity) }}</strong>
            </li>
          </ul>

          <div class="checkout-line">
            <span>{{ CONFIRM_COPY.subtotal }}</span>
            <strong>{{ formattedSubtotal }}</strong>
          </div>

          <div class="checkout-line">
            <span>{{ CONFIRM_COPY.delivery }}</span>
            <strong>{{ formatProductPrice(deliveryFeeKobo) }}</strong>
          </div>

          <div class="confirm-subtotal">
            <span>{{ CONFIRM_COPY.total }}</span>
            <strong>{{ formatProductPrice(totalKobo) }}</strong>
          </div>

          <p v-if="errorMessage" class="confirm-error" role="alert">{{ errorMessage }}</p>

          <button class="confirm-submit" type="submit" :disabled="isSubmitting">
            {{ isSubmitting ? CONFIRM_COPY.submitting : CONFIRM_COPY.submit }}
          </button>
        </aside>
      </form>
    </template>
  </div>
</template>
