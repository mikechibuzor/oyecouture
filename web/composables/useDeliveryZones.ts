import { API_ROUTES } from '~/constants/api'
import type { DeliveryZonesResponse } from '~/types/order'

// Delivery zones and their fees live in the database so the owner can change
// them without a deploy.
export const useDeliveryZones = () =>
  useFetch<DeliveryZonesResponse>(API_ROUTES.deliveryZones, {
    baseURL: useApiBaseUrl(),
    key: 'delivery-zones',
    retry: 0,
  })
