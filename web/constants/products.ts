import type {
  CatalogueFilterOption,
  CategoryFilterOption,
  CategoryPageConfig,
  CategoryTileConfig,
  HomeHeroImageConfig,
  ProductPolicy,
  ProductSectionNavigation,
  ProductSortOption,
} from '~/types/product'
import {
  CatalogueFilterId,
  HomeHeroImageSize,
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
    heroHeading: {
      lead: 'Tailor Your Ambition',
      tail: 'With Effortless Precision',
    },
  },
  productDetail: {
    error: 'This piece could not be loaded. Please return to the collection.',
    loading: 'Loading this piece...',
    notFound: 'This piece is no longer available.',
    orderEyebrow: 'Order details',
    readyToWearHeading: 'Choose your size',
    readyToWearIntro: 'Select a size before adding this piece to your bag.',
    bespokeHeading: 'Made for you',
    bespokeIntro: 'This piece is made to measure and ordered one at a time.',
    policiesEyebrow: 'Good to know',
    policiesHeading: 'Delivery & aftercare',
    policiesIntro: 'Timelines, returns, and care guidance for your piece.',
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
  { id: ProductSortId.PriceAscending, label: 'Price · Low to High' },
  { id: ProductSortId.PriceDescending, label: 'Price · High to Low' },
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

export const CATEGORY_ALL_FILTER = 'All'
export const CATEGORY_ALL_FILTER_OPTION: CategoryFilterOption = {
  label: CATEGORY_ALL_FILTER,
  value: CATEGORY_ALL_FILTER,
}

const categoryFilter = (value: string, label = value): CategoryFilterOption => ({
  label,
  value,
})

export const CATEGORY_PAGE_BY_SECTION = {
  [ProductSection.Suits]: {
    section: ProductSection.Suits,
    breadcrumb: 'Suits',
    heading: 'Suits',
    description:
      'Structured tailoring for the rooms where decisions are made: two-piece, three-piece and blazers, cut to a single measure.',
    categories: [
      categoryFilter('Three-Piece'),
      categoryFilter('Two-Piece'),
      categoryFilter('Blazers'),
    ],
  },
  [ProductSection.Trad]: {
    section: ProductSection.Trad,
    breadcrumb: 'Trad Attires',
    heading: 'Trad Attires',
    description:
      "Heritage cloth, cut for the room you're about to enter: kaftans, agbada and ceremonial two-piece attire, each finished by hand.",
    categories: [
      categoryFilter('Kaftans'),
      categoryFilter('Agbada'),
      categoryFilter('Aso-Oke', 'Two Piece'),
    ],
  },
  [ProductSection.CasualsMen]: {
    section: ProductSection.CasualsMen,
    breadcrumb: 'Casuals · Men',
    heading: 'Casuals / Men',
    description:
      'Off-duty pieces cut with the same discipline as our tailoring. Order online, sized to you.',
    categories: [],
  },
  [ProductSection.CasualsWomen]: {
    section: ProductSection.CasualsWomen,
    breadcrumb: 'Casuals · Women',
    heading: 'Casuals / Women',
    description:
      'Off-duty pieces cut with the same discipline as our tailoring. Order online, sized to you.',
    categories: [],
  },
} as const satisfies Readonly<Record<ProductSection, CategoryPageConfig>>

export const PRODUCT_SECTION_NAVIGATION = {
  [ProductSection.Suits]: {
    label: CATEGORY_PAGE_BY_SECTION[ProductSection.Suits].heading,
    to: StorefrontRoute.Suits,
  },
  [ProductSection.Trad]: {
    label: CATEGORY_PAGE_BY_SECTION[ProductSection.Trad].heading,
    to: StorefrontRoute.TradAttires,
  },
  [ProductSection.CasualsMen]: {
    label: CATEGORY_PAGE_BY_SECTION[ProductSection.CasualsMen].heading,
    to: StorefrontRoute.CasualsMen,
  },
  [ProductSection.CasualsWomen]: {
    label: CATEGORY_PAGE_BY_SECTION[ProductSection.CasualsWomen].heading,
    to: StorefrontRoute.CasualsWomen,
  },
} as const satisfies Readonly<Record<ProductSection, ProductSectionNavigation>>

// Trads and Casuals use the client's own photography from Cloudinary. Suits
// still points at the design-mockup mannequin because there is no suit shoot yet.
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
    imagePath:
      'https://res.cloudinary.com/exgithrw/image/upload/f_auto,q_auto,w_900/oyecouture/campaign-4',
    imageAlt: 'Oyecouture traditional attire',
    to: StorefrontRoute.TradAttires,
  },
  {
    label: 'Shop',
    heading: 'Casuals',
    imagePath:
      'https://res.cloudinary.com/exgithrw/image/upload/f_auto,q_auto,w_900/oyecouture/casual-14',
    imageAlt: 'Oyecouture casual clothing',
    to: StorefrontRoute.CasualsMen,
  },
] as const satisfies readonly CategoryTileConfig[]

export const HOME_HERO_IMAGES = [
  {
    id: 'hero-image-one',
    source: '/images/catalogue/oye-hero-c1.webp',
    alt: 'Oyecouture womenswear portrait',
    size: HomeHeroImageSize.Compact,
  },
  {
    id: 'hero-image-two',
    source: '/images/catalogue/oye-hero-c2.webp',
    alt: 'Oyecouture double-breasted suit',
    size: HomeHeroImageSize.Tall,
  },
  {
    id: 'hero-image-three',
    source: '/images/catalogue/oye-hero-c3.webp',
    alt: 'Oyecouture casual menswear portrait',
    size: HomeHeroImageSize.Medium,
  },
  {
    id: 'hero-image-four',
    source: '/images/catalogue/oye-hero-c4.webp',
    alt: 'Oyecouture ceremonial traditional attire',
    size: HomeHeroImageSize.Large,
  },
  {
    id: 'hero-image-five',
    source: '/images/catalogue/oye-hero-c5.webp',
    alt: 'Oyecouture traditional menswear portrait',
    size: HomeHeroImageSize.Landscape,
  },
] as const satisfies readonly HomeHeroImageConfig[]

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

// Delivery and returns policy shown on every product page. The copy is the
// approved design's verbatim. It sets expectations a made-to-order buyer
// needs before paying, so it is deliberately not per-product.
export const PRODUCT_POLICIES: readonly ProductPolicy[] = [
  {
    id: 'shipping',
    title: 'Shipping & Delivery',
    intro:
      'Every Oyecouture piece is made to order, so delivery works differently than off-the-shelf fashion.',
    lines: [
      { label: 'Production time:', text: '3 week or less from confirmed measurements and payment.' },
      { label: 'Within Nigeria:', text: '1–3 working days after work is completed' },
      {
        label: 'International:',
        text: 'We deliver worldwide. It could take 5-7 days to get delivered to doorstep',
      },
      {
        label: 'Tracking:',
        text: 'You’ll receive a tracking link by email/SMS once your order ships.',
      },
      {
        label: 'Packaging:',
        text: 'Every order arrives in Oyecouture signature packaging, including a garment bag, dust protection, and care card.',
      },
    ],
    footnote:
      'Need it sooner? We can expedite production for an additional 20% of the piece’s price. If you’re working to a date, such as a wedding or event, let us know when you order and we’ll prioritise your piece.',
  },
  {
    id: 'returns',
    title: 'Returns & Exchanges',
    intro:
      'Because Oyecouture pieces are made to your measurements, we don’t offer standard returns. A bespoke jacket cut for you isn’t something we can resell to someone else.',
    linesHeading: 'What we do offer',
    lines: [
      {
        label: 'Manufacturing faults',
        text: 'If something’s wrong with the construction or material, we can reach a consensus to repair, remake, or refund. Reach out within 3 days of receiving your order.',
      },
      {
        label: 'Ready-to-wear items',
        text: '(non-bespoke), if any are sold: exchange or store credit within 3 days, unworn and with tags attached.',
      },
    ],
    footnote:
      'To start an alteration or flag an issue, contact us with your order number. See Contact Us below.',
  },
]
