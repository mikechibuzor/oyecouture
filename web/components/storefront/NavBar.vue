<script setup lang="ts">
import type { IconNode } from 'lucide'
import {
  ChevronDown,
  Heart,
  Menu,
  Search,
  ShoppingBag,
  UserRound,
  X,
} from 'lucide'
import {
  MOBILE_HEADER_ACTIONS,
  MOBILE_MENU_ACTIONS,
  STOREFRONT_ACTION_BY_ID,
  STOREFRONT_COPY,
  STOREFRONT_NAVIGATION,
  isInternalStorefrontRoute,
} from '~/constants/storefront'
import { StorefrontActionId } from '~/types/storefront'

const emit = defineEmits<{
  action: [storefrontActionId: StorefrontActionId]
}>()

const navigationRoot = ref<HTMLElement | null>(null)
const {
  closeNavigation,
  isDropdownOpen,
  isMobileMenuOpen,
  isNavigationItemActive,
  toggleDropdown,
  toggleMobileMenu,
} = useStorefrontNavigation(navigationRoot)

const STOREFRONT_ACTION_ICONS: Readonly<Record<StorefrontActionId, IconNode>> = {
  [StorefrontActionId.Search]: Search,
  [StorefrontActionId.Account]: UserRound,
  [StorefrontActionId.Wishlist]: Heart,
  [StorefrontActionId.Cart]: ShoppingBag,
}

// Forwards shell actions to the future search, account, wishlist, and cart features.
const handleStorefrontAction = (storefrontActionId: StorefrontActionId): void => {
  emit('action', storefrontActionId)
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
              @click="toggleDropdown(navigationItem.id)"
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
          v-for="storefrontAction in STOREFRONT_ACTION_BY_ID"
          :key="storefrontAction.id"
          class="header-action-button"
          type="button"
          :aria-label="storefrontAction.label"
          :title="storefrontAction.label"
          @click="handleStorefrontAction(storefrontAction.id)"
        >
          <StorefrontIcon
            :icon-node="STOREFRONT_ACTION_ICONS[storefrontAction.id]"
            :size="20"
            :stroke-width="1.5"
            aria-hidden="true"
          />
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
            :size="20"
            :stroke-width="1.5"
            aria-hidden="true"
          />
        </button>

        <button
          class="header-action-button"
          type="button"
          :aria-controls="STOREFRONT_COPY.mobileMenuId"
          :aria-expanded="isMobileMenuOpen"
          :aria-label="isMobileMenuOpen ? STOREFRONT_COPY.closeMenu : STOREFRONT_COPY.openMenu"
          :title="isMobileMenuOpen ? STOREFRONT_COPY.closeMenu : STOREFRONT_COPY.openMenu"
          @click="toggleMobileMenu"
        >
          <StorefrontIcon
            v-if="isMobileMenuOpen"
            :icon-node="X"
            :size="23"
            :stroke-width="1.5"
            aria-hidden="true"
          />
          <StorefrontIcon
            v-else
            :icon-node="Menu"
            :size="23"
            :stroke-width="1.5"
            aria-hidden="true"
          />
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
