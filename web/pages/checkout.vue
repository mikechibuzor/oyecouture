<script setup lang="ts">
import { CART_COPY } from '~/constants/cart'
import { CHECKOUT_COPY } from '~/constants/checkout'
import { PRODUCT_COPY } from '~/constants/products'
import type { DeliveryZone, Order } from '~/types/order'
import { StorefrontRoute } from '~/types/storefront'
import { formatProductPrice } from '~/utils/product'

const { clear, isEmpty, lines, subtotalKobo } = useCart()
const { submitOrder } = useOrders()
const { data: zoneData } = await useDeliveryZones()

const customerName = ref('')
const email = ref('')
const phone = ref('')
const address = ref('')
const deliveryZoneName = ref('')

const isSubmitting = ref(false)
const errorMessage = ref('')
const placedOrder = ref<Order | null>(null)

useSeoMeta({ title: `${CHECKOUT_COPY.heading} — Oyecouture Bespoke` })

const zones = computed<readonly DeliveryZone[]>(() => zoneData.value?.zones ?? [])
const selectedZone = computed(() =>
  zones.value.find((zone) => zone.name === deliveryZoneName.value) ?? null,
)
// Shown as an estimate; the server recalculates the fee when the order is placed.
const deliveryFeeKobo = computed(() => selectedZone.value?.feeKobo ?? 0)
const totalKobo = computed(() => subtotalKobo.value + deliveryFeeKobo.value)

const handleSubmit = async (): Promise<void> => {
  if (isSubmitting.value || isEmpty.value) {
    return
  }

  isSubmitting.value = true
  errorMessage.value = ''

  try {
    const { order } = await submitOrder({
      customerName: customerName.value.trim(),
      email: email.value.trim(),
      phone: phone.value.trim(),
      items: lines.value.map((line) => ({
        slug: line.slug,
        size: line.size ?? '',
        quantity: line.quantity,
      })),
      measurements: null,
      address: address.value.trim(),
      deliveryZone: deliveryZoneName.value,
    })

    placedOrder.value = order
    clear()
  } catch (submitError) {
    const apiError = submitError as { data?: { error?: { message?: string } } }
    errorMessage.value = apiError?.data?.error?.message ?? CHECKOUT_COPY.errorFallback
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="confirm-page">
    <section v-if="placedOrder" class="confirm-success">
      <p class="confirm-eyebrow">{{ CHECKOUT_COPY.successEyebrow }}</p>
      <h1>{{ CHECKOUT_COPY.successHeading }}</h1>
      <p class="confirm-order-no">{{ placedOrder.orderNo }}</p>
      <p class="confirm-success-body">{{ CHECKOUT_COPY.successBody }}</p>
      <NuxtLink class="confirm-submit" :to="StorefrontRoute.Collection">
        {{ CHECKOUT_COPY.successAction }}
      </NuxtLink>
    </section>

    <template v-else>
      <header class="confirm-header">
        <nav class="category-breadcrumbs" :aria-label="PRODUCT_COPY.breadcrumbs.label">
          <NuxtLink :to="StorefrontRoute.Home">{{ PRODUCT_COPY.breadcrumbs.home }}</NuxtLink>
          <span aria-hidden="true">/</span>
          <span>{{ CHECKOUT_COPY.breadcrumb }}</span>
        </nav>
        <h1>{{ CHECKOUT_COPY.heading }}</h1>
        <p>{{ CHECKOUT_COPY.description }}</p>
      </header>

      <p v-if="isEmpty" class="confirm-empty">
        {{ CART_COPY.empty }}
        <NuxtLink :to="StorefrontRoute.Collection">{{ CART_COPY.emptyAction }}</NuxtLink>
      </p>

      <form v-else class="confirm-layout" @submit.prevent="handleSubmit">
        <div class="confirm-fields">
          <fieldset>
            <legend>{{ CHECKOUT_COPY.contactHeading }}</legend>

            <label>
              <span>{{ CHECKOUT_COPY.name }}</span>
              <input v-model="customerName" type="text" name="name" autocomplete="name" required>
            </label>
            <label>
              <span>{{ CHECKOUT_COPY.email }}</span>
              <input v-model="email" type="email" name="email" autocomplete="email" required>
            </label>
            <label>
              <span>{{ CHECKOUT_COPY.phone }}</span>
              <input v-model="phone" type="tel" name="phone" autocomplete="tel" required>
            </label>
          </fieldset>

          <fieldset>
            <legend>{{ CHECKOUT_COPY.deliveryHeading }}</legend>

            <label>
              <span>{{ CHECKOUT_COPY.address }}</span>
              <textarea
                v-model="address"
                name="address"
                autocomplete="street-address"
                rows="3"
                required
              />
            </label>

            <label>
              <span>{{ CHECKOUT_COPY.zone }}</span>
              <select v-model="deliveryZoneName" name="deliveryZone" required>
                <option value="" disabled>{{ CHECKOUT_COPY.zonePlaceholder }}</option>
                <option v-for="zone in zones" :key="zone.name" :value="zone.name">
                  {{ zone.name }} — {{ formatProductPrice(zone.feeKobo) }}
                </option>
              </select>
              <small>{{ CHECKOUT_COPY.zoneHint }}</small>
            </label>
          </fieldset>
        </div>

        <aside class="confirm-summary">
          <h2>{{ CHECKOUT_COPY.summaryHeading }}</h2>

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
            <span>{{ CHECKOUT_COPY.subtotal }}</span>
            <strong>{{ formatProductPrice(subtotalKobo) }}</strong>
          </div>
          <div class="checkout-line">
            <span>{{ CHECKOUT_COPY.delivery }}</span>
            <strong>{{ formatProductPrice(deliveryFeeKobo) }}</strong>
          </div>

          <div class="confirm-subtotal">
            <span>{{ CHECKOUT_COPY.total }}</span>
            <strong>{{ formatProductPrice(totalKobo) }}</strong>
          </div>

          <p v-if="errorMessage" class="confirm-error" role="alert">{{ errorMessage }}</p>

          <button class="confirm-submit" type="submit" :disabled="isSubmitting">
            {{ isSubmitting ? CHECKOUT_COPY.submitting : CHECKOUT_COPY.submit }}
          </button>

          <p class="checkout-payment-note">{{ CHECKOUT_COPY.paymentNote }}</p>
        </aside>
      </form>
    </template>
  </div>
</template>
