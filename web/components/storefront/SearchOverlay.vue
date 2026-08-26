<script setup lang="ts">
import type { IconNode } from 'lucide'
import { Search, X } from 'lucide'
import {
  PRODUCT_SEARCH_CONFIG,
  PRODUCT_SEARCH_COPY,
  PRODUCT_SEARCH_ELEMENT_IDS,
  PRODUCT_SEARCH_FOCUSABLE_SELECTOR,
  PRODUCT_SEARCH_FOCUS_SELECTORS,
  PRODUCT_SEARCH_TERMS,
  getProductSearchNoResultsMessage,
  getProductSearchResultsHeading,
} from '~/constants/search'
import { AsyncDataStatus } from '~/types/product'
import type {
  ProductSearchTerm,
  SearchResultTransitionStyle,
} from '~/types/search'

interface SearchOverlayProps {
  isOpen: boolean
}

const props = defineProps<SearchOverlayProps>()
const emit = defineEmits<{
  close: []
}>()
const searchInput = ref<HTMLInputElement | null>(null)
const searchOverlayPanel = ref<HTMLElement | null>(null)
let previouslyFocusedElement: HTMLElement | null = null
let previousBodyOverflow = ''
const {
  clearProductSearch,
  hasSearchQuery,
  loadRecommendedProducts,
  normalizedSearchQuery,
  recommendedProducts,
  recommendationStatus,
  searchQuery,
  searchResults,
  searchStatus,
} = useProductSearch()

const SEARCH_OVERLAY_ICONS = {
  close: X,
  search: Search,
} as const satisfies Readonly<Record<string, IconNode>>

// Selects live query results or curated recommendations for the current overlay state.
const displayedProducts = computed(() =>
  hasSearchQuery.value ? searchResults.value : recommendedProducts.value,
)
// Returns the approved recommendation heading or a count-aware query heading.
const resultsHeading = computed(() =>
  hasSearchQuery.value
    ? getProductSearchResultsHeading(
        searchResults.value.length,
        normalizedSearchQuery.value,
      )
    : PRODUCT_SEARCH_COPY.recommendationsHeading,
)
// Reports whether the current result source is waiting on the product API.
const isLoading = computed(() =>
  hasSearchQuery.value
    ? searchStatus.value === AsyncDataStatus.Pending
    : recommendationStatus.value === AsyncDataStatus.Pending,
)
// Reports whether the current result source failed to load.
const hasError = computed(() =>
  hasSearchQuery.value
    ? searchStatus.value === AsyncDataStatus.Error
    : recommendationStatus.value === AsyncDataStatus.Error,
)
// Reports whether a completed query produced no matching pieces.
const hasNoSearchResults = computed(() =>
  hasSearchQuery.value
  && searchStatus.value === AsyncDataStatus.Success
  && searchResults.value.length === 0,
)
// Builds the live-region message for a completed empty query.
const noResultsMessage = computed(() =>
  getProductSearchNoResultsMessage(normalizedSearchQuery.value),
)
// Selects the correct loading copy for search or recommendations.
const loadingMessage = computed(() =>
  hasSearchQuery.value
    ? PRODUCT_SEARCH_COPY.loading
    : PRODUCT_SEARCH_COPY.loadingRecommendations,
)

// Finds the first visible element that can safely receive restored focus.
const findVisibleFocusTarget = (focusSelector: string): HTMLElement | null =>
  Array.from(document.querySelectorAll<HTMLElement>(focusSelector))
    .find((focusTarget) => focusTarget.offsetParent !== null) ?? null

// Emits the single dismissal event used by close, Escape, backdrop, and result selection.
const closeSearchOverlay = (): void => {
  emit('close')
}

// Clears the current query and returns keyboard focus to the search field.
const handleClearSearch = (): void => {
  clearProductSearch()
  searchInput.value?.focus()
}

// Copies an approved top-search term into the live query field.
const handleSearchTermSelection = (searchTerm: ProductSearchTerm): void => {
  searchQuery.value = searchTerm.query
  searchInput.value?.focus()
}

// Closes the overlay while the selected product continues client-side navigation.
const handleProductSelection = (): void => {
  closeSearchOverlay()
}

// Returns a typed stagger index consumed by the result entrance transition.
const getSearchResultTransitionStyle = (
  searchResultIndex: number,
): SearchResultTransitionStyle => ({
  '--search-result-index': Math.min(searchResultIndex, 5),
})

// Keeps keyboard focus inside the modal and supports immediate Escape dismissal.
const handleOverlayKeydown = (keyboardEvent: KeyboardEvent): void => {
  if (keyboardEvent.key === 'Escape') {
    keyboardEvent.preventDefault()
    closeSearchOverlay()
    return
  }

  if (keyboardEvent.key !== 'Tab' || !searchOverlayPanel.value) {
    return
  }

  // Step 1: collect only controls that are currently visible in the overlay.
  const focusableElements = Array.from(
    searchOverlayPanel.value.querySelectorAll<HTMLElement>(
      PRODUCT_SEARCH_FOCUSABLE_SELECTOR,
    ),
  ).filter((focusableElement) => focusableElement.offsetParent !== null)

  if (focusableElements.length === 0) {
    keyboardEvent.preventDefault()
    return
  }

  // Step 2: wrap focus at either edge while leaving normal middle traversal intact.
  const firstFocusableElement = focusableElements[0]!
  const lastFocusableElement = focusableElements[focusableElements.length - 1]!
  const activeElement = document.activeElement

  if (keyboardEvent.shiftKey && activeElement === firstFocusableElement) {
    keyboardEvent.preventDefault()
    lastFocusableElement.focus()
  } else if (!keyboardEvent.shiftKey && activeElement === lastFocusableElement) {
    keyboardEvent.preventDefault()
    firstFocusableElement.focus()
  }
}

// Locks page scrolling, loads recommendations, and restores focus across overlay state changes.
watch(
  () => props.isOpen,
  async (isOpen) => {
    if (isOpen) {
      // Step 1: remember the trigger and prevent the page beneath from moving.
      previouslyFocusedElement = document.activeElement instanceof HTMLElement
        ? document.activeElement
        : null
      previousBodyOverflow = document.body.style.overflow
      document.body.style.overflow = 'hidden'

      // Step 2: load editorial results and focus the transactional field.
      void loadRecommendedProducts()
      await nextTick()
      searchInput.value?.focus()
      return
    }

    // Step 3: reset transient search state before returning focus to its trigger.
    document.body.style.overflow = previousBodyOverflow
    clearProductSearch()
    await nextTick()
    const focusReturnTarget = previouslyFocusedElement !== null
      && previouslyFocusedElement.offsetParent !== null
      ? previouslyFocusedElement
      : findVisibleFocusTarget(PRODUCT_SEARCH_FOCUS_SELECTORS.trigger)
        ?? findVisibleFocusTarget(PRODUCT_SEARCH_FOCUS_SELECTORS.fallback)
    focusReturnTarget?.focus()
    previouslyFocusedElement = null
  },
)

// Restores global page state if the layout unmounts while search is open.
onBeforeUnmount(() => {
  document.body.style.overflow = previousBodyOverflow
})
</script>

<template>
  <Teleport to="body">
    <Transition name="search-overlay">
      <div
        v-if="isOpen"
        class="search-overlay"
        role="dialog"
        aria-modal="true"
        :aria-label="PRODUCT_SEARCH_COPY.dialogLabel"
        @click.self="closeSearchOverlay"
        @keydown="handleOverlayKeydown"
      >
        <div ref="searchOverlayPanel" class="search-overlay-panel">
          <header class="search-overlay-header">
            <StorefrontBrandLogo />
            <button
              class="search-overlay-close"
              type="button"
              :aria-label="PRODUCT_SEARCH_COPY.closeAction"
              :title="PRODUCT_SEARCH_COPY.closeAction"
              @click="closeSearchOverlay"
            >
              <StorefrontIcon
                :icon-node="SEARCH_OVERLAY_ICONS.close"
                :size="24"
                :stroke-width="1.4"
                aria-hidden="true"
              />
            </button>
          </header>

          <form class="search-overlay-form" role="search" @submit.prevent>
            <StorefrontIcon
              :icon-node="SEARCH_OVERLAY_ICONS.search"
              :size="19"
              :stroke-width="1.4"
              aria-hidden="true"
            />
            <label class="visually-hidden" :for="PRODUCT_SEARCH_ELEMENT_IDS.input">
              {{ PRODUCT_SEARCH_COPY.inputLabel }}
            </label>
            <input
              :id="PRODUCT_SEARCH_ELEMENT_IDS.input"
              ref="searchInput"
              v-model="searchQuery"
              type="search"
              autocomplete="off"
              :maxlength="PRODUCT_SEARCH_CONFIG.maximumQueryLength"
              :placeholder="PRODUCT_SEARCH_COPY.inputPlaceholder"
            >
            <button
              v-if="hasSearchQuery"
              class="search-overlay-clear"
              type="button"
              @click="handleClearSearch"
            >
              {{ PRODUCT_SEARCH_COPY.clearAction }}
            </button>
          </form>

          <section v-if="!hasSearchQuery" class="search-top-terms">
            <h2>{{ PRODUCT_SEARCH_COPY.topSearchesHeading }}</h2>
            <div class="search-top-terms-grid">
              <button
                v-for="searchTerm in PRODUCT_SEARCH_TERMS"
                :key="searchTerm.query"
                type="button"
                @click="handleSearchTermSelection(searchTerm)"
              >
                {{ searchTerm.label }}
              </button>
            </div>
          </section>

          <section
            class="search-results-section"
            :aria-busy="isLoading"
            aria-live="polite"
          >
            <h2>{{ resultsHeading }}</h2>

            <Transition name="search-state" mode="out-in">
              <p v-if="isLoading" key="loading" class="search-state-message">
                {{ loadingMessage }}
              </p>
              <p v-else-if="hasError" key="error" class="search-state-message is-error">
                {{ PRODUCT_SEARCH_COPY.error }}
              </p>
              <p
                v-else-if="hasNoSearchResults"
                key="empty"
                class="search-state-message"
              >
                {{ noResultsMessage }}
              </p>
              <TransitionGroup
                v-else
                key="results"
                name="search-result"
                tag="div"
                class="search-results-grid"
              >
                <StorefrontSearchResultCard
                  v-for="(product, productIndex) in displayedProducts"
                  :key="product.id"
                  :product="product"
                  :style="getSearchResultTransitionStyle(productIndex)"
                  @select="handleProductSelection"
                />
              </TransitionGroup>
            </Transition>
          </section>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
