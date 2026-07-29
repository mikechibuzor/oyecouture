<script setup lang="ts">
import type { IconNode } from 'lucide'
import { ArrowRight } from 'lucide'
import {
  CATEGORY_TILES,
  HOME_HERO_IMAGES,
  PRODUCT_COPY,
} from '~/constants/products'
import type { Product } from '~/types/product'
import { StorefrontRoute } from '~/types/storefront'
import { getProductErrorMessage } from '~/utils/product'

const { data, error, status } = await useProductList()
const HERO_ACTION_ICON: IconNode = ArrowRight
const heroCarouselImages = [...HOME_HERO_IMAGES, ...HOME_HERO_IMAGES] as const

// Returns all API products while preserving a stable empty array during loading.
const products = computed<readonly Product[]>(() => data.value?.products ?? [])
// Maps a structured API failure to safe customer-facing copy.
const errorMessage = computed(() =>
  getProductErrorMessage(error.value, PRODUCT_COPY.catalogue.error),
)
</script>

<template>
  <div class="home-page">
    <section class="home-hero" aria-labelledby="home-hero-title">
      <div class="home-hero-content">
        <h1 id="home-hero-title">
          {{ PRODUCT_COPY.home.heroHeading.lead }}<br>
          {{ PRODUCT_COPY.home.heroHeading.tail }}
        </h1>
        <p>{{ PRODUCT_COPY.home.heroDescription }}</p>
        <NuxtLink :to="StorefrontRoute.Collection">
          <span>{{ PRODUCT_COPY.home.heroAction }}</span>
          <StorefrontIcon
            :icon-node="HERO_ACTION_ICON"
            :size="17"
            :stroke-width="1.5"
            aria-hidden="true"
          />
        </NuxtLink>
      </div>

      <div class="home-hero-media" aria-hidden="true">
        <div class="home-hero-media-track">
          <img
            v-for="(heroImage, heroImageIndex) in heroCarouselImages"
            :key="`${heroImage.id}-${heroImageIndex}`"
            :class="heroImage.size"
            :src="heroImage.source"
            :alt="heroImage.alt"
          >
        </div>
      </div>
    </section>

    <section class="category-tiles-section" aria-labelledby="category-tiles-title">
      <header>
        <p class="eyebrow">{{ PRODUCT_COPY.home.categoryEyebrow }}</p>
        <h2 id="category-tiles-title">{{ PRODUCT_COPY.categoryTilesHeading }}</h2>
      </header>
      <div class="category-tile-grid">
        <CatalogueCategoryTile
          v-for="category in CATEGORY_TILES"
          :key="category.to"
          :category="category"
        />
      </div>
    </section>

    <CatalogueCollection
      :products="products"
      :status="status"
      :error-message="errorMessage"
    />
  </div>
</template>
