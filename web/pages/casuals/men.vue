<script setup lang="ts">
import { CATEGORY_PAGE_BY_SECTION } from '~/constants/products'
import { SEO_COPY } from '~/constants/seo'
import { ProductSection } from '~/types/product'

const categoryConfig = CATEGORY_PAGE_BY_SECTION[ProductSection.CasualsMen]

useStorefrontSeo(SEO_COPY.casualsMen)
</script>

<template>
  <CatalogueCategoryPage :config="categoryConfig" />
</template>
