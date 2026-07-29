import type {
  CatalogueFilterOption,
  CategoryPageConfig,
  CategoryTileConfig,
  ProductSortOption,
} from '~/types/product'
import {
  CatalogueFilterId,
  ProductSection,
  ProductSize,
  ProductSortId,
} from '~/types/product'
import { StorefrontRoute } from '~/types/storefront'

export const PRODUCT_API_QUERY_KEYS = {
  section: 'section',
  search: 'q',
} as const

export const PRODUCT_ASYNC_DATA_KEYS = {
  all: 'products-all',
  detailPrefix: 'product-detail',
  sectionPrefix: 'products-section',
} as const

export const PRODUCT_ROUTE_PREFIX = '/product' as const
export const PRODUCT_MEDIA_PLACEHOLDER_COUNT = 3 as const

export const PRODUCT_PRICE_FORMAT = {
  currency: 'NGN',
  locale: 'en-NG',
} as const

export const PRODUCT_COPY = {
  actions: {
    addToBag: 'Add to Bag',
    nextImage: 'Show next image',
    previousImage: 'Show previous image',
    showImage: 'Show image',
    viewProduct: 'View Product',
    wishlistAdd: 'Add to wishlist',
    wishlistRemove: 'Remove from wishlist',
  },
  bespokeNote:
    'Every piece is finished by hand. Sizing, fabric & delivery details are confirmed with us on WhatsApp after checkout.',
  breadcrumbs: {
    home: 'Home',
    label: 'Breadcrumb',
  },
  catalogue: {
    empty: 'No pieces match this selection.',
    error: 'The collection could not be loaded. Please try again.',
    eyebrow: 'The Collection',
    filterLabel: 'Filter the collection',
    heading: 'All Pieces',
    loading: 'Loading the collection...',
    pieceSingular: 'piece',
    piecePlural: 'pieces',
    sortLabel: 'Sort by',
  },
  categoryCollectionSuffix: 'Collection',
  categoryTilesHeading: 'Shop By Category',
  galleryLabel: 'Product gallery',
  home: {
    categoryEyebrow: 'Explore Oyecouture',
    heroAction: 'Discover Our Collection',
    heroDescription:
      'Invest in pieces that reflect your ambition brilliantly. Bespoke tailoring that speaks before you do.',
    heroHeading: 'Tailor Your Ambition With Effortless Precision',
  },
  imagePlaceholder: 'Photography coming soon',
  productDetail: {
    error: 'This piece could not be loaded. Please return to the collection.',
    loading: 'Loading this piece...',
    notFound: 'This piece is no longer available.',
  },
  size: {
    label: 'Size',
  },
} as const

export const CATALOGUE_FILTER_OPTIONS = [
  { id: CatalogueFilterId.All, label: 'All' },
  { id: CatalogueFilterId.Suits, label: 'Suits' },
  { id: CatalogueFilterId.AsoOke, label: 'Aso-Oke' },
  { id: CatalogueFilterId.Casuals, label: 'Casuals' },
  { id: CatalogueFilterId.Kaftans, label: 'Kaftans' },
] as const satisfies readonly CatalogueFilterOption[]

export const PRODUCT_SORT_OPTIONS = [
  { id: ProductSortId.Featured, label: 'Featured' },
  { id: ProductSortId.PriceAscending, label: 'Price: Low to High' },
  { id: ProductSortId.PriceDescending, label: 'Price: High to Low' },
  { id: ProductSortId.Newest, label: 'Newest' },
] as const satisfies readonly ProductSortOption[]

export const PRODUCT_SIZE_OPTIONS = [
  ProductSize.Small,
  ProductSize.Medium,
  ProductSize.Large,
  ProductSize.ExtraLarge,
] as const

export const DEFAULT_PRODUCT_SIZE = ProductSize.Medium

export const PRODUCT_FILTER_CATEGORIES = {
  asoOke: 'Aso-Oke',
  kaftans: 'Kaftans',
} as const

export const CASUAL_PRODUCT_SECTIONS = [
  ProductSection.CasualsMen,
  ProductSection.CasualsWomen,
] as const

export const CATEGORY_PAGE_BY_SECTION = {
  [ProductSection.Suits]: {
    section: ProductSection.Suits,
    eyebrow: 'Precision Tailoring',
    heading: 'Suits',
    description:
      'Structured silhouettes and considered details, made to a single measure for the rooms that matter.',
  },
  [ProductSection.Trad]: {
    section: ProductSection.Trad,
    eyebrow: 'Heritage, Reconsidered',
    heading: 'Trad Attires',
    description:
      'Aso-oke, agbada and kaftan forms refined through modern proportion and exacting handwork.',
  },
  [ProductSection.CasualsMen]: {
    section: ProductSection.CasualsMen,
    eyebrow: 'Ready To Wear',
    heading: 'Casuals For Men',
    description:
      'Relaxed shirts, trousers and knitwear with the same discipline as our tailoring.',
  },
  [ProductSection.CasualsWomen]: {
    section: ProductSection.CasualsWomen,
    eyebrow: 'Ready To Wear',
    heading: 'Casuals For Women',
    description:
      'Quietly expressive pieces cut for ease, movement and repeat wear.',
  },
} as const satisfies Readonly<Record<ProductSection, CategoryPageConfig>>

export const CATEGORY_TILES = [
  {
    label: 'Shop',
    heading: 'Suits',
    imagePath: '/images/catalogue/oye-cat-suits.webp',
    imageAlt: 'Oyecouture tailored suit',
    to: StorefrontRoute.Suits,
  },
  {
    label: 'Shop',
    heading: 'Trads',
    imagePath: '/images/catalogue/oye-cat-trads.webp',
    imageAlt: 'Oyecouture traditional attire',
    to: StorefrontRoute.TradAttires,
  },
  {
    label: 'Shop',
    heading: 'Casuals',
    imagePath: '/images/catalogue/oye-cat-casuals.webp',
    imageAlt: 'Oyecouture casual clothing',
    to: StorefrontRoute.CasualsMen,
  },
] as const satisfies readonly CategoryTileConfig[]

export const HOME_HERO_IMAGES = [
  {
    source: '/images/catalogue/oye-hero-c1.webp',
    alt: 'Oyecouture traditional menswear',
  },
  {
    source: '/images/catalogue/oye-hero-c2.webp',
    alt: 'Oyecouture ceremonial tailoring',
  },
] as const

export const PRODUCT_FALLBACK_IMAGES: Readonly<
  Partial<Record<string, readonly string[]>>
> = {
  monarch: ['/images/catalogue/oye-p-monarch.webp'],
  sovereign: ['/images/catalogue/oye-p-sovereign.webp'],
  regent: ['/images/catalogue/oye-p-regent.webp'],
  heir: ['/images/catalogue/oye-p-heir.webp'],
  ambassador: [
    '/images/catalogue/oye-p-ambassador-1.webp',
    '/images/catalogue/oye-p-ambassador-2.webp',
  ],
  chairman: ['/images/catalogue/oye-p-chairman.webp'],
  diplomat: ['/images/catalogue/oye-p-diplomat.webp'],
  weekend: ['/images/catalogue/oye-p-weekend.webp'],
  linen: [
    '/images/catalogue/oye-p-linen-1.webp',
    '/images/catalogue/oye-p-linen-2.webp',
    '/images/catalogue/oye-p-linen-3.webp',
  ],
  oxford: [
    '/images/catalogue/oye-p-oxford-1.webp',
    '/images/catalogue/oye-p-oxford-2.webp',
    '/images/catalogue/oye-p-oxford-3.webp',
  ],
  polo: [
    '/images/catalogue/oye-p-polo-1.webp',
    '/images/catalogue/oye-p-polo-2.webp',
    '/images/catalogue/oye-p-polo-3.webp',
  ],
  chino: [
    '/images/catalogue/oye-p-chino-1.webp',
    '/images/catalogue/oye-p-chino-2.webp',
    '/images/catalogue/oye-p-chino-3.webp',
  ],
  jogger: [
    '/images/catalogue/oye-p-jogger-1.webp',
    '/images/catalogue/oye-p-jogger-2.webp',
    '/images/catalogue/oye-p-jogger-3.webp',
  ],
  short: [
    '/images/catalogue/oye-p-short-1.webp',
    '/images/catalogue/oye-p-short-2.webp',
    '/images/catalogue/oye-p-short-3.webp',
  ],
  crew: [
    '/images/catalogue/oye-p-crew-1.webp',
    '/images/catalogue/oye-p-crew-2.webp',
    '/images/catalogue/oye-p-crew-3.webp',
  ],
  cardigan: [
    '/images/catalogue/oye-p-cardigan-1.webp',
    '/images/catalogue/oye-p-cardigan-2.webp',
    '/images/catalogue/oye-p-cardigan-3.webp',
  ],
  wrap: [
    '/images/catalogue/oye-p-wrap-1.webp',
    '/images/catalogue/oye-p-wrap-2.webp',
    '/images/catalogue/oye-p-wrap-3.webp',
  ],
  linenshirt: [
    '/images/catalogue/oye-p-linenshirt-1.webp',
    '/images/catalogue/oye-p-linenshirt-2.webp',
    '/images/catalogue/oye-p-linenshirt-3.webp',
  ],
}

// Builds the canonical client-side route for one product detail page.
export const getProductRoute = (productSlug: string): string =>
  `${PRODUCT_ROUTE_PREFIX}/${encodeURIComponent(productSlug)}`
