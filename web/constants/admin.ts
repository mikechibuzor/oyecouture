import { OrderKind } from '~/types/order'
import type {
  AdminFilterOption,
  AdminMetricCard,
  AdminNavGroup,
  AdminStatusAction,
} from '~/types/admin'
import { AdminFilterType, MetricsPeriod, OrderStatus } from '~/types/admin'

export const ADMIN_ROUTES = {
  login: '/admin/login',
  orders: '/admin',
  orderDetail: (orderNo: string) => `/admin/orders/${encodeURIComponent(orderNo)}`,
} as const

export const ADMIN_API_ROUTES = {
  login: '/admin/login',
  logout: '/admin/logout',
  me: '/admin/me',
  orders: '/admin/orders',
  orderMetrics: '/admin/orders/metrics',
  orderDetail: (orderNo: string) => `/admin/orders/${encodeURIComponent(orderNo)}`,
  orderStatus: (orderNo: string) => `/admin/orders/${encodeURIComponent(orderNo)}/status`,
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
  searchLabel: 'Search orders',
  searchPlaceholder: 'Order number, name, email or phone',
  statusLabel: 'Status',
  kindLabel: 'Type',
  anyStatus: 'All statuses',
  anyKind: 'All types',
  loading: 'Loading orders…',
  loadFailed: 'The orders could not be loaded.',
  retry: 'Try again',
  empty: 'No orders yet.',
  emptyFiltered: 'No orders match these filters.',
  clearFilters: 'Clear filters',
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

// The shape of the finished admin. Only Orders is built; the rest are marked
// pending rather than linking nowhere.
export const ADMIN_NAVIGATION: readonly AdminNavGroup[] = [
  {
    id: 'main',
    label: 'Main',
    items: [
      { id: 'orders', label: 'Orders', route: ADMIN_ROUTES.orders },
      { id: 'products', label: 'Products', isPending: true },
      { id: 'home-content', label: 'Home content', isPending: true },
      { id: 'delivery', label: 'Delivery', isPending: true },
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
