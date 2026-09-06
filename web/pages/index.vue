<script setup lang="ts">
import type { IconNode } from 'lucide'
import { ArrowRight, Pause, Play } from 'lucide'
import {
  HOME_HERO_IMAGES,
  PRODUCT_COPY,
} from '~/constants/products'
import { SEO_COPY } from '~/constants/seo'
import type { Product } from '~/types/product'
import { getProductErrorMessage } from '~/utils/product'
import { buildBusinessSchema, buildWebsiteSchema, jsonLdScript } from '~/utils/seo'

// The homepage title leads with the brand rather than appending it, so a search
// for the shop's own name matches the front of the title.
useStorefrontSeo({ ...SEO_COPY.home, bareTitle: true })

useHead({
  script: [jsonLdScript(buildBusinessSchema()), jsonLdScript(buildWebsiteSchema())],
})

const { data, error, status } = await useProductList()
const { data: homeContent } = await useHomeContent()
const HERO_ACTION_ICON: IconNode = ArrowRight
const HERO_MOTION_ICONS = { pause: Pause, play: Play } as const
const isHeroMotionPaused = ref(false)
const heroCarouselImages = [...HOME_HERO_IMAGES, ...HOME_HERO_IMAGES] as const
const fallbackHomeContent = defaultHomeContent()

// Returns all API products while preserving a stable empty array during loading.
const products = computed<readonly Product[]>(() => data.value?.products ?? [])
// Maps a structured API failure to safe customer-facing copy.
const errorMessage = computed(() =>
  getProductErrorMessage(error.value, PRODUCT_COPY.catalogue.error),
)
const heroContent = computed(() => (homeContent.value ?? fallbackHomeContent).hero)
const categoryContent = computed(() => (homeContent.value ?? fallbackHomeContent).categories)
</script>

<template>
  <div class="home-page">
    <section class="home-hero" aria-labelledby="home-hero-title">
      <div class="home-hero-content">
        <h1 id="home-hero-title">
          {{ heroContent.headingLead }}<br>
          {{ heroContent.headingTail }}
        </h1>
        <p>{{ heroContent.description }}</p>
        <NuxtLink class="home-hero-action" :to="heroContent.actionTo">
          <span>{{ heroContent.actionLabel }}</span>
          <StorefrontIcon
            :icon-node="HERO_ACTION_ICON"
            :size="17"
            :stroke-width="1.5"
            aria-hidden="true"
          />
        </NuxtLink>
      </div>

      <div class="home-hero-media-shell" :class="{ 'is-paused': isHeroMotionPaused }">
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

        <button
          type="button"
          class="home-hero-motion-toggle"
          :aria-pressed="isHeroMotionPaused"
          :aria-label="isHeroMotionPaused ? 'Play image movement' : 'Pause image movement'"
          @click="isHeroMotionPaused = !isHeroMotionPaused"
        >
          <StorefrontIcon
            :icon-node="isHeroMotionPaused ? HERO_MOTION_ICONS.play : HERO_MOTION_ICONS.pause"
            :size="15"
            :stroke-width="1.6"
            aria-hidden="true"
          />
        </button>
      </div>
    </section>

    <section class="category-tiles-section" aria-labelledby="category-tiles-title">
      <header>
        <p class="eyebrow">{{ categoryContent.eyebrow }}</p>
        <h2 id="category-tiles-title">{{ categoryContent.heading }}</h2>
      </header>
      <div class="category-tile-grid">
        <CatalogueCategoryTile
          v-for="category in categoryContent.tiles"
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
