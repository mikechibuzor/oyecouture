import { OrderKind } from '~/types/order'
import { FulfillmentStage } from '~/types/tracking'

export const FULFILLMENT_STAGE_LABELS: Readonly<Record<FulfillmentStage, string>> = {
  [FulfillmentStage.PaymentConfirmed]: 'Payment confirmed',
  [FulfillmentStage.MeasurementsApproved]: 'Measurements approved',
  [FulfillmentStage.PreparingOrder]: 'Preparing order',
  [FulfillmentStage.InProduction]: 'In production',
  [FulfillmentStage.QualityCheck]: 'Quality check',
  [FulfillmentStage.ReadyForDispatch]: 'Ready for dispatch',
  [FulfillmentStage.Dispatched]: 'Dispatched',
  [FulfillmentStage.Delivered]: 'Delivered',
}

export const FULFILLMENT_PATHS: Readonly<Record<OrderKind, readonly FulfillmentStage[]>> = {
  [OrderKind.Bespoke]: [
    FulfillmentStage.PaymentConfirmed,
    FulfillmentStage.MeasurementsApproved,
    FulfillmentStage.InProduction,
    FulfillmentStage.QualityCheck,
    FulfillmentStage.ReadyForDispatch,
    FulfillmentStage.Dispatched,
    FulfillmentStage.Delivered,
  ],
  [OrderKind.ReadyToWear]: [
    FulfillmentStage.PaymentConfirmed,
    FulfillmentStage.PreparingOrder,
    FulfillmentStage.QualityCheck,
    FulfillmentStage.ReadyForDispatch,
    FulfillmentStage.Dispatched,
    FulfillmentStage.Delivered,
  ],
}

export const TRACKING_COPY = {
  eyebrow: 'Order progress',
  heading: 'Track your order',
  loading: 'Loading your order progress…',
  loadFailed: 'This tracking link could not be loaded.',
  retry: 'Try again',
  currentStatus: 'Current status',
  estimatedDelivery: 'Estimated delivery',
  shipment: 'Delivery details',
  carrier: 'Courier',
  trackingNumber: 'Tracking number',
  courierLink: 'Open courier tracking',
  items: 'Your pieces',
  awaitingUpdate: 'We will update this step as your order progresses.',
} as const
