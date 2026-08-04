<script setup lang="ts">
import { ABOUT_BESPOKE_BAND } from '~/constants/about'
import { StorefrontActionId, StorefrontRoute } from '~/types/storefront'

const isSearchOverlayOpen = ref(false)
const { open: openCart } = useCart()
const route = useRoute()
const router = useRouter()

// The About page runs its own wording for the bespoke band; every other page
// falls back to the component's shared copy.
const bespokeBandProps = computed(() =>
  route.path === StorefrontRoute.About ? ABOUT_BESPOKE_BAND : {},
)

// Routes implemented header actions to their owning storefront surfaces.
const handleStorefrontAction = (storefrontActionId: StorefrontActionId): void => {
  if (storefrontActionId === StorefrontActionId.Search) {
    isSearchOverlayOpen.value = true
    return
  }

  if (storefrontActionId === StorefrontActionId.Cart) {
    openCart()
    return
  }

  if (storefrontActionId === StorefrontActionId.Wishlist) {
    void router.push(StorefrontRoute.Wishlist)
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
    <StorefrontBespokeServiceBand v-bind="bespokeBandProps" />
    <StorefrontFooter />
    <CartDrawer />
    <StorefrontSearchOverlay
      :is-open="isSearchOverlayOpen"
      @close="closeSearchOverlay"
    />
  </div>
</template>
