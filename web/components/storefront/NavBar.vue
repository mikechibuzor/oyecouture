<script setup lang="ts">
import type { IconNode } from 'lucide'
import {
  ChevronDown,
  Heart,
  Search,
  ShoppingBag,
  UserRound,
} from 'lucide'
import {
  DESKTOP_HEADER_ACTIONS,
  MOBILE_HEADER_ACTIONS,
  MOBILE_MENU_ACTIONS,
  STOREFRONT_ACTION_BY_ID,
  STOREFRONT_COPY,
  STOREFRONT_NAVIGATION,
  isInternalStorefrontRoute,
} from '~/constants/storefront'
import type { StorefrontNavigationItem } from '~/types/storefront'
import { StorefrontActionId } from '~/types/storefront'

const emit = defineEmits<{
  action: [storefrontActionId: StorefrontActionId]
}>()

const navigationRoot = ref<HTMLElement | null>(null)
const {
  closeDropdown,
  closeNavigation,
  isDropdownOpen,
  isMobileMenuOpen,
  isNavigationItemActive,
  openDropdown,
  toggleDropdown,
  toggleMobileMenu,
} = useStorefrontNavigation(navigationRoot)
const { count: cartCount } = useCart()
const { count: wishlistCount } = useWishlist()

const STOREFRONT_ACTION_ICONS: Readonly<Record<StorefrontActionId, IconNode>> = {
  [StorefrontActionId.Search]: Search,
  [StorefrontActionId.Account]: UserRound,
  [StorefrontActionId.Wishlist]: Heart,
  [StorefrontActionId.Cart]: ShoppingBag,
}

const getStorefrontActionBadgeCount = (
  storefrontActionId: StorefrontActionId,
): number => {
  if (storefrontActionId === StorefrontActionId.Cart) {
    return cartCount.value
  }

  if (storefrontActionId === StorefrontActionId.Wishlist) {
    return wishlistCount.value
  }

  return 0
}

const shouldShowStorefrontActionBadge = (
  storefrontActionId: StorefrontActionId,
): boolean => getStorefrontActionBadgeCount(storefrontActionId) > 0

const getProductSearchTriggerValue = (
  storefrontActionId: StorefrontActionId,
): true | undefined =>
  storefrontActionId === StorefrontActionId.Search ? true : undefined

// Closes mobile navigation and forwards shell actions to their owning features.
const handleStorefrontAction = (storefrontActionId: StorefrontActionId): void => {
  closeNavigation()
  emit('action', storefrontActionId)
}

// Desktop pointer hover/focus opens a dropdown. The full nav item owns the
// hover area so moving from trigger to menu does not flicker closed.
const handleItemPointerEnter = (navigationItem: StorefrontNavigationItem): void => {
  if (navigationItem.children) {
    openDropdown(navigationItem.id)
  }
}

const handleItemPointerLeave = (navigationItem: StorefrontNavigationItem): void => {
  if (navigationItem.children) {
    closeDropdown()
  }
}

const handleItemFocusOut = (event: FocusEvent): void => {
  const currentTarget = event.currentTarget

  if (
    currentTarget instanceof HTMLElement
    && event.relatedTarget instanceof Node
    && currentTarget.contains(event.relatedTarget)
  ) {
    return
  }

  closeDropdown()
}

const handleTriggerClick = (navigationItem: StorefrontNavigationItem): void => {
  if (navigationItem.children) {
    openDropdown(navigationItem.id)
  }
}
</script>

<template>
  <header ref="navigationRoot" class="storefront-header">
    <nav
      class="storefront-navigation"
      :aria-label="STOREFRONT_COPY.navigationLabel"
    >
      <StorefrontBrandLogo />

      <ul class="desktop-navigation-list">
        <li
          v-for="navigationItem in STOREFRONT_NAVIGATION"
          :key="navigationItem.id"
          class="desktop-navigation-item"
          @mouseenter="handleItemPointerEnter(navigationItem)"
          @mouseleave="handleItemPointerLeave(navigationItem)"
          @focusin="handleItemPointerEnter(navigationItem)"
          @focusout="handleItemFocusOut"
        >
          <NuxtLink
            v-if="navigationItem.to"
            class="desktop-navigation-link"
            :class="{ 'is-active': isNavigationItemActive(navigationItem) }"
            :to="navigationItem.to"
            :external="!isInternalStorefrontRoute(navigationItem.to)"
          >
            {{ navigationItem.label }}
          </NuxtLink>

          <template v-else>
            <button
              class="desktop-navigation-link desktop-navigation-toggle"
              type="button"
              :class="{ 'is-active': isNavigationItemActive(navigationItem) }"
              :aria-controls="navigationItem.menuId"
              :aria-expanded="isDropdownOpen(navigationItem.id)"
              @click="handleTriggerClick(navigationItem)"
            >
              <span>{{ navigationItem.label }}</span>
              <StorefrontIcon
                :icon-node="ChevronDown"
                class="navigation-chevron"
                :class="{ 'is-open': isDropdownOpen(navigationItem.id) }"
                :size="14"
                :stroke-width="1.5"
                aria-hidden="true"
              />
            </button>

            <div
              v-show="isDropdownOpen(navigationItem.id)"
              :id="navigationItem.menuId"
              class="desktop-navigation-dropdown"
            >
              <NuxtLink
                v-for="navigationChild in navigationItem.children"
                :key="navigationChild.to"
                :to="navigationChild.to"
                :external="!isInternalStorefrontRoute(navigationChild.to)"
                @click="closeNavigation"
              >
                {{ navigationChild.label }}
              </NuxtLink>
            </div>
          </template>
        </li>
      </ul>

      <div class="desktop-header-actions">
        <button
          v-for="storefrontActionId in DESKTOP_HEADER_ACTIONS"
          :key="storefrontActionId"
          class="header-action-button"
          type="button"
          :data-product-search-trigger="getProductSearchTriggerValue(storefrontActionId)"
          :aria-label="STOREFRONT_ACTION_BY_ID[storefrontActionId].label"
          :title="STOREFRONT_ACTION_BY_ID[storefrontActionId].label"
          @click="handleStorefrontAction(storefrontActionId)"
        >
          <StorefrontIcon
            :icon-node="STOREFRONT_ACTION_ICONS[storefrontActionId]"
            :size="19"
            :stroke-width="1.4"
            aria-hidden="true"
          />
          <ClientOnly>
            <span
              v-if="shouldShowStorefrontActionBadge(storefrontActionId)"
              class="cart-count-badge"
              aria-hidden="true"
            >
              {{ getStorefrontActionBadgeCount(storefrontActionId) }}
            </span>
          </ClientOnly>
        </button>
      </div>

      <div class="mobile-header-actions">
        <button
          v-for="storefrontActionId in MOBILE_HEADER_ACTIONS"
          :key="storefrontActionId"
          class="header-action-button"
          type="button"
          :aria-label="STOREFRONT_ACTION_BY_ID[storefrontActionId].label"
          :title="STOREFRONT_ACTION_BY_ID[storefrontActionId].label"
          @click="handleStorefrontAction(storefrontActionId)"
        >
          <StorefrontIcon
            :icon-node="STOREFRONT_ACTION_ICONS[storefrontActionId]"
            :size="19"
            :stroke-width="1.4"
            aria-hidden="true"
          />
          <ClientOnly>
            <span
              v-if="shouldShowStorefrontActionBadge(storefrontActionId)"
              class="cart-count-badge"
              aria-hidden="true"
            >
              {{ getStorefrontActionBadgeCount(storefrontActionId) }}
            </span>
          </ClientOnly>
        </button>

        <button
          class="header-action-button"
          type="button"
          data-product-search-focus-fallback
          :aria-controls="STOREFRONT_COPY.mobileMenuId"
          :aria-expanded="isMobileMenuOpen"
          :aria-label="isMobileMenuOpen ? STOREFRONT_COPY.closeMenu : STOREFRONT_COPY.openMenu"
          :title="isMobileMenuOpen ? STOREFRONT_COPY.closeMenu : STOREFRONT_COPY.openMenu"
          @click="toggleMobileMenu"
        >
          <StorefrontMenuToggleIcon :is-open="isMobileMenuOpen" />
        </button>
      </div>

      <Transition name="mobile-navigation">
        <div
          v-show="isMobileMenuOpen"
          :id="STOREFRONT_COPY.mobileMenuId"
          class="mobile-navigation-panel"
        >
          <ul class="mobile-navigation-list">
            <li
              v-for="navigationItem in STOREFRONT_NAVIGATION"
              :key="navigationItem.id"
              class="mobile-navigation-item"
            >
              <NuxtLink
                v-if="navigationItem.to"
                class="mobile-navigation-link"
                :class="{ 'is-active': isNavigationItemActive(navigationItem) }"
                :to="navigationItem.to"
                :external="!isInternalStorefrontRoute(navigationItem.to)"
                @click="closeNavigation"
              >
                {{ navigationItem.label }}
              </NuxtLink>

              <template v-else>
                <button
                  class="mobile-navigation-link mobile-navigation-toggle"
                  type="button"
                  :aria-controls="`${navigationItem.menuId}-mobile`"
                  :aria-expanded="isDropdownOpen(navigationItem.id)"
                  @click="toggleDropdown(navigationItem.id)"
                >
                  <span>{{ navigationItem.label }}</span>
                  <StorefrontIcon
                    :icon-node="ChevronDown"
                    class="navigation-chevron"
                    :class="{ 'is-open': isDropdownOpen(navigationItem.id) }"
                    :size="16"
                    :stroke-width="1.5"
                    aria-hidden="true"
                  />
                </button>

                <div
                  v-show="isDropdownOpen(navigationItem.id)"
                  :id="`${navigationItem.menuId}-mobile`"
                  class="mobile-navigation-children"
                >
                  <NuxtLink
                    v-for="navigationChild in navigationItem.children"
                    :key="navigationChild.to"
                    :to="navigationChild.to"
                    :external="!isInternalStorefrontRoute(navigationChild.to)"
                    @click="closeNavigation"
                  >
                    {{ navigationChild.label }}
                  </NuxtLink>
                </div>
              </template>
            </li>
          </ul>

          <div class="mobile-menu-actions">
            <button
              v-for="storefrontActionId in MOBILE_MENU_ACTIONS"
              :key="storefrontActionId"
              type="button"
              :data-product-search-trigger="getProductSearchTriggerValue(storefrontActionId)"
              @click="handleStorefrontAction(storefrontActionId)"
            >
              <StorefrontIcon
                :icon-node="STOREFRONT_ACTION_ICONS[storefrontActionId]"
                :size="19"
                :stroke-width="1.5"
                aria-hidden="true"
              />
              <span>{{ STOREFRONT_ACTION_BY_ID[storefrontActionId].label }}</span>
            </button>
          </div>
        </div>
      </Transition>
    </nav>
  </header>
</template>
