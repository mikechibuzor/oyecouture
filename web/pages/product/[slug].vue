<script setup lang="ts">
import {
  DEFAULT_PRODUCT_SIZE,
  PRODUCT_COPY,
  PRODUCT_SECTION_NAVIGATION,
} from '~/constants/products'
import type { Product } from '~/types/product'
import { AsyncDataStatus, ProductSize } from '~/types/product'
import { StorefrontRoute } from '~/types/storefront'
import { formatProductPrice, getProductErrorMessage } from '~/utils/product'

const route = useRoute()

// Normalizes the dynamic route parameter into one product slug.
const productSlug = computed(() => {
  const routeSlug = route.params.slug
  return Array.isArray(routeSlug) ? (routeSlug[0] ?? '') : (routeSlug ?? '')
})

const { data, error, status } = await useProductBySlug(productSlug.value)
const { add, open: openCart } = useCart()
const selectedSize = ref<ProductSize>(DEFAULT_PRODUCT_SIZE)

// Adds the loaded product to the bag, passing the chosen size for ready-to-wear.
const addProductToBag = (): void => {
  if (!product.value) {
    return
  }

  if (add(product.value, product.value.hasSizes ? selectedSize.value : null)) {
    openCart()
  }
}

// Returns the loaded product when the detail request succeeds.
const product = computed<Product | null>(() => data.value?.product ?? null)
// Reports whether the product detail is still loading.
const isLoading = computed(() =>
  status.value === AsyncDataStatus.Idle
  || status.value === AsyncDataStatus.Pending,
)
// Reports whether the product detail request failed.
const hasError = computed(() => status.value === AsyncDataStatus.Error)
// Maps a structured detail failure to safe customer-facing copy.
const errorMessage = computed(() =>
  getProductErrorMessage(error.value, PRODUCT_COPY.productDetail.error),
)
// Resolves the loaded product's section breadcrumb.
const sectionNavigation = computed(() =>
  product.value ? PRODUCT_SECTION_NAVIGATION[product.value.section] : null,
)
// Formats the loaded product's integer kobo price.
const formattedPrice = computed(() =>
  product.value ? formatProductPrice(product.value.priceKobo) : '',
)
</script>

<template>
  <div class="product-detail-page">
    <p v-if="isLoading" class="product-detail-state" role="status">
      {{ PRODUCT_COPY.productDetail.loading }}
    </p>
    <p v-else-if="hasError" class="product-detail-state is-error" role="alert">
      {{ errorMessage }}
    </p>

    <template v-else-if="product && sectionNavigation">
      <nav class="product-breadcrumbs" :aria-label="PRODUCT_COPY.breadcrumbs.label">
        <NuxtLink :to="StorefrontRoute.Home">
          {{ PRODUCT_COPY.breadcrumbs.home }}
        </NuxtLink>
        <span aria-hidden="true">/</span>
        <NuxtLink :to="sectionNavigation.to">{{ sectionNavigation.label }}</NuxtLink>
        <span aria-hidden="true">/</span>
        <span>{{ product.name }}</span>
      </nav>

      <div class="product-detail-layout">
        <CatalogueProductGallery :product="product" />

        <section class="product-detail-information">
          <p class="product-category">{{ product.category }}</p>
          <h1>{{ product.name }}</h1>
          <p class="product-detail-price">{{ formattedPrice }}</p>
          <p class="product-detail-description">{{ product.description }}</p>

          <CatalogueSizeSelector
            v-if="product.hasSizes"
            v-model="selectedSize"
          />

          <button class="product-detail-add-button" type="button" @click="addProductToBag">
            <CatalogueAddToBagButtonContent />
          </button>

          <CatalogueBespokeNote v-if="!product.hasSizes" />
        </section>
      </div>
    </template>

    <p v-else class="product-detail-state" role="status">
      {{ PRODUCT_COPY.productDetail.notFound }}
    </p>
  </div>
</template>
