import type { OrderKind } from '~/types/order'

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

export interface AdminOrderQuery {
  readonly status?: OrderStatus | ''
  readonly kind?: OrderKind | ''
  readonly q?: string
  readonly limit?: number
  readonly offset?: number
}
