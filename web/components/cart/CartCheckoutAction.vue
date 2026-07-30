<script setup lang="ts">
import type { IconNode } from 'lucide'
import { ArrowRight, MessageCircle } from 'lucide'
import { StorefrontRoute } from '~/types/storefront'
import { buildWhatsAppOrderUrl } from '~/utils/cart'

interface CartCheckoutActionProps {
  isBespoke: boolean
  label: string
}

const props = defineProps<CartCheckoutActionProps>()

const { close, lines, subtotalKobo } = useCart()
const runtimeConfig = useRuntimeConfig()

const ACTION_ICONS = {
  bespoke: MessageCircle,
  readyToWear: ArrowRight,
} as const satisfies Readonly<Record<string, IconNode>>

// Guards against a missing WHATSAPP_NUMBER: without it wa.me would open a broken
// page, so the action is disabled rather than sending the customer nowhere.
const isWhatsAppConfigured = computed(() =>
  runtimeConfig.public.whatsappNumber.replace(/\D/g, '').length > 0,
)

// Prefills the WhatsApp conversation with the bag so the customer never has to
// retype what they want and the owner sees the order immediately.
const whatsAppUrl = computed(() =>
  buildWhatsAppOrderUrl(
    runtimeConfig.public.whatsappNumber,
    lines.value,
    subtotalKobo.value,
  ),
)
const actionIcon = computed(() =>
  props.isBespoke ? ACTION_ICONS.bespoke : ACTION_ICONS.readyToWear,
)
</script>

<template>
  <button
    v-if="isBespoke && !isWhatsAppConfigured"
    class="cart-primary-action"
    type="button"
    disabled
  >
    <span>{{ label }}</span>
  </button>

  <a
    v-else-if="isBespoke"
    class="cart-primary-action"
    :href="whatsAppUrl"
    target="_blank"
    rel="noopener noreferrer"
    @click="close"
  >
    <StorefrontIcon
      :icon-node="actionIcon"
      :size="18"
      :stroke-width="1.5"
      aria-hidden="true"
    />
    <span>{{ label }}</span>
  </a>

  <NuxtLink
    v-else
    class="cart-primary-action"
    :to="StorefrontRoute.Checkout"
    @click="close"
  >
    <span>{{ label }}</span>
    <StorefrontIcon
      :icon-node="actionIcon"
      :size="18"
      :stroke-width="1.5"
      aria-hidden="true"
    />
  </NuxtLink>
</template>
