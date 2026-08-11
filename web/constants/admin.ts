import { OrderKind } from '~/types/order'
import type {
  AdminProductActiveFilter,
  AdminFilterOption,
  AdminMetricCard,
  AdminNavGroup,
  AdminStatusAction,
} from '~/types/admin'
import { AdminFilterType, MetricsPeriod, OrderStatus } from '~/types/admin'
import { ProductSection } from '~/types/product'
import type { SelectOption } from '~/types/ui'

export const ADMIN_ROUTES = {
  login: '/admin/login',
  orders: '/admin',
  delivery: '/admin/delivery',
  orderDetail: (orderNo: string) => `/admin/orders/${encodeURIComponent(orderNo)}`,
  products: '/admin/products',
  productNew: '/admin/products/new',
  productEdit: (slug: string) => `/admin/products/${encodeURIComponent(slug)}`,
} as const

export const ADMIN_API_ROUTES = {
  login: '/admin/login',
  logout: '/admin/logout',
  me: '/admin/me',
  orders: '/admin/orders',
  orderMetrics: '/admin/orders/metrics',
  orderDetail: (orderNo: string) => `/admin/orders/${encodeURIComponent(orderNo)}`,
  orderStatus: (orderNo: string) => `/admin/orders/${encodeURIComponent(orderNo)}/status`,
  deliveryZones: '/admin/delivery-zones',
  deliveryZone: (name: string) => `/admin/delivery-zones/${encodeURIComponent(name)}`,
  products: '/admin/products',
  productDetail: (slug: string) => `/admin/products/${encodeURIComponent(slug)}`,
  productActive: (slug: string) => `/admin/products/${encodeURIComponent(slug)}/active`,
  productImages: (slug: string) => `/admin/products/${encodeURIComponent(slug)}/images`,
  productImagesReorder: (slug: string) => `/admin/products/${encodeURIComponent(slug)}/images/reorder`,
  productImage: (slug: string, imageID: string) =>
    `/admin/products/${encodeURIComponent(slug)}/images/${encodeURIComponent(imageID)}`,
} as const

export const ADMIN_COPY = {
  brand: 'Oyecouture',
  panel: 'Admin',
  signIn: 'Sign in',
  signOut: 'Sign out',
  signingIn: 'Signing in…',
  loginHeading: 'Order desk',
  loginIntro: 'Sign in to manage orders.',
  emailLabel: 'Email',
  passwordLabel: 'Password',
  loginFailed: 'That email and password did not match.',
  loginThrottled: 'Too many attempts. Wait a few minutes and try again.',
  loginUnreachable: 'Could not reach the server. Check your connection.',
  ordersHeading: 'Orders',
  ordersSubheading: 'Manage orders',
  productsHeading: 'Products',
  productsSubheading: 'Manage catalogue availability',
  deliveryHeading: 'Delivery',
  deliverySubheading: 'Manage checkout delivery areas and fees',
  deliveryLoading: 'Loading delivery zones…',
  deliveryLoadFailed: 'Delivery zones could not be loaded.',
  deliveryEmpty: 'No delivery zones yet.',
  deliveryZoneLabel: 'Area',
  deliveryFeeLabel: 'Fee in naira',
  deliveryActiveLabel: 'Available at checkout',
  deliverySave: 'Save',
  deliverySaving: 'Saving…',
  deliverySavedNotice: 'Delivery zone saved',
  deliverySaveFailed: 'Delivery zone could not be saved.',
  newProduct: 'New product',
  editProduct: 'Edit',
  closeProductEditor: 'Close product editor',
  productNewHeading: 'New product',
  productEditHeading: 'Edit product',
  saveProduct: 'Save product',
  savingProduct: 'Saving…',
  productSaved: 'Product saved.',
  productCreatedNotice: 'Product created',
  productUpdatedNotice: 'Product saved',
  productArchivedNotice: 'Product archived',
  productRestoredNotice: 'Product restored',
  productSaveFailed: 'The product could not be saved.',
  productDuplicateSlug: 'That slug already exists.',
  productDetailLoading: 'Loading product…',
  productDetailFailed: 'The product could not be loaded.',
  productImagesHeading: 'Images',
  productImagesIntro: 'Upload, reorder, and remove product gallery images.',
  uploadProductImage: 'Upload image',
  uploadingProductImage: 'Uploading…',
  productImageUploadFailed: 'The image could not be uploaded.',
  productImageUploadedNotice: 'Image uploaded',
  productImageUpdateFailed: 'The image order could not be updated.',
  productImageUpdatedNotice: 'Image order updated',
  productImageDeleteFailed: 'The image could not be removed.',
  productImageDeletedNotice: 'Image removed',
  productImagesEmpty: 'No images yet.',
  productImageMoveUp: 'Move up',
  productImageMoveDown: 'Move down',
  productImageRemove: 'Remove',
  backToProducts: 'Back to products',
  slugLabel: 'Slug',
  nameLabel: 'Name',
  categoryLabel: 'Category',
  priceNairaLabel: 'Price in naira',
  descriptionLabel: 'Description',
  tagLabel: 'Tag',
  sortOrderLabel: 'Sort order',
  hasSizesLabel: 'Ready-to-wear sizes',
  activeLabel: 'Visible on storefront',
  searchLabel: 'Search orders',
  searchPlaceholder: 'Order number, name, email or phone',
  productSearchLabel: 'Search products',
  productSearchPlaceholder: 'Name, slug, category or description',
  productSectionLabel: 'Section',
  productActiveLabel: 'Availability',
  statusLabel: 'Status',
  kindLabel: 'Type',
  anyStatus: 'All statuses',
  anyKind: 'All types',
  loading: 'Loading orders…',
  loadFailed: 'The orders could not be loaded.',
  productLoading: 'Loading products…',
  productLoadFailed: 'The products could not be loaded.',
  retry: 'Try again',
  empty: 'No orders yet.',
  emptyFiltered: 'No orders match these filters.',
  productEmpty: 'No products yet.',
  productEmptyFiltered: 'No products match these filters.',
  productUpdateFailed: 'The product could not be updated.',
  clearFilters: 'Clear filters',
  archiveProduct: 'Archive',
  restoreProduct: 'Restore',
  archivingProduct: 'Archiving…',
  restoringProduct: 'Restoring…',
  previous: 'Previous',
  next: 'Next',
  pendingBadge: 'Soon',
  navigationLabel: 'Admin sections',
  collapseSidebar: 'Collapse sidebar',
  expandSidebar: 'Expand sidebar',
  periodLabel: 'Summary period',
  filter: 'Filter',
  filterEmpty: 'All',
  appliedFiltersLabel: 'Applied filters',
  clearFilter: 'Clear filter',
  backToOrders: 'Back to orders',
  orderDetailLoading: 'Loading order…',
  orderDetailFailed: 'The order could not be loaded.',
  customerDetails: 'Customer',
  orderItems: 'Items',
  measurements: 'Measurements',
  noMeasurements: 'No measurements submitted.',
  paymentLink: 'Payment link',
  copyPaymentLink: 'Copy payment link',
  copied: 'Copied',
  whatsappReply: 'WhatsApp reply',
  statusActions: 'Status actions',
  noStatusActions: 'No manual status actions are available.',
  statusUpdateFailed: 'The status could not be updated.',
  statusUpdating: 'Updating…',
  subtotal: 'Subtotal',
  delivery: 'Delivery',
  total: 'Total',
} as const

// The summary strip covers this window by default.
export const ADMIN_DEFAULT_PERIOD = MetricsPeriod.AllTime

export const METRICS_PERIOD_LABELS: Readonly<Record<MetricsPeriod, string>> = {
  [MetricsPeriod.Today]: 'Today',
  [MetricsPeriod.LastWeek]: 'Last week',
  [MetricsPeriod.Last14Days]: 'Last 14 days',
  [MetricsPeriod.Last60Days]: 'Last 60 days',
  [MetricsPeriod.Last90Days]: 'Last 90 days',
  [MetricsPeriod.AllTime]: 'All time',
}

export const METRICS_PERIOD_OPTIONS: readonly MetricsPeriod[] = [
  MetricsPeriod.Today,
  MetricsPeriod.LastWeek,
  MetricsPeriod.Last14Days,
  MetricsPeriod.Last60Days,
  MetricsPeriod.Last90Days,
  MetricsPeriod.AllTime,
]

// Ten rows a page, matching the house admin pattern.
export const ADMIN_ORDERS_PAGE_SIZE = 10

// Product rows are denser than orders, so the catalogue desk can show more at once.
export const ADMIN_PRODUCTS_PAGE_SIZE = 25

// Typing in the search box should not fire a request per keystroke.
export const ADMIN_SEARCH_DEBOUNCE_MS = 250

// The summary strip. The first card clears the status filter; the rest apply
// theirs, so the numbers double as one-click filters.
export const ADMIN_METRIC_CARDS: readonly AdminMetricCard[] = [
  { id: 'total', label: 'All orders', statusFilter: null },
  {
    id: 'pendingConfirmation',
    label: 'Pending confirmation',
    statusFilter: OrderStatus.PendingConfirmation,
  },
  {
    id: 'awaitingPayment',
    label: 'Awaiting payment',
    statusFilter: OrderStatus.AwaitingPayment,
  },
  { id: 'paid', label: 'Paid', statusFilter: OrderStatus.Paid },
  { id: 'revenueKobo', label: 'Revenue', statusFilter: null, isMoney: true },
]

// The shape of the finished admin. Unbuilt sections are marked pending rather
// than linking nowhere.
export const ADMIN_NAVIGATION: readonly AdminNavGroup[] = [
  {
    id: 'main',
    label: 'Main',
    items: [
      { id: 'orders', label: 'Orders', route: ADMIN_ROUTES.orders },
      { id: 'products', label: 'Products', route: ADMIN_ROUTES.products },
      { id: 'home-content', label: 'Home content', isPending: true },
      { id: 'delivery', label: 'Delivery', route: ADMIN_ROUTES.delivery },
    ],
  },
  {
    id: 'other',
    label: 'Other',
    items: [{ id: 'settings', label: 'Settings', isPending: true }],
  },
]

// Status labels and tones. Tone drives the badge colour only.
export const ORDER_STATUS_LABELS: Readonly<Record<OrderStatus, string>> = {
  [OrderStatus.PendingConfirmation]: 'Pending confirmation',
  [OrderStatus.AwaitingPayment]: 'Awaiting payment',
  [OrderStatus.Paid]: 'Paid',
  [OrderStatus.InProgress]: 'In progress',
  [OrderStatus.Shipped]: 'Shipped',
  [OrderStatus.Delivered]: 'Delivered',
  [OrderStatus.Cancelled]: 'Cancelled',
  [OrderStatus.Refunded]: 'Refunded',
}

export const ORDER_KIND_LABELS: Readonly<Record<OrderKind, string>> = {
  [OrderKind.Bespoke]: 'Bespoke',
  [OrderKind.ReadyToWear]: 'Ready to wear',
}

// Filter options, ordered the way the owner works through a day.
export const ORDER_STATUS_OPTIONS: readonly OrderStatus[] = [
  OrderStatus.PendingConfirmation,
  OrderStatus.AwaitingPayment,
  OrderStatus.Paid,
  OrderStatus.InProgress,
  OrderStatus.Shipped,
  OrderStatus.Delivered,
  OrderStatus.Cancelled,
  OrderStatus.Refunded,
]

export const ORDER_KIND_OPTIONS: readonly OrderKind[] = [
  OrderKind.Bespoke,
  OrderKind.ReadyToWear,
]

export const PRODUCT_SECTION_LABELS: Readonly<Record<ProductSection, string>> = {
  [ProductSection.Suits]: 'Suits',
  [ProductSection.Trad]: 'Trad attires',
  [ProductSection.CasualsMen]: 'Casuals · Men',
  [ProductSection.CasualsWomen]: 'Casuals · Women',
}

export const PRODUCT_SECTION_OPTIONS: readonly SelectOption<ProductSection | ''>[] = [
  { value: '', label: 'All sections' },
  ...Object.values(ProductSection).map((section) => ({
    value: section,
    label: PRODUCT_SECTION_LABELS[section],
  })),
]

export const PRODUCT_FORM_SECTION_OPTIONS: readonly SelectOption<ProductSection>[] =
  Object.values(ProductSection).map((section) => ({
    value: section,
    label: PRODUCT_SECTION_LABELS[section],
  }))

export const PRODUCT_CATEGORY_OPTIONS_BY_SECTION: Readonly<
  Record<ProductSection, readonly SelectOption<string>[]>
> = {
  [ProductSection.Suits]: [
    { value: 'Three-Piece', label: 'Three-Piece' },
    { value: 'Two-Piece', label: 'Two-Piece' },
    { value: 'Blazers', label: 'Blazers' },
  ],
  [ProductSection.Trad]: [
    { value: 'Agbada', label: 'Agbada' },
    { value: 'Aso-Oke', label: 'Aso-Oke' },
    { value: 'Kaftans', label: 'Kaftans' },
  ],
  [ProductSection.CasualsMen]: [
    { value: 'Shirts', label: 'Shirts' },
    { value: 'Trousers', label: 'Trousers' },
    { value: 'Knitwear', label: 'Knitwear' },
  ],
  [ProductSection.CasualsWomen]: [
    { value: 'Sets', label: 'Sets' },
    { value: 'Tops', label: 'Tops' },
  ],
}

export const PRODUCT_ACTIVE_LABELS: Readonly<Record<AdminProductActiveFilter, string>> = {
  '': 'All availability',
  'true': 'Visible',
  'false': 'Archived',
}

export const PRODUCT_ACTIVE_OPTIONS: readonly SelectOption<AdminProductActiveFilter>[] = [
  { value: '', label: PRODUCT_ACTIVE_LABELS[''] },
  { value: 'true', label: PRODUCT_ACTIVE_LABELS['true'] },
  { value: 'false', label: PRODUCT_ACTIVE_LABELS['false'] },
]

export const ADMIN_PRODUCT_FILTER_OPTIONS: readonly AdminFilterOption[] = [
  {
    key: 'section',
    label: ADMIN_COPY.productSectionLabel,
    type: AdminFilterType.Select,
    items: Object.values(ProductSection).map((section) => ({
      value: section,
      label: PRODUCT_SECTION_LABELS[section],
    })),
  },
  {
    key: 'active',
    label: ADMIN_COPY.productActiveLabel,
    type: AdminFilterType.Select,
    items: [
      { value: 'true', label: PRODUCT_ACTIVE_LABELS['true'] },
      { value: 'false', label: PRODUCT_ACTIVE_LABELS['false'] },
    ],
  },
]

export const ADMIN_ORDER_FILTER_OPTIONS: readonly AdminFilterOption[] = [
  {
    key: 'status',
    label: ADMIN_COPY.statusLabel,
    type: AdminFilterType.Select,
    items: ORDER_STATUS_OPTIONS.map((status) => ({
      value: status,
      label: ORDER_STATUS_LABELS[status],
    })),
  },
  {
    key: 'kind',
    label: ADMIN_COPY.kindLabel,
    type: AdminFilterType.Select,
    items: ORDER_KIND_OPTIONS.map((kind) => ({
      value: kind,
      label: ORDER_KIND_LABELS[kind],
    })),
  },
  { key: 'from', label: 'From date', type: AdminFilterType.Date },
  { key: 'to', label: 'To date', type: AdminFilterType.Date },
]

export const ADMIN_STATUS_ACTIONS: Readonly<Record<OrderStatus, readonly AdminStatusAction[]>> = {
  [OrderStatus.PendingConfirmation]: [
    { status: OrderStatus.AwaitingPayment, label: 'Send payment link' },
    { status: OrderStatus.Cancelled, label: 'Cancel order' },
  ],
  [OrderStatus.AwaitingPayment]: [
    { status: OrderStatus.Cancelled, label: 'Cancel order' },
  ],
  [OrderStatus.Paid]: [
    { status: OrderStatus.InProgress, label: 'Start fulfilment' },
    { status: OrderStatus.Cancelled, label: 'Cancel order' },
    { status: OrderStatus.Refunded, label: 'Mark refunded' },
  ],
  [OrderStatus.InProgress]: [
    { status: OrderStatus.Shipped, label: 'Mark shipped' },
    { status: OrderStatus.Cancelled, label: 'Cancel order' },
  ],
  [OrderStatus.Shipped]: [
    { status: OrderStatus.Delivered, label: 'Mark delivered' },
  ],
  [OrderStatus.Delivered]: [],
  [OrderStatus.Cancelled]: [],
  [OrderStatus.Refunded]: [],
}
