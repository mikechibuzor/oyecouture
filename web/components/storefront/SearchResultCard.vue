<script setup lang="ts">
import { getProductRoute } from '~/constants/products'
import type { Product, ProductMediaItem } from '~/types/product'
import { formatProductPrice } from '~/utils/product'

interface SearchResultCardProps {
  product: Product
}

const props = defineProps<SearchResultCardProps>()
const emit = defineEmits<{
  select: []
}>()
const { mediaItems } = useProductMedia(() => props.product)

// Returns the first available product photograph or its stable placeholder.
const primaryMediaItem = computed<ProductMediaItem>(() => mediaItems.value[0]!)
// Returns the canonical detail route for the selected search result.
const productRoute = computed(() => getProductRoute(props.product.slug))
// Formats the integer kobo value for customer-facing display.
const formattedPrice = computed(() => formatProductPrice(props.product.priceKobo))

// Notifies the overlay to close after a result starts navigation.
const handleProductSelection = (): void => {
  emit('select')
}
</script>

<template>
  <NuxtLink
    class="search-result-card"
    :to="productRoute"
    @click="handleProductSelection"
  >
    <div class="search-result-media">
      <CatalogueProductMedia :media-item="primaryMediaItem" />
    </div>
    <div class="search-result-content">
      <h3>{{ product.name }}</h3>
      <p class="search-result-category">{{ product.category }}</p>
      <p class="search-result-price">{{ formattedPrice }}</p>
    </div>
  </NuxtLink>
</template>
