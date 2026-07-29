import { API_ROUTES } from '~/constants/api'
import { PRODUCT_API_QUERY_KEYS } from '~/constants/products'
import {
  PRODUCT_SEARCH_CONFIG,
  PRODUCT_SEARCH_RECOMMENDED_SLUGS,
} from '~/constants/search'
import type { Product, ProductListResponse } from '~/types/product'
import { AsyncDataStatus } from '~/types/product'

// Manages curated recommendations and race-safe debounced product searches.
export const useProductSearch = () => {
  const runtimeConfig = useRuntimeConfig()
  const searchQuery = ref('')
  const searchResults = ref<readonly Product[]>([])
  const recommendedProducts = ref<readonly Product[]>([])
  const searchStatus = ref<AsyncDataStatus>(AsyncDataStatus.Idle)
  const recommendationStatus = ref<AsyncDataStatus>(AsyncDataStatus.Idle)
  let searchDebounceTimer: ReturnType<typeof setTimeout> | undefined
  let searchRequestController: AbortController | null = null
  let recommendationRequestController: AbortController | null = null

  // Normalizes whitespace before a query reaches the public API.
  const normalizedSearchQuery = computed(() => searchQuery.value.trim())
  // Reports whether the customer has entered a meaningful query.
  const hasSearchQuery = computed(() => normalizedSearchQuery.value.length > 0)

  // Cancels the current debounce window before scheduling another search.
  const clearSearchDebounce = (): void => {
    if (searchDebounceTimer) {
      clearTimeout(searchDebounceTimer)
      searchDebounceTimer = undefined
    }
  }

  // Cancels an obsolete in-flight query so stale responses cannot replace newer results.
  const abortSearchRequest = (): void => {
    searchRequestController?.abort()
    searchRequestController = null
  }

  // Requests products with an optional server-side search filter.
  const requestProducts = (
    requestSignal: AbortSignal,
    requestedSearchQuery?: string,
  ): Promise<ProductListResponse> =>
    $fetch<ProductListResponse>(API_ROUTES.products, {
      baseURL: runtimeConfig.public.apiBase,
      query: requestedSearchQuery
        ? { [PRODUCT_API_QUERY_KEYS.search]: requestedSearchQuery }
        : undefined,
      retry: 0,
      signal: requestSignal,
    })

  // Executes the latest normalized query and ignores every superseded response.
  const executeSearch = async (): Promise<void> => {
    const requestedSearchQuery = normalizedSearchQuery.value

    if (!requestedSearchQuery) {
      return
    }

    // Step 1: own a fresh abort signal for this exact query.
    abortSearchRequest()
    const requestController = new AbortController()
    searchRequestController = requestController
    searchStatus.value = AsyncDataStatus.Pending

    try {
      // Step 2: let the API search name, category, section, and description.
      const productResponse = await requestProducts(
        requestController.signal,
        requestedSearchQuery,
      )

      // Step 3: publish only when this request still owns the active search.
      if (
        requestController.signal.aborted
        || searchRequestController !== requestController
      ) {
        return
      }

      searchResults.value = productResponse.products
      searchStatus.value = AsyncDataStatus.Success
    } catch {
      if (requestController.signal.aborted) {
        return
      }

      searchResults.value = []
      searchStatus.value = AsyncDataStatus.Error
    }
  }

  // Loads the six approved recommendations once and preserves their editorial order.
  const loadRecommendedProducts = async (): Promise<void> => {
    if (
      recommendationStatus.value === AsyncDataStatus.Pending
      || recommendationStatus.value === AsyncDataStatus.Success
    ) {
      return
    }

    recommendationRequestController?.abort()
    const requestController = new AbortController()
    recommendationRequestController = requestController
    recommendationStatus.value = AsyncDataStatus.Pending

    try {
      // Step 1: load the active public catalogue from the same endpoint as search.
      const productResponse = await requestProducts(requestController.signal)
      const productsBySlug = new Map(
        productResponse.products.map((product) => [product.slug, product]),
      )

      // Step 2: map the editorial slug order to the current active product records.
      recommendedProducts.value = PRODUCT_SEARCH_RECOMMENDED_SLUGS
        .map((productSlug) => productsBySlug.get(productSlug))
        .filter((product): product is Product => Boolean(product))
        .slice(0, PRODUCT_SEARCH_CONFIG.recommendationLimit)
      recommendationStatus.value = AsyncDataStatus.Success
    } catch {
      if (requestController.signal.aborted) {
        return
      }

      recommendedProducts.value = []
      recommendationStatus.value = AsyncDataStatus.Error
    }
  }

  // Clears customer input and returns the live search state to idle.
  const clearProductSearch = (): void => {
    clearSearchDebounce()
    abortSearchRequest()
    searchQuery.value = ''
    searchResults.value = []
    searchStatus.value = AsyncDataStatus.Idle
  }

  // Debounces normalized input and cancels work that no longer matches the field.
  watch(normalizedSearchQuery, (updatedSearchQuery) => {
    clearSearchDebounce()
    abortSearchRequest()

    if (!updatedSearchQuery) {
      searchResults.value = []
      searchStatus.value = AsyncDataStatus.Idle
      return
    }

    searchStatus.value = AsyncDataStatus.Pending
    searchDebounceTimer = setTimeout(() => {
      void executeSearch()
    }, PRODUCT_SEARCH_CONFIG.debounceMilliseconds)
  })

  // Releases timers and requests when the owning overlay leaves the application.
  onScopeDispose(() => {
    clearSearchDebounce()
    abortSearchRequest()
    recommendationRequestController?.abort()
  })

  return {
    clearProductSearch,
    hasSearchQuery,
    loadRecommendedProducts,
    normalizedSearchQuery,
    recommendedProducts,
    recommendationStatus,
    searchQuery,
    searchResults,
    searchStatus,
  }
}
