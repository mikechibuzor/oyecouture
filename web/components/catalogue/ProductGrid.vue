<script setup lang="ts">
import { PRODUCT_COPY } from '~/constants/products'
import type { AsyncDataStatusValue, Product } from '~/types/product'
import { AsyncDataStatus } from '~/types/product'

interface ProductGridProps {
  products: readonly Product[]
  status: AsyncDataStatusValue
  errorMessage: string
}

const props = defineProps<ProductGridProps>()

// Reports whether product data is currently loading.
const isLoading = computed(() =>
  props.status === AsyncDataStatus.Idle
  || props.status === AsyncDataStatus.Pending,
)
// Reports whether the product request failed.
const hasError = computed(() => props.status === AsyncDataStatus.Error)
// Reports whether a completed request has no visible products.
const isEmpty = computed(() =>
  props.status === AsyncDataStatus.Success && props.products.length === 0,
)
</script>

<template>
  <div class="product-grid-region" aria-live="polite">
    <template v-if="isLoading">
      <p class="visually-hidden" role="status">{{ PRODUCT_COPY.catalogue.loading }}</p>
      <div class="product-grid is-loading" aria-hidden="true">
        <article v-for="skeletonIndex in 4" :key="skeletonIndex" class="product-card-skeleton">
          <span class="product-card-skeleton-media" />
          <span class="product-card-skeleton-line is-short" />
          <span class="product-card-skeleton-line" />
          <span class="product-card-skeleton-line is-price" />
        </article>
      </div>
    </template>
    <p v-else-if="hasError" class="catalogue-state is-error" role="alert">
      {{ errorMessage }}
    </p>
    <p v-else-if="isEmpty" class="catalogue-state" role="status">
      {{ PRODUCT_COPY.catalogue.empty }}
    </p>
    <div v-else class="product-grid">
      <CatalogueProductCard
        v-for="product in products"
        :key="product.id"
        :product="product"
      />
    </div>
  </div>
</template>
