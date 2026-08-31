import type { OrderKind } from '~/types/order'

export enum FulfillmentStage {
  PaymentConfirmed = 'payment_confirmed',
  MeasurementsApproved = 'measurements_approved',
  PreparingOrder = 'preparing_order',
  InProduction = 'in_production',
  QualityCheck = 'quality_check',
  ReadyForDispatch = 'ready_for_dispatch',
  Dispatched = 'dispatched',
  Delivered = 'delivered',
}

export interface FulfillmentEvent {
  readonly id: string
  readonly stage: FulfillmentStage
  readonly label: string
  readonly publicNote?: string
  readonly actorEmail?: string
  readonly occurredAt: string
}

export interface TrackingShipment {
  readonly carrierName?: string
  readonly carrierTrackingNumber?: string
  readonly carrierTrackingUrl?: string
}

export interface AdminOrderTracking {
  readonly orderNo: string
  readonly kind: OrderKind
  readonly currentStage: FulfillmentStage
  readonly currentStageLabel: string
  readonly version: number
  readonly estimatedDeliveryDate?: string
  readonly shipment: TrackingShipment
  readonly trackingUrl?: string
  readonly events: readonly FulfillmentEvent[]
  readonly allowedNextStages: readonly FulfillmentStage[]
}

export interface FulfillmentTransitionInput {
  readonly stage: FulfillmentStage
  readonly expectedVersion: number
  readonly publicNote: string
  readonly estimatedDeliveryDate: string
  readonly carrierName: string
  readonly carrierTrackingNumber: string
  readonly carrierTrackingUrl: string
}

export interface TrackingItem {
  readonly name: string
  readonly size?: string
  readonly quantity: number
}

export interface PublicOrderTracking {
  readonly orderNo: string
  readonly kind: OrderKind
  readonly currentStage: FulfillmentStage
  readonly currentStageLabel: string
  readonly estimatedDeliveryDate?: string
  readonly shipment: TrackingShipment
  readonly items: readonly TrackingItem[]
  readonly events: readonly FulfillmentEvent[]
}
