<script setup lang="ts">
import { PRODUCT_COPY } from '~/constants/products'
import type { ProductMediaItem } from '~/types/product'

interface ProductMediaProps {
  mediaItem: ProductMediaItem
  eager?: boolean
}

withDefaults(defineProps<ProductMediaProps>(), {
  eager: false,
})
</script>

<template>
  <div class="product-media">
    <img
      v-if="mediaItem.source"
      :src="mediaItem.source"
      :alt="mediaItem.alt"
      :loading="eager ? 'eager' : 'lazy'"
      decoding="async"
    >
    <div v-else class="product-media-placeholder" role="img" :aria-label="mediaItem.alt">
      <span>{{ mediaItem.alt }}</span>
      <small>{{ PRODUCT_COPY.imagePlaceholder }}</small>
    </div>
  </div>
</template>
