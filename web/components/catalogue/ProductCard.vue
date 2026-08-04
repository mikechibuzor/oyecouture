<script setup lang="ts">
import type { IconNode } from 'lucide'
import { ArrowUpRight, ChevronLeft, ChevronRight, Heart } from 'lucide'
import {
  DEFAULT_PRODUCT_SIZE,
  PRODUCT_COPY,
  PRODUCT_SIZE_OPTIONS,
  getProductRoute,
} from '~/constants/products'
import type { Product, ProductMediaItem } from '~/types/product'
import type { ProductSize } from '~/types/product'
import { formatProductPrice } from '~/utils/product'

interface ProductCardProps {
  product: Product
}

const props = defineProps<ProductCardProps>()
const activeMediaIndex = ref(0)
const { mediaItems } = useProductMedia(() => props.product, true)
const { add, open: openCart } = useCart()
const { has: hasWishlistedProduct, toggle: toggleWishlistedProduct } = useWishlist()
const { preloadImages } = useImagePreloader()
const selectedSize = ref<ProductSize>(DEFAULT_PRODUCT_SIZE)

const CARD_ICONS = {
  next: ChevronRight,
  previous: ChevronLeft,
  view: ArrowUpRight,
  wishlist: Heart,
} as const satisfies Readonly<Record<string, IconNode>>

// Returns the media item selected in this card carousel.
const activeMediaItem = computed<ProductMediaItem>(
  () => mediaItems.value[activeMediaIndex.value]!,
)
// Returns the canonical route for this product.
const productRoute = computed(() => getProductRoute(props.product.slug))
// Formats this product's integer kobo price for display.
const formattedPrice = computed(() => formatProductPrice(props.product.priceKobo))
// Reports whether the current product has more than one media frame.
const hasMultipleMediaItems = computed(() => mediaItems.value.length > 1)
const mediaSources = computed(() =>
  mediaItems.value.map((mediaItem) => mediaItem.source),
)
const adjacentMediaSources = computed(() => {
  if (!hasMultipleMediaItems.value) {
    return []
  }

  const previousMediaIndex = (
    activeMediaIndex.value - 1 + mediaItems.value.length
  ) % mediaItems.value.length
  const nextMediaIndex = (activeMediaIndex.value + 1) % mediaItems.value.length

  return [...new Set([previousMediaIndex, nextMediaIndex])]
    .map((mediaIndex) => mediaItems.value[mediaIndex]?.source)
})
// Reports whether this card's product is saved in the persisted guest wishlist.
const isWishlisted = computed(() => hasWishlistedProduct(props.product.slug))
// Maps wishlist state to the correct accessible action label.
const wishlistActionLabel = computed(() =>
  isWishlisted.value
    ? PRODUCT_COPY.actions.wishlistRemove
    : PRODUCT_COPY.actions.wishlistAdd,
)

// Selects one known media position in the card carousel.
const selectMediaItem = (mediaIndex: number): void => {
  activeMediaIndex.value = mediaIndex
}

// Advances the card carousel and wraps at the final image.
const showNextMediaItem = (): void => {
  activeMediaIndex.value = (activeMediaIndex.value + 1) % mediaItems.value.length
}

// Rewinds the card carousel and wraps at the first image.
const showPreviousMediaItem = (): void => {
  activeMediaIndex.value = (
    activeMediaIndex.value - 1 + mediaItems.value.length
  ) % mediaItems.value.length
}

// Toggles this product in the persisted guest wishlist.
const toggleWishlist = (): void => {
  toggleWishlistedProduct(props.product)
}

const preloadAdjacentMediaImages = (): void => {
  preloadImages(adjacentMediaSources.value)
}

const preloadAllMediaImages = (): void => {
  preloadImages(mediaSources.value)
}

// Adds the product to the bag. A sized product cannot be added from the grid
// without a size, so the customer is sent to the detail page to choose one.
const addProductToBag = (): void => {
  const size = props.product.hasSizes ? selectedSize.value : null

  if (add(props.product, size)) {
    openCart()
  }
}

// Builds an accessible label for one carousel position.
const getMediaButtonLabel = (mediaIndex: number): string =>
  `${PRODUCT_COPY.actions.showImage} ${mediaIndex + 1}`

const isMediaItemActive = (mediaIndex: number): boolean =>
  mediaIndex === activeMediaIndex.value

const isSizeSelected = (size: ProductSize): boolean =>
  selectedSize.value === size

const selectSize = (size: ProductSize): void => {
  selectedSize.value = size
}

watch(adjacentMediaSources, preloadAdjacentMediaImages, {
  flush: 'post',
  immediate: true,
})
</script>

<template>
  <article
    class="product-card"
    :data-product-slug="product.slug"
    @focusin="preloadAllMediaImages"
    @mouseenter="preloadAllMediaImages"
  >
    <div class="product-card-media">
      <CatalogueProductMedia :media-item="activeMediaItem" />

      <span v-if="product.tag" class="product-tag">{{ product.tag }}</span>

      <button
        class="product-wishlist-button"
        :class="{ 'is-active': isWishlisted }"
        type="button"
        :aria-label="wishlistActionLabel"
        :title="wishlistActionLabel"
        @click="toggleWishlist"
      >
        <StorefrontIcon
          :icon-node="CARD_ICONS.wishlist"
          :size="18"
          :stroke-width="1.5"
          aria-hidden="true"
        />
      </button>

      <template v-if="hasMultipleMediaItems">
        <button
          class="product-carousel-button is-previous"
          type="button"
          :aria-label="PRODUCT_COPY.actions.previousImage"
          :title="PRODUCT_COPY.actions.previousImage"
          @click="showPreviousMediaItem"
        >
          <StorefrontIcon
            :icon-node="CARD_ICONS.previous"
            :size="16"
            :stroke-width="1.5"
            aria-hidden="true"
          />
        </button>
        <button
          class="product-carousel-button is-next"
          type="button"
          :aria-label="PRODUCT_COPY.actions.nextImage"
          :title="PRODUCT_COPY.actions.nextImage"
          @click="showNextMediaItem"
        >
          <StorefrontIcon
            :icon-node="CARD_ICONS.next"
            :size="16"
            :stroke-width="1.5"
            aria-hidden="true"
          />
        </button>

        <div class="product-carousel-dots">
          <button
            v-for="(mediaItem, mediaIndex) in mediaItems"
            :key="mediaItem.id"
            type="button"
            :class="{ 'is-active': isMediaItemActive(mediaIndex) }"
            :aria-label="getMediaButtonLabel(mediaIndex)"
            @click="selectMediaItem(mediaIndex)"
          />
        </div>
      </template>
    </div>

    <div class="product-card-content">
      <p v-if="!product.hasSizes" class="product-category">{{ product.category }}</p>
      <h3>
        <NuxtLink :to="productRoute">{{ product.name }}</NuxtLink>
      </h3>
      <p class="product-price">{{ formattedPrice }}</p>

      <div
        v-if="product.hasSizes"
        class="product-card-sizes"
        role="group"
        :aria-label="PRODUCT_COPY.size.label"
      >
        <button
          v-for="size in PRODUCT_SIZE_OPTIONS"
          :key="size"
          type="button"
          :class="{ 'is-active': isSizeSelected(size) }"
          :aria-pressed="isSizeSelected(size)"
          @click="selectSize(size)"
        >
          {{ size }}
        </button>
      </div>

      <div class="product-card-actions">
        <NuxtLink class="product-view-link" :to="productRoute">
          <span>{{ PRODUCT_COPY.actions.viewProduct }}</span>
          <StorefrontIcon
            :icon-node="CARD_ICONS.view"
            :size="15"
            :stroke-width="1.5"
            aria-hidden="true"
          />
        </NuxtLink>
        <button class="product-add-button" type="button" @click="addProductToBag">
          <CatalogueAddToBagButtonContent :icon-size="15" />
        </button>
      </div>
    </div>
  </article>
</template>
