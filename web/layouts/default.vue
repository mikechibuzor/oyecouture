<script setup lang="ts">
import { StorefrontActionId } from '~/types/storefront'

const isSearchOverlayOpen = ref(false)

// Routes implemented header actions while leaving future account, wishlist, and cart work inert.
const handleStorefrontAction = (storefrontActionId: StorefrontActionId): void => {
  if (storefrontActionId === StorefrontActionId.Search) {
    isSearchOverlayOpen.value = true
  }
}

// Dismisses the search overlay from every supported close path.
const closeSearchOverlay = (): void => {
  isSearchOverlayOpen.value = false
}
</script>

<template>
  <div class="storefront-shell">
    <StorefrontAnnouncementBar />
    <StorefrontNavBar @action="handleStorefrontAction" />
    <main class="storefront-main">
      <slot />
    </main>
    <StorefrontFooter />
    <StorefrontSearchOverlay
      :is-open="isSearchOverlayOpen"
      @close="closeSearchOverlay"
    />
  </div>
</template>
