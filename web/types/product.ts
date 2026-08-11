import type { StorefrontRoute } from '~/types/storefront'

export enum ProductSection {
  Suits = 'suits',
  Trad = 'trad',
  CasualsMen = 'casuals-men',
  CasualsWomen = 'casuals-women',
}

export enum ProductTag {
  None = '',
  Bespoke = 'Bespoke',
  New = 'New',
}

export enum ProductSize {
  Small = 'S',
  Medium = 'M',
  Large = 'L',
  ExtraLarge = 'XL',
}

export enum CatalogueFilterId {
  All = 'all',
  Suits = 'suits',
  AsoOke = 'aso-oke',
  Casuals = 'casuals',
  Kaftans = 'kaftans',
}

export enum ProductSortId {
  Featured = 'featured',
  PriceAscending = 'price-ascending',
  PriceDescending = 'price-descending',
  Newest = 'newest',
}

export enum HomeHeroImageSize {
  Compact = 'is-compact',
  Tall = 'is-tall',
  Medium = 'is-medium',
  Large = 'is-large',
  Landscape = 'is-landscape',
}

export enum AsyncDataStatus {
  Idle = 'idle',
  Pending = 'pending',
  Success = 'success',
  Error = 'error',
}

export type AsyncDataStatusValue = `${AsyncDataStatus}`

export interface ProductImage {
  readonly id: string
  readonly url: string
  readonly position: number
}

export interface Product {
  readonly id: string
  readonly slug: string
  readonly name: string
  readonly section: ProductSection
  readonly category: string
  readonly priceKobo: number
  readonly description: string
  readonly tag: ProductTag
  readonly hasSizes: boolean
  readonly sortOrder: number
  readonly createdAt: string
  readonly images: readonly ProductImage[]
}

export interface ProductListResponse {
  readonly products: readonly Product[]
  readonly count: number
}

export interface ProductDetailResponse {
  readonly product: Product
}

export interface ProductApiErrorDetail {
  readonly code: string
  readonly message: string
}

export interface ProductApiErrorResponse {
  readonly error: ProductApiErrorDetail
}

export interface ProductFetchError {
  readonly data?: ProductApiErrorResponse
}

export interface CatalogueFilterOption {
  readonly id: CatalogueFilterId
  readonly label: string
}

export interface ProductSortOption {
  readonly id: ProductSortId
  readonly label: string
}

export interface CategoryPageConfig {
  readonly section: ProductSection
  readonly breadcrumb: string
  readonly heading: string
  readonly description: string
  readonly categories: readonly string[]
}

export interface CategoryTileConfig {
  readonly label: string
  readonly heading: string
  readonly imagePath: string
  readonly imageAlt: string
  readonly to: string
}

export interface HomeHeroImageConfig {
  readonly id: string
  readonly source: string
  readonly alt: string
  readonly size: HomeHeroImageSize
}

export interface ProductSectionNavigation {
  readonly label: string
  readonly to: StorefrontRoute
}

export interface ProductMediaItem {
  readonly id: string
  readonly source: string | null
  readonly alt: string
}

// One line inside a product disclosure panel. A `label` renders as an inline
// lead-in and carries its own punctuation — the shipping lines read
// "Production time: …" while the returns lines read "Manufacturing faults — …",
// so the separator cannot be hardcoded at the template.
export interface ProductPolicyLine {
  readonly label?: string
  readonly text: string
}

export interface ProductPolicy {
  readonly id: string
  readonly title: string
  readonly intro: string
  readonly lines: readonly ProductPolicyLine[]
  // Rendered as a sub-heading above the lines, when the section needs one.
  readonly linesHeading?: string
  readonly footnote?: string
}
