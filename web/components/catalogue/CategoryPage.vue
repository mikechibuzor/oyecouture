<script setup lang="ts">
import { PRODUCT_COPY } from '~/constants/products'
import type { CategoryPageConfig, Product } from '~/types/product'
import { getProductErrorMessage } from '~/utils/product'

interface CategoryPageProps {
  config: CategoryPageConfig
}

const props = defineProps<CategoryPageProps>()
const { data, error, status } = await useProductList(props.config.section)

// Returns the section products while preserving a stable empty array during loading.
const products = computed<readonly Product[]>(() => data.value?.products ?? [])
// Maps a structured API failure to safe customer-facing copy.
const errorMessage = computed(() =>
  getProductErrorMessage(error.value, PRODUCT_COPY.catalogue.error),
)
// Builds a unique collection heading for this category.
const collectionHeading = computed(() =>
  `${props.config.heading} ${PRODUCT_COPY.categoryCollectionSuffix}`,
)
</script>

<template>
  <div class="category-page">
    <header class="category-page-header">
      <div>
        <p class="eyebrow">{{ config.eyebrow }}</p>
        <h1>{{ config.heading }}</h1>
        <p>{{ config.description }}</p>
      </div>
    </header>

    <CatalogueCollection
      :products="products"
      :status="status"
      :error-message="errorMessage"
      :show-filters="false"
      :heading="collectionHeading"
    />
  </div>
</template>
