import type { CSSProperties } from 'vue'

export interface ProductSearchTerm {
  readonly label: string
  readonly query: string
}

export interface SearchResultTransitionStyle extends CSSProperties {
  readonly '--search-result-index': number
}
