<script setup lang="ts">
import { BESPOKE_SERVICE_COPY } from '~/constants/storefront'
import { StorefrontRoute } from '~/types/storefront'

// The About page runs its own heading and adds a line of copy; every other page
// uses the shared wording.
interface BespokeServiceBandProps {
  heading?: string
  description?: string
}

withDefaults(defineProps<BespokeServiceBandProps>(), {
  heading: BESPOKE_SERVICE_COPY.heading,
  description: '',
})
</script>

<template>
  <section class="bespoke-service-band" aria-labelledby="bespoke-service-title">
    <!-- The Oyecouture symbol, taken from the approved design. -->
    <svg
      class="bespoke-service-mark"
      width="36"
      height="36"
      viewBox="0 0 100 100"
      fill="none"
      aria-hidden="true"
    >
      <circle cx="50" cy="50" r="45" stroke="currentColor" stroke-width="4" />
      <path d="M50 50 L50 5 A45 45 0 0 1 95 50 Z" fill="currentColor" />
      <path d="M50 50 L50 95" stroke="currentColor" stroke-width="4" />
      <path d="M50 50 L5 50" stroke="currentColor" stroke-width="4" />
      <path d="M50 50 L16 82" stroke="currentColor" stroke-width="4" />
    </svg>

    <p class="eyebrow">{{ BESPOKE_SERVICE_COPY.eyebrow }}</p>
    <h2 id="bespoke-service-title">{{ BESPOKE_SERVICE_COPY.heading }}</h2>

    <NuxtLink class="bespoke-service-action" :to="StorefrontRoute.BookFitting">
      {{ BESPOKE_SERVICE_COPY.action }}
    </NuxtLink>
  </section>
</template>
