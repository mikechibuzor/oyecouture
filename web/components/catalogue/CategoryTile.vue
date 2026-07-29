<script setup lang="ts">
import type { IconNode } from 'lucide'
import { ArrowRight } from 'lucide'
import type { CategoryTileConfig } from '~/types/product'

interface CategoryTileProps {
  category: CategoryTileConfig
}

defineProps<CategoryTileProps>()

const CATEGORY_TILE_ICON: IconNode = ArrowRight
</script>

<template>
  <NuxtLink class="category-tile" :to="category.to">
    <img :src="category.imagePath" :alt="category.imageAlt" loading="lazy">
    <span class="category-tile-overlay">
      <small>{{ category.label }}</small>
      <strong>{{ category.heading }}</strong>
      <StorefrontIcon
        :icon-node="CATEGORY_TILE_ICON"
        :size="22"
        :stroke-width="1.5"
        aria-hidden="true"
      />
    </span>
  </NuxtLink>
</template>
