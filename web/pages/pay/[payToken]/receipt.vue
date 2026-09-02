<script setup lang="ts">
import { PAYMENT_COPY, RECEIPT_COPY } from '~/constants/payment'
import type { PaymentOrder } from '~/types/payment'

// A printable receipt rather than a generated PDF. The server is a single core
// running a distroless image, so headless Chrome or wkhtmltopdf would be a large
// dependency for one button. Every browser prints to PDF, and the print
// stylesheet below is what makes the output look like a document rather than a
// screenshot of a web page.
const route = useRoute()
const { verifyPayment, fetchPayment } = usePayments()

const payToken = computed(() => {
  const rawToken = route.params.payToken
  return Array.isArray(rawToken) ? rawToken[0] ?? '' : rawToken ?? ''
})

const reference = computed(() => {
  const rawReference = route.query.reference
  return Array.isArray(rawReference) ? rawReference[0] ?? '' : rawReference ?? ''
})

const order = ref<PaymentOrder | null>(null)
const paymentChannel = ref('')
const isLoading = ref(true)
const loadFailed = ref(false)

useSeoMeta({
  title: `${RECEIPT_COPY.title} | Oyecouture Bespoke`,
  robots: 'noindex, nofollow',
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
        month: 'long',
        year: 'numeric',
      }).format(parsed)
})

const paymentLabel = computed(() => {
  const channelLabels: Record<string, string> = {
    card: 'Card',
    bank: 'Bank',
    bank_transfer: 'Bank transfer',
    ussd: 'USSD',
    qr: 'QR',
    mobile_money: 'Mobile money',
  }

  const label = channelLabels[paymentChannel.value]
  return label ? `Paystack · ${label} · Verified` : 'Paystack · Verified'
})

const deliveryLabel = computed(() =>
  (order.value?.deliveryFeeKobo ?? 0) > 0
    ? formatProductPrice(order.value?.deliveryFeeKobo ?? 0)
    : RECEIPT_COPY.complimentary,
)

// Verify when a reference is present so the receipt reflects a confirmed
// payment; otherwise fall back to the order summary, which is enough to reprint
// a receipt later from a link that has lost its query string.
const load = async (): Promise<void> => {
  isLoading.value = true

  try {
    if (reference.value) {
      const { result } = await verifyPayment(payToken.value, reference.value)
      order.value = result.order
      paymentChannel.value = result.channel ?? ''
    } else {
      const summary = await fetchPayment(payToken.value)
      order.value = summary.order
    }
  } catch {
    loadFailed.value = true
  } finally {
    isLoading.value = false
  }
}

const print = (): void => {
  window.print()
}

onMounted(load)
</script>

<template>
  <div class="receipt-page">
    <p v-if="isLoading" class="receipt-status">{{ PAYMENT_COPY.loading }}</p>
    <p v-else-if="loadFailed || !order" class="receipt-status">{{ PAYMENT_COPY.loadFailed }}</p>

    <template v-else>
      <!-- Screen-only controls. The print stylesheet removes them so they never
           appear on the saved document. -->
      <div class="receipt-controls">
        <button type="button" class="receipt-print" @click="print">
          {{ RECEIPT_COPY.printAction }}
        </button>
        <NuxtLink class="receipt-back" :to="`/pay/${payToken}/callback?reference=${reference}`">
          {{ RECEIPT_COPY.back }}
        </NuxtLink>
      </div>
      <p class="receipt-hint">{{ RECEIPT_COPY.printHint }}</p>

      <article class="receipt">
        <header class="receipt-banner">{{ RECEIPT_COPY.banner }}</header>

        <div class="receipt-brand">
          <StorefrontBrandLogo />
        </div>

        <div class="receipt-body">
          <p class="receipt-eyebrow">{{ RECEIPT_COPY.eyebrow }}</p>
          <h1>{{ RECEIPT_COPY.heading }}</h1>

          <section class="receipt-details">
            <h2>{{ RECEIPT_COPY.orderDetails }}</h2>
            <dl>
              <div>
                <dt>{{ RECEIPT_COPY.orderReference }}</dt>
                <dd>{{ order.orderNo }}</dd>
              </div>
              <div>
                <dt>{{ RECEIPT_COPY.orderDate }}</dt>
                <dd>{{ orderDateLabel }}</dd>
              </div>
              <div>
                <dt>{{ RECEIPT_COPY.payment }}</dt>
                <dd>{{ paymentLabel }}</dd>
              </div>
            </dl>
          </section>

          <section class="receipt-pieces">
            <h2>{{ RECEIPT_COPY.yourPieces }}</h2>

            <ul>
              <li v-for="(item, index) in order.items" :key="`${item.name}-${index}`">
                <div>
                  <p class="receipt-piece-name">{{ item.name }}</p>
                  <p class="receipt-piece-meta">
                    <span v-if="item.size">{{ item.size }} · </span>
                    <span>Qty {{ item.quantity }}</span>
                  </p>
                </div>
                <p class="receipt-piece-price">{{ formatProductPrice(item.lineKobo) }}</p>
              </li>
            </ul>
          </section>

          <dl class="receipt-totals">
            <div>
              <dt>{{ RECEIPT_COPY.subtotal }}</dt>
              <dd>{{ formatProductPrice(order.subtotalKobo) }}</dd>
            </div>
            <div>
              <dt>{{ RECEIPT_COPY.tailoring }}</dt>
              <dd>{{ RECEIPT_COPY.complimentary }}</dd>
            </div>
            <div>
              <dt>{{ RECEIPT_COPY.delivery }}</dt>
              <dd>{{ deliveryLabel }}</dd>
            </div>
            <div class="receipt-total">
              <dt>{{ RECEIPT_COPY.totalPaid }}</dt>
              <dd>{{ formatProductPrice(order.totalKobo) }}</dd>
            </div>
          </dl>

          <section v-if="order.address || order.phone" class="receipt-delivering">
            <h2>{{ RECEIPT_COPY.deliveringTo }}</h2>
            <p>{{ order.customerName }}</p>
            <p v-if="order.address">{{ order.address }}</p>
            <p v-if="order.phone">{{ order.phone }}</p>
          </section>

          <p class="receipt-signoff">{{ RECEIPT_COPY.signOff }}</p>
          <p class="receipt-brandmark">{{ RECEIPT_COPY.brandName }}</p>
        </div>

        <footer class="receipt-footer">
          <p>{{ RECEIPT_COPY.returnsNote }}</p>
          <p>{{ RECEIPT_COPY.address }}</p>
          <p>{{ RECEIPT_COPY.receiptFor(order.orderNo) }}</p>
        </footer>
      </article>
    </template>
  </div>
</template>
