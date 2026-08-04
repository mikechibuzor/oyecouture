import type { MeasurementsRequest, OrderItem, OrderKind } from '~/types/order'
import type { Product, ProductImage, ProductSection } from '~/types/product'

// The order lifecycle, mirroring the database CHECK constraint.
export enum OrderStatus {
  PendingConfirmation = 'pending_confirmation',
  AwaitingPayment = 'awaiting_payment',
  Paid = 'paid',
  InProgress = 'in_progress',
  Shipped = 'shipped',
  Delivered = 'delivered',
  Cancelled = 'cancelled',
  Refunded = 'refunded',
}

export interface AdminSession {
  readonly email: string
  readonly role: string
}

// One row of the orders list. Contact details are admin-only.
export interface AdminOrderSummary {
  readonly orderNo: string
  readonly kind: OrderKind
  readonly status: OrderStatus
  readonly customerName: string
  readonly email: string
  readonly phone: string
  readonly itemCount: number
  readonly totalKobo: number
  readonly createdAt: string
}

export interface AdminOrderPage {
  readonly orders: readonly AdminOrderSummary[]
  readonly total: number
  readonly limit: number
  readonly offset: number
}

export interface AdminOrderDetail {
  readonly orderNo: string
  readonly kind: OrderKind
  readonly status: OrderStatus
  readonly customerName: string
  readonly email: string
  readonly phone: string
  readonly address?: string
  readonly deliveryZone?: string
  readonly deliveryFeeKobo: number
  readonly subtotalKobo: number
  readonly totalKobo: number
  readonly payToken: string
  readonly paymentProvider?: string
  readonly paymentRef?: string
  readonly items: readonly OrderItem[]
  readonly measurements?: MeasurementsRequest
  readonly createdAt: string
}

export interface AdminStatusAction {
  readonly status: OrderStatus
  readonly label: string
}

// The window the summary strip covers. It filters the metrics only — the list
// below always shows every order matching the current filters.
export enum MetricsPeriod {
  Today = 'today',
  LastWeek = 'last-week',
  Last14Days = 'last-14-days',
  Last60Days = 'last-60-days',
  Last90Days = 'last-90-days',
  AllTime = 'all-time',
}

// The summary strip above the orders list. Revenue counts only orders that
// have actually been paid for.
export interface AdminOrderMetrics {
  readonly total: number
  readonly pendingConfirmation: number
  readonly awaitingPayment: number
  readonly paid: number
  readonly revenueKobo: number
}

// A metric card. `statusFilter` is what clicking it applies to the list;
// `null` means the card clears the status filter instead.
export interface AdminMetricCard {
  readonly id: string
  readonly label: string
  readonly statusFilter: OrderStatus | null
  readonly isMoney?: boolean
}

// One entry in the side navigation. Sections not yet built are marked pending
// so the shell shows the shape of the admin without offering dead links.
export interface AdminNavItem {
  readonly id: string
  readonly label: string
  readonly route?: string
  readonly isPending?: boolean
}

export interface AdminNavGroup {
  readonly id: string
  readonly label: string
  readonly items: readonly AdminNavItem[]
}

export enum AdminFilterType {
  Select = 'select',
  Date = 'date',
}

export type AdminOrderFilterKey = 'status' | 'kind' | 'from' | 'to'
export type AdminProductFilterKey = 'section' | 'active'
export type AdminFilterKey = AdminOrderFilterKey | AdminProductFilterKey

export interface AdminFilterItem {
  readonly value: string
  readonly label: string
}

export interface AdminFilterOption {
  readonly key: AdminFilterKey
  readonly label: string
  readonly type: AdminFilterType
  readonly items?: readonly AdminFilterItem[]
}

export type AdminOrderFilterValues = Readonly<Record<AdminOrderFilterKey, string>>
export type AdminProductFilterValues = Readonly<Record<AdminProductFilterKey, string>>
export type AdminFilterValues = Readonly<Record<string, string>>

export interface AdminOrderQuery {
  readonly status?: OrderStatus | ''
  readonly kind?: OrderKind | ''
  readonly q?: string
  readonly from?: string
  readonly to?: string
  readonly limit?: number
  readonly offset?: number
}

export interface AdminProduct extends Product {
  readonly active: boolean
}

export interface AdminProductPage {
  readonly products: readonly AdminProduct[]
  readonly total: number
  readonly limit: number
  readonly offset: number
}

export interface AdminProductDetailResponse {
  readonly product: AdminProduct
}

export interface AdminProductImagesResponse {
  readonly images: readonly ProductImage[]
}

export interface AdminProductDraft {
  slug: string
  name: string
  section: ProductSection
  category: string
  priceKobo: number
  description: string
  tag: string
  hasSizes: boolean
  active: boolean
  sortOrder: number
}

export type AdminProductActiveFilter = '' | 'true' | 'false'

export interface AdminProductQuery {
  readonly section?: ProductSection | ''
  readonly active?: AdminProductActiveFilter
  readonly q?: string
  readonly limit?: number
  readonly offset?: number
}
