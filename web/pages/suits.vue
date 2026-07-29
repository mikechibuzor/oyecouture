<script setup lang="ts">
import { CATEGORY_PAGE_BY_SECTION } from '~/constants/products'
import { ProductSection } from '~/types/product'

const categoryConfig = CATEGORY_PAGE_BY_SECTION[ProductSection.Suits]
</script>

<template>
  <CatalogueCategoryPage :config="categoryConfig" />
</template>
