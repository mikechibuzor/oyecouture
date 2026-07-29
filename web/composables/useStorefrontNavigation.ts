import type { Ref } from 'vue'
import type { StorefrontNavigationItem } from '~/types/storefront'
import { StorefrontNavigationItemId } from '~/types/storefront'

export const useStorefrontNavigation = (
  navigationRoot: Readonly<Ref<HTMLElement | null>>,
) => {
  const route = useRoute()
  const activeDropdownId = ref<StorefrontNavigationItemId | null>(null)
  const isMobileMenuOpen = ref(false)

  // Reports whether a dropdown is currently expanded.
  const isDropdownOpen = (navigationItemId: StorefrontNavigationItemId): boolean =>
    activeDropdownId.value === navigationItemId

  // Reports whether the current route belongs to a navigation item.
  const isNavigationItemActive = (
    navigationItem: StorefrontNavigationItem,
  ): boolean => {
    if (navigationItem.to === route.path) {
      return true
    }

    return Boolean(
      navigationItem.children?.some(
        (navigationChild) => navigationChild.to.split('#')[0] === route.path,
      ),
    )
  }

  // Opens or closes one dropdown while ensuring only one menu is expanded.
  const toggleDropdown = (
    navigationItemId: StorefrontNavigationItemId,
  ): void => {
    activeDropdownId.value = isDropdownOpen(navigationItemId)
      ? null
      : navigationItemId
  }

  // Closes every desktop or mobile navigation surface.
  const closeNavigation = (): void => {
    activeDropdownId.value = null
    isMobileMenuOpen.value = false
  }

  // Opens or closes the mobile menu and resets stale dropdown state when closing.
  const toggleMobileMenu = (): void => {
    isMobileMenuOpen.value = !isMobileMenuOpen.value

    if (!isMobileMenuOpen.value) {
      activeDropdownId.value = null
    }
  }

  // Closes navigation when focus moves to a pointer target outside the header.
  const handleDocumentPointerDown = (event: PointerEvent): void => {
    if (
      event.target instanceof Node
      && !navigationRoot.value?.contains(event.target)
    ) {
      closeNavigation()
    }
  }

  // Closes navigation when the user presses Escape.
  const handleDocumentKeyDown = (event: KeyboardEvent): void => {
    if (event.key === 'Escape') {
      closeNavigation()
    }
  }

  onMounted(() => {
    document.addEventListener('pointerdown', handleDocumentPointerDown)
    document.addEventListener('keydown', handleDocumentKeyDown)
  })

  onBeforeUnmount(() => {
    document.removeEventListener('pointerdown', handleDocumentPointerDown)
    document.removeEventListener('keydown', handleDocumentKeyDown)
  })

  watch(() => route.fullPath, closeNavigation)

  return {
    activeDropdownId: readonly(activeDropdownId),
    closeNavigation,
    isDropdownOpen,
    isMobileMenuOpen: readonly(isMobileMenuOpen),
    isNavigationItemActive,
    toggleDropdown,
    toggleMobileMenu,
  }
}
