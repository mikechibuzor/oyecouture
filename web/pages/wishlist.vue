<script setup lang="ts">
import type { IconNode } from 'lucide'
import { ArrowRight, Heart, ShoppingBag, X } from 'lucide'
import {
  DEFAULT_PRODUCT_SIZE,
  PRODUCT_COPY,
} from '~/constants/products'
import { AsyncDataStatus } from '~/types/product'
import type { Product } from '~/types/product'
import { StorefrontRoute } from '~/types/storefront'
import { getProductErrorMessage } from '~/utils/product'

const WISHLIST_ICONS = {
  heart: Heart,
  shop: ArrowRight,
  bag: ShoppingBag,
  remove: X,
} as const satisfies Readonly<Record<string, IconNode>>

const { data, error, status } = await useProductList()
const { add: addToCart, open: openCart } = useCart()
const {
  clear: clearWishlist,
  count: wishlistCount,
  items: wishlistItems,
  isEmpty: isWishlistEmpty,
  isLoaded: isWishlistLoaded,
  remove: removeFromWishlist,
  updateProduct: updateWishlistProduct,
} = useWishlist()

const productsBySlug = computed(() =>
  new Map((data.value?.products ?? []).map((product) => [product.slug, product])),
)

const wishlistedProducts = computed<Product[]>(() =>
  wishlistItems.value
    .map((wishlistItem) =>
      productsBySlug.value.get(wishlistItem.slug) ?? wishlistItem.product,
    )
    .filter((product): product is Product => Boolean(product)),
)
const hasWishlistedProducts = computed(() => wishlistedProducts.value.length > 0)

const headerSubcopy = computed(() => {
  if (wishlistCount.value === 0) {
    return "Tap the heart on any piece to keep it here. It is a private rail of the cuts you're considering."
  }

  return `${wishlistCount.value} saved ${wishlistCount.value === 1 ? 'piece' : 'pieces'} ready when you are.`
})

const isWaitingForStoredWishlist = computed(() =>
  !isWishlistLoaded.value,
)
const needsSavedProductData = computed(() =>
  isWishlistLoaded.value && !isWishlistEmpty.value,
)
const isLoading = computed(() =>
  needsSavedProductData.value
  && !hasWishlistedProducts.value
  && (
    status.value === AsyncDataStatus.Idle
    || status.value === AsyncDataStatus.Pending
  ),
)
const isWishlistBusy = computed(() =>
  isWaitingForStoredWishlist.value || isLoading.value,
)
const hasError = computed(() =>
  needsSavedProductData.value
  && !hasWishlistedProducts.value
  && status.value === AsyncDataStatus.Error,
)
const errorMessage = computed(() =>
  getProductErrorMessage(error.value, PRODUCT_COPY.catalogue.error),
)
const showEmptyState = computed(() =>
  !isWaitingForStoredWishlist.value && isWishlistEmpty.value,
)
const showWishlistGrid = computed(() =>
  !isLoading.value && hasWishlistedProducts.value,
)
const showUnavailableState = computed(() =>
  needsSavedProductData.value
  && !isLoading.value
  && !hasError.value
  && !hasWishlistedProducts.value,
)

const addAllToBag = (): void => {
  let addedAnyProduct = false

  for (const product of wishlistedProducts.value) {
    const size = product.hasSizes ? DEFAULT_PRODUCT_SIZE : null

    if (addToCart(product, size)) {
      removeFromWishlist(product.slug)
      addedAnyProduct = true
    }
  }

  if (addedAnyProduct) {
    openCart()
  }
}

const removeWishlistedProduct = (productSlug: string): void => {
  removeFromWishlist(productSlug)
}

watch(
  () => data.value?.products,
  (products) => {
    for (const product of products ?? []) {
      updateWishlistProduct(product)
    }
  },
  { immediate: true },
)

useSeoMeta({
  title: 'Wishlist | Oyecouture Bespoke',
  description:
    "Your saved Oyecouture pieces form a private rail of the cuts you're considering.",
})
</script>

<template>
  <div class="wishlist-page">
    <header class="wishlist-header">
      <nav class="wishlist-breadcrumbs" :aria-label="PRODUCT_COPY.breadcrumbs.label">
        <NuxtLink :to="StorefrontRoute.Home">{{ PRODUCT_COPY.breadcrumbs.home }}</NuxtLink>
        <span aria-hidden="true">/</span>
        <span>Wishlist</span>
      </nav>

      <div class="wishlist-header-row">
        <div>
          <p class="wishlist-eyebrow">Saved For Later</p>
          <h1>Your Wishlist</h1>
          <p>{{ headerSubcopy }}</p>
        </div>

        <div v-if="showWishlistGrid" class="wishlist-header-actions">
          <button type="button" @click="clearWishlist">Clear wishlist</button>
          <button type="button" class="is-primary" @click="addAllToBag">
            <StorefrontIcon
              :icon-node="WISHLIST_ICONS.bag"
              :size="15"
              :stroke-width="1.5"
              aria-hidden="true"
            />
            <span>Add all to bag</span>
          </button>
        </div>
      </div>
    </header>

    <section class="wishlist-body" aria-live="polite">
      <p v-if="isWishlistBusy" class="wishlist-state" role="status">
        {{ PRODUCT_COPY.catalogue.loading }}
      </p>

      <p v-else-if="hasError" class="wishlist-state is-error" role="alert">
        {{ errorMessage }}
      </p>

      <div v-else-if="showEmptyState" class="wishlist-empty">
        <span class="wishlist-empty-icon" aria-hidden="true">
          <StorefrontIcon
            :icon-node="WISHLIST_ICONS.heart"
            :size="38"
            :stroke-width="1.3"
          />
        </span>
        <h2>Nothing saved yet</h2>
        <p>
          Tap the heart on any piece to keep it here. It is a private rail of the cuts
          you're considering.
        </p>
        <NuxtLink :to="StorefrontRoute.Collection">
          <span>Explore the collection</span>
          <StorefrontIcon
            :icon-node="WISHLIST_ICONS.shop"
            :size="14"
            :stroke-width="2"
            aria-hidden="true"
          />
        </NuxtLink>
      </div>

      <div v-else-if="showWishlistGrid" class="product-grid wishlist-grid">
        <div
          v-for="product in wishlistedProducts"
          :key="product.id"
          class="wishlist-grid-item"
        >
          <CatalogueProductCard :product="product" />
          <button
            type="button"
            class="wishlist-remove-button"
            :aria-label="`Remove ${product.name} from wishlist`"
            :title="`Remove ${product.name} from wishlist`"
            @click="removeWishlistedProduct(product.slug)"
          >
            <StorefrontIcon
              :icon-node="WISHLIST_ICONS.remove"
              :size="14"
              :stroke-width="1.7"
              aria-hidden="true"
            />
            <span>Remove</span>
          </button>
        </div>
      </div>

      <div v-else-if="showUnavailableState" class="wishlist-empty">
        <span class="wishlist-empty-icon" aria-hidden="true">
          <StorefrontIcon
            :icon-node="WISHLIST_ICONS.heart"
            :size="38"
            :stroke-width="1.3"
          />
        </span>
        <h2>Saved pieces unavailable</h2>
        <p>
          These saved pieces are no longer in the current collection. Clear your
          wishlist or explore the latest pieces.
        </p>
        <NuxtLink :to="StorefrontRoute.Collection">
          <span>Explore the collection</span>
          <StorefrontIcon
            :icon-node="WISHLIST_ICONS.shop"
            :size="14"
            :stroke-width="2"
            aria-hidden="true"
          />
        </NuxtLink>
      </div>
    </section>
  </div>
</template>
