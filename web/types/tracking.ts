import type { OrderKind } from '~/types/order'

export enum FulfillmentStage {
  PaymentConfirmed = 'payment_confirmed',
  MeasurementsApproved = 'measurements_approved',
  PreparingOrder = 'preparing_order',
  InProduction = 'in_production',
  QualityCheck = 'quality_check',
  ReadyForDispatch = 'ready_for_dispatch',
  Dispatched = 'dispatched',
  // Optional, and only for a local delivery by Oyecouture's own rider.
  OutForDelivery = 'out_for_delivery',
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
  // Handover record, accepted by the API only on the delivered transition.
  readonly signedForBy: string
  readonly signedAt: string
  readonly deliveryLocation: string
  // Rider details, accepted only on the out-for-delivery transition.
  readonly riderName: string
  readonly riderContact: string
  readonly deliveryWindow: string
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
