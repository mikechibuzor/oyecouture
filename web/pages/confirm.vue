<script setup lang="ts">
import type { IconNode } from 'lucide'
import { MessageCircle } from 'lucide'
import { CONFIRM_COPY, MEASUREMENT_FIELDS } from '~/constants/confirm'
import { CART_COPY } from '~/constants/cart'
import { PRODUCT_COPY } from '~/constants/products'
import type { DeliveryZone, MeasurementsRequest, Order } from '~/types/order'
import { StorefrontRoute } from '~/types/storefront'
import type { SelectOption } from '~/types/ui'
import { buildWhatsAppOrderUrl } from '~/utils/cart'
import { formatProductPrice } from '~/utils/product'

const { clear, isEmpty, lines, subtotalKobo } = useCart()
const { submitOrder } = useOrders()
const { data: zoneData } = await useDeliveryZones()
const runtimeConfig = useRuntimeConfig()

const WHATSAPP_ICON: IconNode = MessageCircle

const customerName = ref('')
const email = ref('')
const phone = ref('')
const address = ref('')
const deliveryZoneName = ref('')
const showMeasurements = ref(false)
const notes = ref('')
const measurementValues = ref<Record<string, string>>({})

const isSubmitting = ref(false)
const errorMessage = ref('')
const placedOrder = ref<Order | null>(null)
// Kept because the bag is cleared once the order is placed, and the success
// screen still needs the lines to build the WhatsApp message.
const placedLines = ref<typeof lines.value>([])
const placedSubtotalKobo = ref(0)
const placedDeliveryDetails = ref({ address: '', deliveryZone: '' })

useSeoMeta({ title: `${CONFIRM_COPY.heading} — Oyecouture Bespoke` })

const formattedSubtotal = computed(() => formatProductPrice(subtotalKobo.value))
const zones = computed<readonly DeliveryZone[]>(() => zoneData.value?.zones ?? [])
const zoneSelectOptions = computed<SelectOption[]>(() =>
  zones.value.map((zone) => ({
    value: zone.name,
    label: zone.name,
  })),
)
const configuredWhatsAppNumber = computed(() =>
  String(runtimeConfig.public.whatsappNumber ?? ''),
)
const whatsAppUrl = computed(() =>
  placedOrder.value
    ? buildWhatsAppOrderUrl(
        configuredWhatsAppNumber.value,
        placedLines.value,
        placedSubtotalKobo.value,
        placedOrder.value.orderNo,
        placedDeliveryDetails.value,
      )
    : '',
)
const isWhatsAppConfigured = computed(
  () => configuredWhatsAppNumber.value.replace(/\D/g, '').length > 0,
)

// Converts the free-text measurement inputs into numbers, dropping blanks so an
// untouched field is sent as null rather than 0.
const buildMeasurements = (): MeasurementsRequest | null => {
  if (!showMeasurements.value) {
    return null
  }

  const parsed = Object.fromEntries(
    MEASUREMENT_FIELDS.map((field) => {
      const raw = measurementValues.value[field.key]?.trim()
      const value = raw ? Number(raw) : Number.NaN

      return [field.key, Number.isFinite(value) ? value : null]
    }),
  ) as Omit<MeasurementsRequest, 'notes'>

  return { ...parsed, notes: notes.value.trim() }
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

  try {
    const { order } = await submitOrder({
      customerName: customerName.value.trim(),
      email: email.value.trim(),
      phone: phone.value.trim(),
      items: submittedLines.map((line) => ({
        slug: line.slug,
        size: line.size ?? '',
        quantity: line.quantity,
      })),
      measurements: buildMeasurements(),
      address: address.value.trim(),
      deliveryZone: deliveryZoneName.value,
    })

    placedOrder.value = order
    placedLines.value = submittedLines
    placedSubtotalKobo.value = submittedSubtotal
    placedDeliveryDetails.value = {
      address: address.value.trim(),
      deliveryZone: deliveryZoneName.value,
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

          <fieldset class="confirm-measurements">
            <legend>{{ CONFIRM_COPY.measurementsHeading }}</legend>
            <p class="confirm-measurements-hint">{{ CONFIRM_COPY.measurementsHint }}</p>

            <label class="confirm-toggle">
              <input v-model="showMeasurements" type="checkbox">
              <span>{{ CONFIRM_COPY.measurementsToggle }}</span>
            </label>

            <div v-if="showMeasurements" class="confirm-measurement-grid">
              <label v-for="field in MEASUREMENT_FIELDS" :key="field.key">
                <span>{{ field.label }}</span>
                <input
                  v-model="measurementValues[field.key]"
                  type="number"
                  inputmode="decimal"
                  step="0.1"
                  min="0"
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

          <div class="confirm-subtotal">
            <span>{{ CONFIRM_COPY.subtotal }}</span>
            <strong>{{ formattedSubtotal }}</strong>
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
