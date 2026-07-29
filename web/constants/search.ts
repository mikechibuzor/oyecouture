import type { ProductSearchTerm } from '~/types/search'

export const PRODUCT_SEARCH_COPY = {
  clearAction: 'Clear',
  closeAction: 'Close search',
  dialogLabel: 'Search products',
  error: 'The collection could not be searched. Please try again.',
  inputLabel: 'Search the collection',
  inputPlaceholder: 'Search the collection',
  loading: 'Searching the collection...',
  loadingRecommendations: 'Loading recommendations...',
  noResultsPrefix: 'No pieces match',
  recommendationsHeading: 'Recommended for you',
  resultPlural: 'results for',
  resultSingular: 'result for',
  topSearchesHeading: 'Top Searches',
} as const

export const PRODUCT_SEARCH_CONFIG = {
  debounceMilliseconds: 250,
  maximumQueryLength: 100,
  recommendationLimit: 6,
} as const

export const PRODUCT_SEARCH_ELEMENT_IDS = {
  input: 'storefront-product-search',
} as const

export const PRODUCT_SEARCH_FOCUS_SELECTORS = {
  fallback: '[data-product-search-focus-fallback]',
  trigger: '[data-product-search-trigger]',
} as const

export const PRODUCT_SEARCH_RECOMMENDED_SLUGS = [
  'agbada',
  'monarch',
  'senator',
  'linen',
  'classic',
  'wrap',
] as const

export const PRODUCT_SEARCH_TERMS = [
  { label: 'Suits', query: 'Suits' },
  { label: 'Agbada', query: 'Agbada' },
  { label: 'Kaftans', query: 'Kaftans' },
  { label: 'Aso-Oke', query: 'Aso-Oke' },
  { label: 'Casuals', query: 'Casuals' },
  { label: 'Blazers', query: 'Blazers' },
] as const satisfies readonly ProductSearchTerm[]

export const PRODUCT_SEARCH_FOCUSABLE_SELECTOR = [
  'a[href]',
  'button:not([disabled])',
  'input:not([disabled])',
].join(',')

// Builds the approved singular or plural results heading for a submitted query.
export const getProductSearchResultsHeading = (
  resultCount: number,
  searchQuery: string,
): string => {
  const resultLabel = resultCount === 1
    ? PRODUCT_SEARCH_COPY.resultSingular
    : PRODUCT_SEARCH_COPY.resultPlural

  return `${resultCount} ${resultLabel} "${searchQuery}"`
}

// Builds the approved empty-state message without duplicating display copy.
export const getProductSearchNoResultsMessage = (searchQuery: string): string =>
  `${PRODUCT_SEARCH_COPY.noResultsPrefix} "${searchQuery}".`
