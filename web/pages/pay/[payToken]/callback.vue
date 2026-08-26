<script setup lang="ts">
import { PAYMENT_COPY } from '~/constants/payment'
import { StorefrontRoute } from '~/types/storefront'

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

useSeoMeta({ title: `${PAYMENT_COPY.verifying} | Oyecouture Bespoke` })

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
  <div class="confirm-page">
    <section class="confirm-success">
      <p class="confirm-eyebrow">{{ PAYMENT_COPY.breadcrumb }}</p>
      <h1>
        {{ isVerifying ? PAYMENT_COPY.verifying : isPaid ? PAYMENT_COPY.verified : PAYMENT_COPY.pending }}
      </h1>
      <p class="confirm-success-body">{{ message }}</p>

      <NuxtLink class="confirm-submit" :to="isPaid ? StorefrontRoute.Collection : `/pay/${payToken}`">
        {{ isPaid ? PAYMENT_COPY.continueShopping : PAYMENT_COPY.retry }}
      </NuxtLink>
    </section>
  </div>
</template>
