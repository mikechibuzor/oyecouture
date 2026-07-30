export enum OrderKind {
  Bespoke = 'bespoke',
  ReadyToWear = 'ready-to-wear',
}

// Only slug, size and quantity are sent. The server re-prices every line from
// the database, so the storefront can never influence what is charged.
export interface OrderItemRequest {
  readonly slug: string
  readonly size: string
  readonly quantity: number
}

export interface MeasurementsRequest {
  neckCm: number | null
  chestCm: number | null
  waistCm: number | null
  hipsCm: number | null
  sleeveCm: number | null
  inseamCm: number | null
  lengthCm: number | null
  heightCm: number | null
  weightKg: number | null
  notes: string
}

export interface DeliveryZone {
  readonly name: string
  readonly feeKobo: number
}

export interface DeliveryZonesResponse {
  readonly zones: readonly DeliveryZone[]
}

export interface CreateOrderRequest {
  readonly customerName: string
  readonly email: string
  readonly phone: string
  readonly items: readonly OrderItemRequest[]
  readonly measurements: MeasurementsRequest | null
  // Ready-to-wear only; bespoke delivery is agreed on WhatsApp.
  readonly address?: string
  readonly deliveryZone?: string
}

export interface OrderItem {
  readonly slug: string
  readonly name: string
  readonly size?: string
  readonly quantity: number
  readonly priceKobo: number
  readonly lineKobo: number
}

export interface Order {
  readonly orderNo: string
  readonly kind: OrderKind
  readonly status: string
  readonly customerName: string
  readonly items: readonly OrderItem[]
  readonly subtotalKobo: number
  readonly deliveryZone?: string
  readonly deliveryFeeKobo: number
  readonly totalKobo: number
  readonly payToken: string
  readonly createdAt: string
}

export interface CreateOrderResponse {
  readonly order: Order
}
