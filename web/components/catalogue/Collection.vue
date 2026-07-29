<script setup lang="ts">
import { PRODUCT_COPY } from '~/constants/products'
import type { AsyncDataStatusValue, Product } from '~/types/product'

interface CatalogueCollectionProps {
  products: readonly Product[]
  status: AsyncDataStatusValue
  errorMessage: string
  showFilters?: boolean
  eyebrow?: string
  heading?: string
}

const props = withDefaults(defineProps<CatalogueCollectionProps>(), {
  showFilters: true,
  eyebrow: PRODUCT_COPY.catalogue.eyebrow,
  heading: PRODUCT_COPY.catalogue.heading,
})

const {
  activeFilterId,
  activeSortId,
  pieceCountLabel,
  visibleProducts,
} = useCatalogueProducts(() => props.products)
</script>

<template>
  <section id="collection" class="catalogue-collection" :aria-labelledby="`${heading}-title`">
    <header class="catalogue-collection-heading">
      <p class="eyebrow">{{ eyebrow }}</p>
      <h2 :id="`${heading}-title`">{{ heading }}</h2>
    </header>

    <CatalogueToolbar
      v-model:active-filter-id="activeFilterId"
      v-model:active-sort-id="activeSortId"
      :piece-count-label="pieceCountLabel"
      :show-filters="showFilters"
    />
    <CatalogueProductGrid
      :products="visibleProducts"
      :status="status"
      :error-message="errorMessage"
    />
  </section>
</template>
