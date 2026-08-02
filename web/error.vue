<script setup lang="ts">
import type { NuxtError } from '#app'
import { StorefrontRoute } from '~/types/storefront'

interface ErrorPageProps {
  error: NuxtError
}

const props = defineProps<ErrorPageProps>()

const isNotFound = computed(() => props.error.statusCode === 404)
const heading = computed(() =>
  isNotFound.value ? 'This page has moved on' : 'Something went wrong',
)
const body = computed(() =>
  isNotFound.value
    ? 'The piece or page you were looking for is no longer here. The collection is still where you left it.'
    : 'We hit an unexpected problem. Please try again in a moment.',
)

// Clears the error before navigating, otherwise the error page persists.
const returnHome = () => clearError({ redirect: StorefrontRoute.Home })
</script>

<template>
  <div class="error-page">
    <p class="confirm-eyebrow">{{ error.statusCode }}</p>
    <h1>{{ heading }}</h1>
    <p class="error-body">{{ body }}</p>
    <button class="cart-primary-action error-action" type="button" @click="returnHome">
      Back to the collection
    </button>
  </div>
</template>

<style>
/* The error page renders outside the default layout, so it carries its own
   minimal styling rather than depending on the storefront shell. */
.error-page {
  display: grid;
  min-height: 100vh;
  padding: var(--space-2xl);
  align-content: center;
  justify-items: center;
  background: var(--color-canvas);
  text-align: center;
}

.error-page h1 {
  margin: 0 0 var(--space-md);
  font-family: var(--font-family-display);
  font-size: var(--text-3xl);
  font-weight: var(--font-weight-bold);
  letter-spacing: var(--tracking-tight);
}

.error-body {
  max-width: 30rem;
  margin: 0 0 var(--space-2xl);
  color: var(--color-body-copy);
  font-size: var(--text-md);
  font-weight: var(--font-weight-light);
  line-height: 1.7;
}

.error-action {
  width: auto;
  padding: 0 var(--space-2xl);
}
</style>
