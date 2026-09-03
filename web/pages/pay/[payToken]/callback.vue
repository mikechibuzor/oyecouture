<script setup lang="ts">
import { CONFIRMATION_COPY, PAYMENT_COPY } from '~/constants/payment'
import { FULFILLMENT_PATHS, FULFILLMENT_STAGE_LABELS, TRACKING_COPY } from '~/constants/tracking'
import { FulfillmentStage } from '~/types/tracking'
import { OrderKind } from '~/types/order'
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
// The same stages, labels and wording the order tracker uses, rather than a
// second set of prose describing the same journey. A customer who lands here
// and later opens their tracking link sees one story told once.
const confirmationSteps = computed(() => {
  const kind = order.value?.kind === OrderKind.Bespoke
    ? OrderKind.Bespoke
    : OrderKind.ReadyToWear

  return FULFILLMENT_PATHS[kind].map((stage) => ({
    stage,
    label: FULFILLMENT_STAGE_LABELS[stage],
    // Payment is the only step that has happened by the time this page loads,
    // so the next one carries the tracker's awaiting-update line.
    isDone: stage === FulfillmentStage.PaymentConfirmed,
    isCurrent: stage === FULFILLMENT_PATHS[kind][1],
  }))
})

const isPaid = ref(false)
const message = ref<string>(PAYMENT_COPY.verifying)
const order = ref<PaymentOrder | null>(null)
const paymentChannel = ref('')

useSeoMeta({
  // Reactive, so the tab stops saying "Confirming payment..." once it is
  // confirmed. Passing a getter is what makes useSeoMeta re-evaluate it.
  title: () =>
    `${isPaid.value ? PAYMENT_COPY.verified : PAYMENT_COPY.verifying} | Oyecouture Bespoke`,
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

// Paystack sends the customer back the instant they finish their side, but a
// bank transfer or USSD payment is still settling then - order OYE-1002
// redirected 50 seconds before its transfer cleared. Checking once and calling
// it a failure told a paying customer their payment had not worked, so the page
// keeps checking instead. The window is generous because the alternative is
// telling someone who has been debited that nothing happened.
const settleWindowMs = 3 * 60 * 1000
const settlePollMs = 4000

const isSettling = ref(false)
const hasGivenUpWaiting = ref(false)
let settlePollTimer: ReturnType<typeof setTimeout> | undefined
let stopWaitingAt = 0

const checkOnce = async (): Promise<boolean> => {
  const { result } = await verifyPayment(payToken.value, reference.value)
  order.value = result.order
  paymentChannel.value = result.channel ?? ''
  isPaid.value = result.paid

  return result.paid
}

const scheduleNextCheck = (): void => {
  if (Date.now() >= stopWaitingAt) {
    isSettling.value = false
    hasGivenUpWaiting.value = true
    message.value = PAYMENT_COPY.settlingSlowBody
    return
  }

  settlePollTimer = setTimeout(() => {
    void poll()
  }, settlePollMs)
}

// A failed check while settling is not news: the reference simply is not
// payable yet. Only the very first attempt distinguishes a real failure.
const poll = async (): Promise<void> => {
  try {
    if (await checkOnce()) {
      isSettling.value = false
      hasGivenUpWaiting.value = false
      message.value = PAYMENT_COPY.verified
      return
    }
  } catch {
    // Swallowed deliberately: keep watching rather than reporting failure.
  }

  scheduleNextCheck()
}

const verify = async (): Promise<void> => {
  if (!reference.value) {
    isVerifying.value = false
    message.value = PAYMENT_COPY.failed
    return
  }

  isVerifying.value = true
  hasGivenUpWaiting.value = false

  try {
    if (await checkOnce()) {
      message.value = PAYMENT_COPY.verified
      return
    }

    // Not paid yet rather than not paid at all. Start watching.
    isSettling.value = true
    message.value = PAYMENT_COPY.settlingBody
    stopWaitingAt = Date.now() + settleWindowMs
    scheduleNextCheck()
  } catch {
    // The first check erroring is ambiguous: an async payment that has not
    // reached the provider yet looks the same as a genuine failure. Wait it
    // out rather than accusing a paying customer of not paying.
    isSettling.value = true
    message.value = PAYMENT_COPY.settlingBody
    stopWaitingAt = Date.now() + settleWindowMs
    scheduleNextCheck()
  } finally {
    isVerifying.value = false
  }
}

onBeforeUnmount(() => clearTimeout(settlePollTimer))

onMounted(verify)
</script>

<template>
  <div class="confirmation">
    <!-- Verifying and failure keep the original plain treatment: the celebratory
         layout below is only honest once the money has actually cleared. -->
    <section v-if="isVerifying || !isPaid" class="confirmation-pending">
      <p class="confirmation-eyebrow">{{ PAYMENT_COPY.breadcrumb }}</p>
      <h1>
        <template v-if="isVerifying">{{ PAYMENT_COPY.verifying }}</template>
        <template v-else-if="isSettling">{{ PAYMENT_COPY.settling }}</template>
        <template v-else-if="hasGivenUpWaiting">{{ PAYMENT_COPY.settlingSlow }}</template>
        <template v-else>{{ PAYMENT_COPY.pending }}</template>
      </h1>
      <p class="confirmation-pending-body">{{ message }}</p>

      <!-- While a transfer is settling there is nothing for the customer to do
           and nothing to retry: the old "Try again" sent them back to a
           checkout link Paystack had already closed, which reads as a second
           failure on top of the first. -->
      <button
        v-if="hasGivenUpWaiting"
        type="button"
        class="confirmation-button is-primary"
        @click="verify"
      >
        {{ PAYMENT_COPY.checkAgain }}
      </button>

      <NuxtLink
        v-else-if="!isVerifying && !isSettling"
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
                v-for="step in confirmationSteps"
                :key="step.stage"
                :class="{ 'is-done': step.isDone, 'is-current': step.isCurrent }"
              >
                <span class="confirmation-step-marker" aria-hidden="true" />
                <div>
                  <p class="confirmation-step-title">{{ step.label }}</p>
                  <p v-if="step.isCurrent" class="confirmation-step-body">
                    {{ TRACKING_COPY.awaitingUpdate }}
                  </p>
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
