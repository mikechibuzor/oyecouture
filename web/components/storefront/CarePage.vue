<script setup lang="ts">
import { PRODUCT_COPY } from '~/constants/products'
import { StorefrontRoute } from '~/types/storefront'

interface CareSection {
  readonly title: string
  readonly body: string
}

interface CarePageProps {
  breadcrumb: string
  heading: string
  intro: string
  sections: readonly CareSection[]
  // Closing line beneath the sections, where the copy has one.
  footnote?: string
}

const props = withDefaults(defineProps<CarePageProps>(), {
  footnote: '',
})

// The care pages set their own search metadata here rather than page by page,
// since every one of them is this component with different copy. Going through
// the shared composable adds the canonical and link-preview tags a bare
// useSeoMeta leaves out.
useStorefrontSeo({ title: props.heading, description: props.intro })
</script>

<template>
  <div class="care-page">
    <header class="category-page-header">
      <div class="category-page-header-inner">
        <nav class="category-breadcrumbs" :aria-label="PRODUCT_COPY.breadcrumbs.label">
          <NuxtLink :to="StorefrontRoute.Home">{{ PRODUCT_COPY.breadcrumbs.home }}</NuxtLink>
          <span aria-hidden="true">/</span>
          <span>{{ breadcrumb }}</span>
        </nav>
        <h1>{{ heading }}</h1>
        <p>{{ intro }}</p>
      </div>
    </header>

    <section class="care-body">
      <article v-for="section in sections" :key="section.title" class="care-section">
        <h2>{{ section.title }}</h2>
        <p>{{ section.body }}</p>
      </article>

      <p v-if="footnote" class="care-footnote">{{ footnote }}</p>

    </section>
  </div>
</template>
