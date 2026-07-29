<script setup lang="ts">
import type { IconNode } from 'lucide'
import { Image as ImageIcon } from 'lucide'
import type { ProductMediaItem } from '~/types/product'

interface ProductMediaProps {
  mediaItem: ProductMediaItem
  eager?: boolean
}

withDefaults(defineProps<ProductMediaProps>(), {
  eager: false,
})

const PRODUCT_MEDIA_PLACEHOLDER_ICON: IconNode = ImageIcon
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
      <StorefrontIcon
        :icon-node="PRODUCT_MEDIA_PLACEHOLDER_ICON"
        :size="28"
        :stroke-width="1.4"
        aria-hidden="true"
      />
      <span>{{ mediaItem.alt }}</span>
    </div>
  </div>
</template>
