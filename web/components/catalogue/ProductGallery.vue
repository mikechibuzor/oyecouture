<script setup lang="ts">
import type { IconNode } from 'lucide'
import { ArrowLeft, ArrowRight } from 'lucide'
import { PRODUCT_COPY } from '~/constants/products'
import type { Product, ProductMediaItem } from '~/types/product'

interface ProductGalleryProps {
  product: Product
}

const props = defineProps<ProductGalleryProps>()
const activeMediaIndex = ref(0)
const { mediaItems } = useProductMedia(() => props.product, true)

const GALLERY_ICONS = {
  next: ArrowRight,
  previous: ArrowLeft,
} as const satisfies Readonly<Record<string, IconNode>>

// Returns the currently selected gallery media item.
const activeMediaItem = computed<ProductMediaItem>(
  () => mediaItems.value[activeMediaIndex.value]!,
)

// Selects one known thumbnail in the product gallery.
const selectMediaItem = (mediaIndex: number): void => {
  activeMediaIndex.value = mediaIndex
}

// Advances the gallery and wraps at the final media item.
const showNextMediaItem = (): void => {
  activeMediaIndex.value = (activeMediaIndex.value + 1) % mediaItems.value.length
}

// Rewinds the gallery and wraps at the first media item.
const showPreviousMediaItem = (): void => {
  activeMediaIndex.value = (
    activeMediaIndex.value - 1 + mediaItems.value.length
  ) % mediaItems.value.length
}

// Builds an accessible label for one gallery thumbnail.
const getThumbnailLabel = (mediaIndex: number): string =>
  `${PRODUCT_COPY.actions.showImage} ${mediaIndex + 1}`
</script>

<template>
  <section class="product-gallery" :aria-label="PRODUCT_COPY.galleryLabel">
    <div class="product-gallery-stage">
      <CatalogueProductMedia :media-item="activeMediaItem" eager />
      <span v-if="product.tag" class="product-tag">{{ product.tag }}</span>

      <button
        class="product-gallery-arrow is-previous"
        type="button"
        :aria-label="PRODUCT_COPY.actions.previousImage"
        :title="PRODUCT_COPY.actions.previousImage"
        @click="showPreviousMediaItem"
      >
        <StorefrontIcon
          :icon-node="GALLERY_ICONS.previous"
          :size="20"
          :stroke-width="1.5"
          aria-hidden="true"
        />
      </button>
      <button
        class="product-gallery-arrow is-next"
        type="button"
        :aria-label="PRODUCT_COPY.actions.nextImage"
        :title="PRODUCT_COPY.actions.nextImage"
        @click="showNextMediaItem"
      >
        <StorefrontIcon
          :icon-node="GALLERY_ICONS.next"
          :size="20"
          :stroke-width="1.5"
          aria-hidden="true"
        />
      </button>
    </div>

    <div class="product-gallery-thumbnails">
      <button
        v-for="(mediaItem, mediaIndex) in mediaItems"
        :key="mediaItem.id"
        type="button"
        :class="{ 'is-active': mediaIndex === activeMediaIndex }"
        :aria-label="getThumbnailLabel(mediaIndex)"
        :aria-pressed="mediaIndex === activeMediaIndex"
        @click="selectMediaItem(mediaIndex)"
      >
        <CatalogueProductMedia :media-item="mediaItem" />
      </button>
    </div>
  </section>
</template>
