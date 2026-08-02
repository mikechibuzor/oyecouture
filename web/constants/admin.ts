import { OrderKind } from '~/types/order'
import { OrderStatus } from '~/types/admin'

export const ADMIN_ROUTES = {
  login: '/admin/login',
  orders: '/admin',
} as const

export const ADMIN_API_ROUTES = {
  login: '/admin/login',
  logout: '/admin/logout',
  me: '/admin/me',
  orders: '/admin/orders',
} as const

export const ADMIN_COPY = {
  brand: 'Oyecouture',
  panel: 'Order desk',
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
} as const

// The page size the list requests; the API clamps anything larger.
export const ADMIN_ORDERS_PAGE_SIZE = 25

// Typing in the search box should not fire a request per keystroke.
export const ADMIN_SEARCH_DEBOUNCE_MS = 250

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
