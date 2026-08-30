<script setup lang="ts">
import type { IconNode } from 'lucide'
import { Minus, Plus, ShoppingBag, X } from 'lucide'
import { CART_COPY, CART_MAXIMUM_QUANTITY } from '~/constants/cart'
import { PRODUCT_COPY, getProductRoute } from '~/constants/products'
import { CartCheckoutRoute } from '~/types/cart'
import { StorefrontRoute } from '~/types/storefront'
import { formatProductPrice } from '~/utils/product'

const {
  checkoutRoute,
  close,
  count,
  isEmpty,
  isOpen,
  lines,
  remove,
  setQuantity,
  subtotalKobo,
} = useCart()

const CART_ICONS = {
  bag: ShoppingBag,
  close: X,
  decrease: Minus,
  increase: Plus,
  remove: X,
} as const satisfies Readonly<Record<string, IconNode>>

const closeButton = ref<HTMLButtonElement | null>(null)
const drawer = ref<HTMLElement | null>(null)

const isBespokeOrder = computed(
  () => checkoutRoute.value === CartCheckoutRoute.Bespoke,
)
const formattedSubtotal = computed(() => formatProductPrice(subtotalKobo.value))
// Bespoke orders collect measurements before choosing payment or WhatsApp.
const primaryActionLabel = computed(() =>
  isBespokeOrder.value ? CART_COPY.confirmBespoke : CART_COPY.checkout,
)
const orderNote = computed(() =>
  isBespokeOrder.value ? CART_COPY.bespokeNote : CART_COPY.readyToWearNote,
)

const closeDrawer = (): void => {
  close()
}

useModalFocus({
  isOpen,
  container: drawer,
  initialFocus: closeButton,
  onEscape: closeDrawer,
})
</script>

<template>
  <Teleport to="body">
    <Transition name="cart-backdrop">
      <div v-if="isOpen" class="cart-backdrop" @click="closeDrawer" />
    </Transition>

    <Transition name="cart-drawer">
      <aside
        v-if="isOpen"
        ref="drawer"
        class="cart-drawer"
        role="dialog"
        aria-modal="true"
        :aria-label="CART_COPY.title"
        tabindex="-1"
      >
        <header class="cart-drawer-header">
          <h2>{{ CART_COPY.title }} · {{ count }}</h2>
          <button
            ref="closeButton"
            type="button"
            class="cart-icon-button"
            :aria-label="CART_COPY.close"
            @click="closeDrawer"
          >
            <StorefrontIcon
              :icon-node="CART_ICONS.close"
              :size="22"
              :stroke-width="1.5"
              aria-hidden="true"
            />
          </button>
        </header>

        <div class="cart-drawer-body">
          <div v-if="isEmpty" class="cart-empty">
            <StorefrontIcon
              :icon-node="CART_ICONS.bag"
              :size="44"
              :stroke-width="1.2"
              aria-hidden="true"
            />
            <p>{{ CART_COPY.empty }}</p>
            <NuxtLink
              class="cart-primary-action"
              :to="StorefrontRoute.Home"
              @click="closeDrawer"
            >
              {{ CART_COPY.emptyAction }}
            </NuxtLink>
          </div>

          <ul v-else class="cart-line-list">
            <li v-for="line in lines" :key="line.key" class="cart-line">
              <NuxtLink
                class="cart-line-media"
                :to="getProductRoute(line.slug)"
                @click="closeDrawer"
              >
                <img v-if="line.imageUrl" :src="line.imageUrl" :alt="line.name">
                <span v-else aria-hidden="true">OYE</span>
              </NuxtLink>

              <div class="cart-line-detail">
                <p class="cart-line-category">{{ line.category }}</p>
                <h3>
                  <NuxtLink :to="getProductRoute(line.slug)" @click="closeDrawer">
                    {{ line.name }}
                  </NuxtLink>
                </h3>
                <p v-if="line.size" class="cart-line-size">
                  {{ CART_COPY.size }}: {{ line.size }}
                </p>

                <div class="cart-line-footer">
                  <div class="cart-quantity" role="group" :aria-label="CART_COPY.quantity">
                    <button
                      type="button"
                      :aria-label="CART_COPY.decrease"
                      @click="setQuantity(line.key, line.quantity - 1)"
                    >
                      <StorefrontIcon
                        :icon-node="CART_ICONS.decrease"
                        :size="14"
                        :stroke-width="1.5"
                        aria-hidden="true"
                      />
                    </button>
                    <span aria-live="polite">{{ line.quantity }}</span>
                    <button
                      type="button"
                      :aria-label="CART_COPY.increase"
                      :disabled="line.quantity >= CART_MAXIMUM_QUANTITY"
                      @click="setQuantity(line.key, line.quantity + 1)"
                    >
                      <StorefrontIcon
                        :icon-node="CART_ICONS.increase"
                        :size="14"
                        :stroke-width="1.5"
                        aria-hidden="true"
                      />
                    </button>
                  </div>

                  <p class="cart-line-price">
                    {{ formatProductPrice(line.priceKobo * line.quantity) }}
                  </p>
                </div>
              </div>

              <button
                type="button"
                class="cart-icon-button cart-line-remove"
                :aria-label="`${CART_COPY.remove}: ${line.name}`"
                @click="remove(line.key)"
              >
                <StorefrontIcon
                  :icon-node="CART_ICONS.remove"
                  :size="16"
                  :stroke-width="1.5"
                  aria-hidden="true"
                />
              </button>
            </li>
          </ul>
        </div>

        <footer v-if="!isEmpty" class="cart-drawer-footer">
          <div class="cart-subtotal">
            <span>{{ CART_COPY.subtotal }}</span>
            <strong>{{ formattedSubtotal }}</strong>
          </div>
          <p class="cart-note">{{ orderNote }}</p>
          <CartCheckoutAction :is-bespoke="isBespokeOrder" :label="primaryActionLabel" />
        </footer>
      </aside>
    </Transition>
  </Teleport>
</template>
