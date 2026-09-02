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
import { FULFILLMENT_STAGE_LABELS } from '~/constants/tracking'

export const ADMIN_ROUTES = {
  login: '/admin/login',
  orders: '/admin',
  delivery: '/admin/delivery',
  homeContent: '/admin/home-content',
  orderDetail: (orderNo: string) => `/admin/orders/${encodeURIComponent(orderNo)}`,
  products: '/admin/products',
  productNew: '/admin/products/new',
  productEdit: (slug: string) => `/admin/products/${encodeURIComponent(slug)}`,
  settings: '/admin/settings',
  team: '/admin/team',
} as const

export const ADMIN_API_ROUTES = {
  login: '/admin/login',
  logout: '/admin/logout',
  me: '/admin/me',
  settings: '/admin/settings',
  orders: '/admin/orders',
  orderMetrics: '/admin/orders/metrics',
  orderDetail: (orderNo: string) => `/admin/orders/${encodeURIComponent(orderNo)}`,
  orderStatus: (orderNo: string) => `/admin/orders/${encodeURIComponent(orderNo)}/status`,
  orderTracking: (orderNo: string) =>
    `/admin/orders/${encodeURIComponent(orderNo)}/tracking`,
  deliveryZones: '/admin/delivery-zones',
  deliveryZone: (name: string) => `/admin/delivery-zones/${encodeURIComponent(name)}`,
  homeContent: '/admin/home-content',
  team: '/admin/team',
  teamMember: (id: string) => `/admin/team/${encodeURIComponent(id)}`,
  password: '/admin/password',
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
  deliveryActiveLabel: 'Shown at checkout',
  deliveryVisibleStatus: 'Shown',
  deliveryHiddenStatus: 'Hidden',
  deliveryNewHeading: 'Add delivery area',
  deliveryNewNameLabel: 'Area name',
  deliveryNewFeeLabel: 'Fee in naira',
  deliveryCreate: 'Add area',
  deliveryCreating: 'Adding…',
  deliveryHide: 'Hide area',
  deliveryShow: 'Show area',
  deliverySave: 'Save',
  deliverySaving: 'Saving…',
  deliverySavedNotice: 'Delivery zone saved',
  deliveryCreatedNotice: 'Delivery zone added',
  deliverySaveFailed: 'Delivery zone could not be saved.',
  deliveryCreateFailed: 'Delivery zone could not be added.',
  closeDeliveryEditor: 'Close delivery editor',
  teamHeading: 'Team',
  teamSubheading: 'Manage who can sign in and what they can reach',
  teamLoading: 'Loading team…',
  teamLoadFailed: 'The team could not be loaded.',
  teamEmailLabel: 'Email',
  teamRoleLabel: 'Role',
  teamStatusLabel: 'Status',
  teamAddedLabel: 'Added',
  teamActiveStatus: 'Active',
  teamRevokedStatus: 'Revoked',
  teamRoleOwner: 'Owner',
  teamRoleStaff: 'Staff',
  teamRoleOwnerHint: 'Full access, including the team, settings, delivery and prices.',
  teamRoleStaffHint: 'Orders and fulfilment only.',
  teamInviteHeading: 'Add an admin',
  teamTemporaryPasswordLabel: 'Temporary password',
  teamTemporaryPasswordHint:
    'At least 12 characters. Share it privately and ask them to change it on their first sign-in.',
  teamCreate: 'Add admin',
  teamCreating: 'Adding…',
  teamCreatedNotice: 'Admin added',
  teamRevoke: 'Revoke access',
  teamRestore: 'Restore access',
  teamDelete: 'Delete',
  teamSavedNotice: 'Team updated',
  teamDeletedNotice: 'Admin deleted',
  teamYou: 'You',
  teamSelfHint: 'You cannot change your own access.',
  teamDeleteConfirm:
    'Delete this admin permanently? Revoking access instead keeps their name on the order history they touched.',
  passwordHeading: 'Your password',
  passwordIntro: 'Change the password you use to sign in.',
  currentPasswordLabel: 'Current password',
  newPasswordLabel: 'New password',
  confirmPasswordLabel: 'Confirm new password',
  passwordMismatch: 'The new passwords do not match.',
  changePassword: 'Change password',
  changingPassword: 'Changing…',
  passwordChangedNotice: 'Password changed',
  settingsHeading: 'Settings',
  settingsSubheading: 'Manage global storefront details',
  settingsLoading: 'Loading settings…',
  settingsLoadFailed: 'Settings could not be loaded.',
  settingsSavedNotice: 'Settings saved',
  settingsSaveFailed: 'Settings could not be saved.',
  saveSettings: 'Save settings',
  savingSettings: 'Saving…',
  announcementLabel: 'Announcement bar',
  announcementHint: 'Shown across the top of the storefront.',
  whatsappNumberLabel: 'WhatsApp number',
  whatsappNumberHint: 'Use international format, e.g. +2348094939264.',
  homeContentHeading: 'Home content',
  homeContentSubheading: 'Manage homepage copy and category tiles',
  homeContentLoading: 'Loading home content…',
  homeContentLoadFailed: 'Home content could not be loaded.',
  homeContentSavedNotice: 'Home content saved',
  homeContentSaveFailed: 'Home content could not be saved.',
  saveHomeContent: 'Save home content',
  savingHomeContent: 'Saving…',
  homeHeroPanel: 'Hero',
  homeHeroPanelIntro: 'The first message and action customers see on the homepage.',
  homeCategoriesPanel: 'Category tiles',
  homeCategoriesPanelIntro: 'Organise the collection links and imagery below the hero.',
  homeHeroLeadLabel: 'Hero heading line 1',
  homeHeroTailLabel: 'Hero heading line 2',
  homeHeroDescriptionLabel: 'Hero description',
  homeHeroActionLabel: 'CTA label',
  homeHeroActionToLabel: 'CTA link',
  homeCategoryEyebrowLabel: 'Section eyebrow',
  homeCategoryHeadingLabel: 'Section heading',
  homeTileLabelLabel: 'Tile label',
  homeTileHeadingLabel: 'Tile heading',
  homeTileImagePathLabel: 'Image URL/path',
  homeTileImageAltLabel: 'Image alt text',
  homeTileToLabel: 'Tile link',
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
  productIdentityHeading: 'Product identity',
  productIdentityIntro: 'The name, URL, and description customers see.',
  productCatalogueHeading: 'Catalogue placement',
  productCatalogueIntro: 'Control where the piece appears and how it is priced.',
  productAvailabilityHeading: 'Ordering & visibility',
  productAvailabilityIntro: 'Set its checkout type and storefront availability.',
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
  rowsPerPage: 'Rows per page',
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
  orderOverviewEyebrow: 'Overview',
  orderOverviewHeading: 'Customer & order state',
  orderOverviewIntro: 'Contact details and the manual actions available for this order.',
  orderProgressEyebrow: 'Fulfilment',
  orderProgressHeading: 'Progress & customer updates',
  orderProgressIntro: 'Review the timeline and publish the next customer-facing milestone.',
  orderCommercialEyebrow: 'Commercial',
  orderCommercialHeading: 'Items & payment',
  orderCommercialIntro: 'Review the order value, payment link, and reply shortcuts.',
  orderTailoringEyebrow: 'Tailoring',
  orderTailoringHeading: 'Fit information',
  orderTailoringIntro: 'Measurements supplied for made-to-measure production.',
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
  fulfillment: 'Order progress',
  fulfillmentLoading: 'Loading order progress…',
  fulfillmentUnavailable: 'Order progress starts after payment is confirmed.',
  fulfillmentLoadFailed: 'Order progress could not be loaded.',
  fulfillmentNoteLabel: 'Customer update',
  fulfillmentNoteHint: 'Optional. This appears on the customer tracking page and email.',
  estimatedDeliveryLabel: 'Estimated delivery date',
  carrierNameLabel: 'Courier name',
  carrierTrackingNumberLabel: 'Courier tracking number',
  carrierTrackingUrlLabel: 'Courier tracking link',
  signedForByLabel: 'Signed for by',
  signedAtLabel: 'Date and time received',
  deliveryLocationLabel: 'Delivered to',
  riderNameLabel: 'Rider',
  riderContactLabel: 'Rider contact number',
  deliveryWindowLabel: 'Delivery window',
  riderHint: 'Shown on the out-for-delivery email so the customer knows who to expect.',
  deliveryHandoverHint:
    'Optional. Shown on the delivery email as the record of who received the parcel.',
  advanceFulfillment: 'Move to',
  advancingFulfillment: 'Updating…',
  fulfillmentUpdatedNotice: 'Order progress updated',
  fulfillmentUpdateFailed: 'Order progress could not be updated.',
  copyTrackingLink: 'Copy tracking link',
  trackingLinkCopied: 'Tracking link copied',
  customerTrackingLink: 'Customer tracking link',
  settingsPanelHeading: 'Storefront essentials',
  settingsPanelIntro: 'Keep the global announcement and primary customer contact current.',
  deliveryDetailsHeading: 'Area details',
  deliveryDetailsIntro: 'Name the delivery area and set the fee shown at checkout.',
  deliveryVisibilityHeading: 'Checkout visibility',
  deliveryVisibilityIntro:
    'Hidden areas remain in admin but stop appearing at checkout and bespoke confirmation.',
  deliveryPreviewHeading: 'Customer preview',
} as const

export { FULFILLMENT_STAGE_LABELS }

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

// The API clamps limits at 100, so the UI exposes only supported page sizes.
export const ADMIN_ROWS_PER_PAGE_OPTIONS: readonly SelectOption<string>[] = [
  { value: '10', label: '10' },
  { value: '25', label: '25' },
  { value: '50', label: '50' },
  { value: '100', label: '100' },
]

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
      { id: 'home-content', label: 'Home content', route: ADMIN_ROUTES.homeContent },
      { id: 'delivery', label: 'Delivery', route: ADMIN_ROUTES.delivery },
    ],
  },
  {
    id: 'other',
    label: 'Other',
    items: [
      { id: 'settings', label: 'Settings', route: ADMIN_ROUTES.settings },
      // Owner-only. Hidden for staff by the sidebar, and refused by the API
      // regardless, so hiding it is convenience rather than the control.
      { id: 'team', label: 'Team', route: ADMIN_ROUTES.team, ownerOnly: true },
    ],
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
    { status: OrderStatus.Cancelled, label: 'Cancel order' },
    { status: OrderStatus.Refunded, label: 'Mark refunded' },
  ],
  [OrderStatus.InProgress]: [
    { status: OrderStatus.Cancelled, label: 'Cancel order' },
  ],
  [OrderStatus.Shipped]: [],
  [OrderStatus.Delivered]: [],
  [OrderStatus.Cancelled]: [],
  [OrderStatus.Refunded]: [],
}
