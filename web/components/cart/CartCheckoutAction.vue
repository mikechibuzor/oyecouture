<script setup lang="ts">
import type { IconNode } from 'lucide'
import { ArrowRight, MessageCircle } from 'lucide'
import { StorefrontRoute } from '~/types/storefront'

interface CartCheckoutActionProps {
  isBespoke: boolean
  label: string
}

const props = defineProps<CartCheckoutActionProps>()

const { close } = useCart()

const ACTION_ICONS = {
  bespoke: MessageCircle,
  readyToWear: ArrowRight,
} as const satisfies Readonly<Record<string, IconNode>>

const actionIcon = computed(() =>
  props.isBespoke ? ACTION_ICONS.bespoke : ACTION_ICONS.readyToWear,
)
</script>

<template>
  <NuxtLink
    v-if="isBespoke"
    class="cart-primary-action"
    :to="StorefrontRoute.Confirm"
    @click="close"
  >
    <StorefrontIcon
      :icon-node="actionIcon"
      :size="18"
      :stroke-width="1.5"
      aria-hidden="true"
    />
    <span>{{ label }}</span>
  </NuxtLink>

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
