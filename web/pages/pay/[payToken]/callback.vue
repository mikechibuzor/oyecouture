<script setup lang="ts">
import { CONFIRMATION_COPY, CONFIRMATION_STEPS, PAYMENT_COPY } from '~/constants/payment'
import { StorefrontRoute } from '~/types/storefront'
import type { PaymentOrder } from '~/types/payment'

const route = useRoute()
const { verifyPayment } = usePayments()

const payToken = computed(() => {
  const rawToken = route.params.payToken
  return Array.isArray(rawToken) ? rawToken[0] ?? '' : rawToken ?? ''
})

const reference = computed(() => {
  const rawReference = route.query.reference
  return Array.isArray(rawReference) ? rawReference[0] ?? '' : rawReference ?? ''
})

const isVerifying = ref(true)
const isPaid = ref(false)
const message = ref<string>(PAYMENT_COPY.verifying)
const order = ref<PaymentOrder | null>(null)
const paymentChannel = ref('')

useSeoMeta({
  title: `${PAYMENT_COPY.verifying} | Oyecouture Bespoke`,
  // A confirmation page carries a customer's address and order value, so it
  // must never reach an index even though the token makes it unguessable.
  robots: 'noindex, nofollow',
})

// Only the first name, because "Thank you, Adeola" reads as a person writing and
// "Thank you, Adeola Bankole" reads as a database.
const firstName = computed(() => (order.value?.customerName ?? '').trim().split(/\s+/)[0] ?? '')

const itemCount = computed(() =>
  (order.value?.items ?? []).reduce((runningTotal, item) => runningTotal + item.quantity, 0),
)

// Paystack returns channels as snake_case identifiers; the design shows them as
// words. Unknown channels fall back to the provider alone rather than printing
// a raw identifier at a customer.
const paymentMethodLabel = computed(() => {
  const channelLabels: Record<string, string> = {
    card: 'Card',
    bank: 'Bank',
    bank_transfer: 'Bank transfer',
    ussd: 'USSD',
    qr: 'QR',
    mobile_money: 'Mobile money',
  }

  const label = channelLabels[paymentChannel.value]
  return label ? `Paystack · ${label}` : 'Paystack'
})

const orderDateLabel = computed(() => {
  const created = order.value?.createdAt
  if (!created) {
    return ''
  }

  const parsed = new Date(created)
  return Number.isNaN(parsed.getTime())
    ? ''
    : new Intl.DateTimeFormat('en-NG', {
        day: 'numeric',
        month: 'short',
        year: 'numeric',
      }).format(parsed)
})

// Delivery is shown as a price only when one was actually charged. A zero fee is
// a free delivery, not a missing value, so it reads as "Complimentary".
const deliveryLabel = computed(() =>
  (order.value?.deliveryFeeKobo ?? 0) > 0
    ? formatProductPrice(order.value?.deliveryFeeKobo ?? 0)
    : CONFIRMATION_COPY.deliveryComplimentary,
)

const receiptLink = computed(() =>
  reference.value
    ? `/pay/${payToken.value}/receipt?reference=${encodeURIComponent(reference.value)}`
    : `/pay/${payToken.value}/receipt`,
)

const verify = async (): Promise<void> => {
  if (!reference.value) {
    isVerifying.value = false
    message.value = PAYMENT_COPY.failed
    return
  }

  isVerifying.value = true

  try {
    const { result } = await verifyPayment(payToken.value, reference.value)
    isPaid.value = result.paid
    order.value = result.order
    paymentChannel.value = result.channel ?? ''
    message.value = result.paid ? PAYMENT_COPY.verified : result.message || PAYMENT_COPY.pending
  } catch {
    isPaid.value = false
    message.value = PAYMENT_COPY.failed
  } finally {
    isVerifying.value = false
  }
}

onMounted(verify)
</script>

<template>
  <div class="confirmation">
    <!-- Verifying and failure keep the original plain treatment: the celebratory
         layout below is only honest once the money has actually cleared. -->
    <section v-if="isVerifying || !isPaid" class="confirmation-pending">
      <p class="confirmation-eyebrow">{{ PAYMENT_COPY.breadcrumb }}</p>
      <h1>{{ isVerifying ? PAYMENT_COPY.verifying : PAYMENT_COPY.pending }}</h1>
      <p class="confirmation-pending-body">{{ message }}</p>

      <NuxtLink
        v-if="!isVerifying"
        class="confirmation-button is-primary"
        :to="`/pay/${payToken}`"
      >
        {{ PAYMENT_COPY.retry }}
      </NuxtLink>
    </section>

    <template v-else-if="order">
      <header class="confirmation-header">
        <StorefrontBrandLogo />
        <p class="confirmation-order-no">Order {{ order.orderNo }}</p>
      </header>

      <section class="confirmation-hero">
        <span class="confirmation-badge" aria-hidden="true">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M4 12.5 9.5 18 20 6.5" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </span>
        <p class="confirmation-eyebrow">{{ CONFIRMATION_COPY.eyebrow }}</p>
        <h1>{{ CONFIRMATION_COPY.greeting(firstName) }}</h1>
        <p class="confirmation-intro">{{ CONFIRMATION_COPY.intro }}</p>
      </section>

      <dl class="confirmation-facts">
        <div>
          <dt>{{ CONFIRMATION_COPY.orderReference }}</dt>
          <dd class="is-strong">{{ order.orderNo }}</dd>
        </div>
        <div>
          <dt>{{ CONFIRMATION_COPY.orderDate }}</dt>
          <dd>{{ orderDateLabel }}</dd>
        </div>
        <div>
          <dt>{{ CONFIRMATION_COPY.paymentMethod }}</dt>
          <dd>{{ paymentMethodLabel }}</dd>
        </div>
        <div>
          <dt>{{ CONFIRMATION_COPY.amountPaid }}</dt>
          <dd>{{ formatProductPrice(order.totalKobo) }}</dd>
        </div>
      </dl>

      <div class="confirmation-body">
        <div class="confirmation-main">
          <section class="confirmation-card">
            <div class="confirmation-card-head">
              <h2>{{ CONFIRMATION_COPY.yourPieces }}</h2>
              <p>{{ CONFIRMATION_COPY.itemCount(itemCount) }}</p>
            </div>

            <ul class="confirmation-pieces">
              <li v-for="(item, index) in order.items" :key="`${item.name}-${index}`">
                <div>
                  <p class="confirmation-piece-name">{{ item.name }}</p>
                  <p class="confirmation-piece-meta">
                    <span v-if="item.size">{{ item.size }}</span>
                    <span v-if="item.size"> · </span>
                    <span>Qty {{ item.quantity }}</span>
                  </p>
                </div>
                <p class="confirmation-piece-price">{{ formatProductPrice(item.lineKobo) }}</p>
              </li>
            </ul>
          </section>

          <section class="confirmation-card">
            <h2>{{ CONFIRMATION_COPY.whatHappensNext }}</h2>

            <ol class="confirmation-steps">
              <!-- The first step is always complete on this page, and the second
                   is always what happens next, so the states are fixed rather
                   than read from the order. -->
              <li
                v-for="(step, index) in CONFIRMATION_STEPS"
                :key="step.title"
                :class="{ 'is-done': index === 0, 'is-current': index === 1 }"
              >
                <span class="confirmation-step-marker" aria-hidden="true" />
                <div>
                  <p class="confirmation-step-title">{{ step.title }}</p>
                  <p class="confirmation-step-body">{{ step.body }}</p>
                </div>
              </li>
            </ol>
          </section>
        </div>

        <aside class="confirmation-aside">
          <section class="confirmation-summary">
            <h2>{{ CONFIRMATION_COPY.paymentSummary }}</h2>

            <dl>
              <div>
                <dt>{{ CONFIRMATION_COPY.subtotal }}</dt>
                <dd>{{ formatProductPrice(order.subtotalKobo) }}</dd>
              </div>
              <div>
                <dt>{{ CONFIRMATION_COPY.tailoring }}</dt>
                <dd>{{ CONFIRMATION_COPY.tailoringValue }}</dd>
              </div>
              <div>
                <dt>{{ CONFIRMATION_COPY.delivery }}</dt>
                <dd>{{ deliveryLabel }}</dd>
              </div>
            </dl>

            <div class="confirmation-summary-total">
              <p>{{ CONFIRMATION_COPY.amountPaid }}</p>
              <p class="confirmation-summary-amount">{{ formatProductPrice(order.totalKobo) }}</p>
            </div>

            <p class="confirmation-summary-note">{{ CONFIRMATION_COPY.chargeNote }}</p>
          </section>

          <section class="confirmation-card confirmation-delivery">
            <h2>{{ CONFIRMATION_COPY.deliveryDetails }}</h2>

            <div>
              <p class="confirmation-label">{{ CONFIRMATION_COPY.recipient }}</p>
              <p>{{ order.customerName }}</p>
            </div>
            <div v-if="order.address">
              <p class="confirmation-label">{{ CONFIRMATION_COPY.address }}</p>
              <p>{{ order.address }}</p>
            </div>
            <div v-if="order.email || order.phone">
              <p class="confirmation-label">{{ CONFIRMATION_COPY.contact }}</p>
              <p v-if="order.email">{{ order.email }}</p>
              <p v-if="order.phone">{{ order.phone }}</p>
            </div>
          </section>

          <div class="confirmation-actions">
            <NuxtLink class="confirmation-button is-primary" :to="receiptLink">
              {{ CONFIRMATION_COPY.downloadReceipt }}
            </NuxtLink>
            <NuxtLink class="confirmation-button is-secondary" :to="StorefrontRoute.Contact">
              {{ CONFIRMATION_COPY.contactAtelier }}
            </NuxtLink>
            <NuxtLink class="confirmation-button is-quiet" :to="StorefrontRoute.Collection">
              {{ PAYMENT_COPY.continueShopping }}
            </NuxtLink>
          </div>
        </aside>
      </div>

      <footer class="confirmation-footer">
        <p>{{ CONFIRMATION_COPY.footerNote }}</p>
        <p>{{ CONFIRMATION_COPY.footerQuestion }}</p>
      </footer>
    </template>
  </div>
</template>
