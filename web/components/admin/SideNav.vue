<script setup lang="ts">
import type { IconNode } from 'lucide'
import {
  ChevronsLeft,
  ClipboardList,
  House,
  Menu,
  Settings,
  Shirt,
  Truck,
  X,
} from 'lucide'
import { ADMIN_COPY, ADMIN_NAVIGATION } from '~/constants/admin'
import type { AdminSession } from '~/types/admin'

interface AdminSideNavProps {
  session: AdminSession
}

defineProps<AdminSideNavProps>()

const emit = defineEmits<{ signOut: [] }>()

const route = useRoute()
const isCollapsed = ref(false)
const isMobileOpen = ref(false)
const sidebar = ref<HTMLElement | null>(null)
const mobileToggle = ref<HTMLButtonElement | null>(null)

const NAV_ICONS: Readonly<Record<string, IconNode>> = {
  orders: ClipboardList,
  products: Shirt,
  'home-content': House,
  delivery: Truck,
  settings: Settings,
}

const closeMobileNavigation = (): void => {
  isMobileOpen.value = false
}

useModalFocus({
  isOpen: isMobileOpen,
  container: sidebar,
  initialFocus: mobileToggle,
  onEscape: closeMobileNavigation,
})

watch(() => route.path, closeMobileNavigation)

// Only an exact match counts: every pending section would otherwise light up
// alongside Orders, which owns the admin root.
const isCurrentRoute = (itemRoute?: string): boolean =>
  itemRoute !== undefined && route.path === itemRoute
</script>

<template>
  <div
    v-if="isMobileOpen"
    class="admin-mobile-nav-backdrop"
    aria-hidden="true"
    @click="closeMobileNavigation"
  />
  <aside
    ref="sidebar"
    class="admin-sidebar"
    :class="{ 'is-collapsed': isCollapsed, 'is-mobile-open': isMobileOpen }"
    :role="isMobileOpen ? 'dialog' : undefined"
    :aria-modal="isMobileOpen ? 'true' : undefined"
    :aria-label="isMobileOpen ? ADMIN_COPY.navigationLabel : undefined"
    :tabindex="isMobileOpen ? -1 : undefined"
  >
    <div class="admin-sidebar-head">
      <p class="admin-sidebar-brand">
        <span class="admin-sidebar-mark" aria-hidden="true">O</span>
        <span class="admin-sidebar-wordmark">{{ ADMIN_COPY.brand }}</span>
      </p>

      <button
        type="button"
        class="admin-sidebar-collapse"
        :aria-label="isCollapsed ? ADMIN_COPY.expandSidebar : ADMIN_COPY.collapseSidebar"
        :aria-expanded="!isCollapsed"
        @click="isCollapsed = !isCollapsed"
      >
        <StorefrontIcon :icon-node="ChevronsLeft" :size="16" :stroke-width="1.75" />
      </button>

      <button
        ref="mobileToggle"
        type="button"
        class="admin-sidebar-mobile-toggle"
        :aria-label="isMobileOpen ? 'Close admin navigation' : 'Open admin navigation'"
        :aria-expanded="isMobileOpen"
        @click="isMobileOpen = !isMobileOpen"
      >
        <StorefrontIcon
          :icon-node="isMobileOpen ? X : Menu"
          :size="19"
          :stroke-width="1.75"
          aria-hidden="true"
        />
      </button>
    </div>

    <nav class="admin-sidebar-nav" :aria-label="ADMIN_COPY.navigationLabel">
      <div v-for="group in ADMIN_NAVIGATION" :key="group.id" class="admin-nav-group">
        <p class="admin-nav-group-label">{{ group.label }}</p>

        <ul>
          <li v-for="item in group.items" :key="item.id">
            <!-- The label is hidden in the collapsed rail, so the accessible
                 name and the hover tooltip both come from the attributes. -->
            <NuxtLink
              v-if="item.route"
              :to="item.route"
              class="admin-nav-item"
              :class="{ 'is-current': isCurrentRoute(item.route) }"
              :aria-current="isCurrentRoute(item.route) ? 'page' : undefined"
              :aria-label="item.label"
              :title="item.label"
              @click="closeMobileNavigation"
            >
              <StorefrontIcon
                :icon-node="NAV_ICONS[item.id]"
                :size="18"
                :stroke-width="1.75"
              />
              <span class="admin-nav-item-label">{{ item.label }}</span>
            </NuxtLink>

            <!-- Not yet built: shown so the shape of the admin is visible, but
                 inert rather than a link that goes nowhere. -->
            <span
              v-else
              class="admin-nav-item is-pending"
              aria-disabled="true"
              :title="`${item.label}: ${ADMIN_COPY.pendingBadge.toLowerCase()}`"
            >
              <StorefrontIcon
                :icon-node="NAV_ICONS[item.id]"
                :size="18"
                :stroke-width="1.75"
              />
              <span class="admin-nav-item-label">{{ item.label }}</span>
              <span class="admin-nav-pending">{{ ADMIN_COPY.pendingBadge }}</span>
            </span>
          </li>
        </ul>
      </div>
    </nav>

    <div class="admin-sidebar-foot">
      <p class="admin-sidebar-user">
        <span class="admin-sidebar-avatar" aria-hidden="true">
          {{ session.email.slice(0, 1).toUpperCase() }}
        </span>
        <span class="admin-sidebar-user-detail">
          <span class="admin-sidebar-user-email">{{ session.email }}</span>
          <span class="admin-sidebar-user-role">{{ session.role }}</span>
        </span>
      </p>

      <button type="button" class="admin-sidebar-signout" @click="emit('signOut')">
        {{ ADMIN_COPY.signOut }}
      </button>
    </div>
  </aside>
</template>
