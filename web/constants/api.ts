export const API_ROUTES = {
  contact: '/contact',
  deliveryZones: '/delivery-zones',
  health: '/healthz',
  orders: '/orders',
  payment: (payToken: string) => `/payments/${encodeURIComponent(payToken)}`,
  paymentInitialize: (payToken: string) =>
    `/payments/${encodeURIComponent(payToken)}/initialize`,
  paymentVerify: (payToken: string) => `/payments/${encodeURIComponent(payToken)}/verify`,
  products: '/products',
} as const

export const ASYNC_DATA_KEYS = {
  apiHealth: 'api-health',
} as const
