<script setup lang="ts">
import { PRODUCT_COPY, PRODUCT_SIZE_OPTIONS } from '~/constants/products'
import { ProductSize } from '~/types/product'

const selectedSize = defineModel<ProductSize>({ required: true })

// Selects one available ready-to-wear size.
const selectSize = (productSize: ProductSize): void => {
  selectedSize.value = productSize
}
</script>

<template>
  <fieldset class="product-size-selector">
    <legend>{{ PRODUCT_COPY.size.label }}</legend>
    <div>
      <button
        v-for="productSize in PRODUCT_SIZE_OPTIONS"
        :key="productSize"
        type="button"
        :class="{ 'is-active': selectedSize === productSize }"
        :aria-pressed="selectedSize === productSize"
        @click="selectSize(productSize)"
      >
        {{ productSize }}
      </button>
    </div>
  </fieldset>
</template>
