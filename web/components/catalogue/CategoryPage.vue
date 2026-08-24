<script setup lang="ts">
import {
  CATEGORY_ALL_FILTER,
  CATEGORY_ALL_FILTER_OPTION,
  PRODUCT_COPY,
} from '~/constants/products'
import type { CategoryPageConfig, Product } from '~/types/product'
import { StorefrontRoute } from '~/types/storefront'
import { getProductErrorMessage } from '~/utils/product'

interface CategoryPageProps {
  config: CategoryPageConfig
}

const props = defineProps<CategoryPageProps>()
const { data, error, status } = await useProductList(props.config.section)
const activeCategory = ref<string>(CATEGORY_ALL_FILTER)
const categoryOptions = computed(() => [CATEGORY_ALL_FILTER_OPTION, ...props.config.categories])

// Resets the chip when the customer moves between category routes.
watch(() => props.config.section, () => {
  activeCategory.value = CATEGORY_ALL_FILTER
})

// Returns the section products while preserving a stable empty array during loading.
const products = computed<readonly Product[]>(() => data.value?.products ?? [])
// Narrows the section to one sub-category when a chip other than ALL is active.
const visibleProducts = computed(() =>
  activeCategory.value === CATEGORY_ALL_FILTER
    ? products.value
    : products.value.filter((product) => product.category === activeCategory.value),
)
// Formats the visible result count with the correct noun.
const pieceCountLabel = computed(() => {
  const productCount = visibleProducts.value.length
  const pieceLabel = productCount === 1
    ? PRODUCT_COPY.catalogue.pieceSingular
    : PRODUCT_COPY.catalogue.piecePlural

  return `${productCount} ${pieceLabel}`
})
// Maps a structured API failure to safe customer-facing copy.
const errorMessage = computed(() =>
  getProductErrorMessage(error.value, PRODUCT_COPY.catalogue.error),
)

const isCategoryActive = (category: string): boolean =>
  activeCategory.value === category

const selectCategory = (category: string): void => {
  activeCategory.value = category
}
</script>

<template>
  <div class="category-page">
    <header class="category-page-header">
      <div class="category-page-header-inner">
        <nav class="category-breadcrumbs" :aria-label="PRODUCT_COPY.breadcrumbs.label">
          <NuxtLink :to="StorefrontRoute.Home">{{ PRODUCT_COPY.breadcrumbs.home }}</NuxtLink>
          <span aria-hidden="true">/</span>
          <span>{{ config.breadcrumb }}</span>
        </nav>

        <h1>{{ config.heading }}</h1>
        <p>{{ config.description }}</p>
      </div>
    </header>

    <section class="category-page-body">
      <div class="category-toolbar">
        <div
          v-if="config.categories.length > 0"
          class="catalogue-filter-list"
          role="group"
          :aria-label="PRODUCT_COPY.catalogue.filterLabel"
        >
          <button
            v-for="category in categoryOptions"
            :key="category.value"
            type="button"
            :class="{ 'is-active': isCategoryActive(category.value) }"
            :aria-pressed="isCategoryActive(category.value)"
            @click="selectCategory(category.value)"
          >
            {{ category.label }}
          </button>
        </div>

        <span class="category-piece-count">{{ pieceCountLabel }}</span>
      </div>

      <CatalogueProductGrid
        :products="visibleProducts"
        :status="status"
        :error-message="errorMessage"
      />
    </section>
  </div>
</template>
