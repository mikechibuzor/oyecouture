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
const { mediaItems } = useProductMedia(() => props.product)
const { preloadImages } = useImagePreloader()

const GALLERY_ICONS = {
  next: ArrowRight,
  previous: ArrowLeft,
} as const satisfies Readonly<Record<string, IconNode>>

// Returns the currently selected gallery media item.
const activeMediaItem = computed<ProductMediaItem>(
  () => mediaItems.value[activeMediaIndex.value]!,
)
const productMediaTransitionStyle = computed(() => ({
  viewTransitionName: `product-${props.product.slug}`,
}))
const hasMultipleMediaItems = computed(() => mediaItems.value.length > 1)
const adjacentMediaSources = computed(() => {
  if (mediaItems.value.length <= 1) {
    return []
  }

  const previousIndex = (
    activeMediaIndex.value - 1 + mediaItems.value.length
  ) % mediaItems.value.length
  const nextIndex = (activeMediaIndex.value + 1) % mediaItems.value.length

  return [...new Set([previousIndex, nextIndex])]
    .map((mediaIndex) => mediaItems.value[mediaIndex]?.source)
})

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

const isMediaItemActive = (mediaIndex: number): boolean =>
  mediaIndex === activeMediaIndex.value

watch(adjacentMediaSources, (sources) => preloadImages(sources), {
  flush: 'post',
  immediate: true,
})
</script>

<template>
  <section class="product-gallery" :aria-label="PRODUCT_COPY.galleryLabel">
    <div class="product-gallery-stage">
      <Transition name="product-gallery-media">
        <CatalogueProductMedia
          :key="activeMediaItem.id"
          :media-item="activeMediaItem"
          :style="productMediaTransitionStyle"
          eager
        />
      </Transition>
      <span v-if="product.tag" class="product-tag">{{ product.tag }}</span>

      <button
        v-if="hasMultipleMediaItems"
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
        v-if="hasMultipleMediaItems"
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

    <div v-if="hasMultipleMediaItems" class="product-gallery-thumbnails">
      <button
        v-for="(mediaItem, mediaIndex) in mediaItems"
        :key="mediaItem.id"
        type="button"
        :class="{ 'is-active': isMediaItemActive(mediaIndex) }"
        :aria-label="getThumbnailLabel(mediaIndex)"
        :aria-pressed="isMediaItemActive(mediaIndex)"
        @click="selectMediaItem(mediaIndex)"
      >
        <CatalogueProductMedia :media-item="mediaItem" />
      </button>
    </div>
  </section>
</template>
