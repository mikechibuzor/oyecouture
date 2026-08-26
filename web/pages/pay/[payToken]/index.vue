<script setup lang="ts">
import { CART_COPY } from '~/constants/cart'
import { PAYMENT_COPY } from '~/constants/payment'
import { PRODUCT_COPY } from '~/constants/products'
import { StorefrontRoute } from '~/types/storefront'
import { formatProductPrice } from '~/utils/product'

const route = useRoute()
const { fetchPayment, initializePayment } = usePayments()

const payToken = computed(() => {
  const rawToken = route.params.payToken
  return Array.isArray(rawToken) ? rawToken[0] ?? '' : rawToken ?? ''
})

const isStarting = ref(false)
const startError = ref('')

const {
  data: payment,
  error: loadError,
  pending: isLoading,
  refresh,
} = await useAsyncData(`payment-${payToken.value}`, () => fetchPayment(payToken.value))

useSeoMeta({ title: `${PAYMENT_COPY.heading} | Oyecouture Bespoke` })

const canStartPayment = computed(() => Boolean(payment.value?.payable && !isStarting.value))

const retryLoad = (): void => {
  void refresh()
}

const startPayment = async (): Promise<void> => {
  if (!canStartPayment.value || typeof window === 'undefined') {
    return
  }

  isStarting.value = true
  startError.value = ''

  try {
    const { payment: initializedPayment } = await initializePayment(payToken.value)
    window.location.href = initializedPayment.authorizationUrl
  } catch (initializeError) {
    const apiError = initializeError as { data?: { error?: { message?: string } } }
    startError.value = apiError?.data?.error?.message ?? PAYMENT_COPY.failed
    isStarting.value = false
  }
}
</script>

<template>
  <div class="confirm-page">
    <header class="confirm-header">
      <nav class="category-breadcrumbs" :aria-label="PRODUCT_COPY.breadcrumbs.label">
        <NuxtLink :to="StorefrontRoute.Home">{{ PRODUCT_COPY.breadcrumbs.home }}</NuxtLink>
        <span aria-hidden="true">/</span>
        <span>{{ PAYMENT_COPY.breadcrumb }}</span>
      </nav>
      <h1>{{ PAYMENT_COPY.heading }}</h1>
      <p>{{ PAYMENT_COPY.description }}</p>
    </header>

    <p v-if="isLoading" class="confirm-empty">{{ PAYMENT_COPY.loading }}</p>

    <section v-else-if="loadError || !payment" class="confirm-success">
      <h1>{{ PAYMENT_COPY.loadFailed }}</h1>
      <button type="button" class="confirm-submit" @click="retryLoad">
        {{ PAYMENT_COPY.retry }}
      </button>
    </section>

    <section v-else class="confirm-layout">
      <div class="confirm-fields">
        <fieldset>
          <legend>{{ PAYMENT_COPY.order }}</legend>
          <div class="checkout-line">
            <span>{{ PAYMENT_COPY.order }}</span>
            <strong>{{ payment.order.orderNo }}</strong>
          </div>
          <div class="checkout-line">
            <span>{{ PAYMENT_COPY.customer }}</span>
            <strong>{{ payment.order.customerName }}</strong>
          </div>
          <p class="checkout-payment-note">{{ payment.message }}</p>
        </fieldset>

        <fieldset>
          <legend>{{ PAYMENT_COPY.items }}</legend>
          <div
            v-for="item in payment.order.items"
            :key="`${item.name}-${item.size}-${item.quantity}`"
            class="checkout-line"
          >
            <span>
              {{ item.name }}
              <template v-if="item.size">({{ CART_COPY.size }} {{ item.size }})</template>
              x{{ item.quantity }}
            </span>
            <strong>{{ formatProductPrice(item.lineKobo) }}</strong>
          </div>
        </fieldset>
      </div>

      <aside class="confirm-summary">
        <h2>{{ PAYMENT_COPY.total }}</h2>
        <div class="confirm-subtotal">
          <span>{{ PAYMENT_COPY.total }}</span>
          <strong>{{ formatProductPrice(payment.order.totalKobo) }}</strong>
        </div>

        <p v-if="startError" class="confirm-error" role="alert">{{ startError }}</p>

        <button
          v-if="payment.payable"
          type="button"
          class="confirm-submit"
          :disabled="isStarting"
          @click="startPayment"
        >
          {{ isStarting ? PAYMENT_COPY.starting : PAYMENT_COPY.payAction }}
        </button>

        <NuxtLink v-else class="confirm-submit" :to="StorefrontRoute.Collection">
          {{ payment.order.status === 'paid' ? PAYMENT_COPY.continueShopping : PAYMENT_COPY.notPayable }}
        </NuxtLink>
      </aside>
    </section>
  </div>
</template>
