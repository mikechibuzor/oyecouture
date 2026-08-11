import { ADMIN_API_ROUTES } from '~/constants/admin'
import type {
  AdminDeliveryZone,
  AdminDeliveryZoneCreate,
  AdminDeliveryZoneResponse,
  AdminDeliveryZonesResponse,
  AdminDeliveryZoneUpdate,
} from '~/types/admin'

export const useAdminDeliveryZones = () => {
  const apiBaseUrl = useApiBaseUrl()

  const fetchDeliveryZones = (): Promise<AdminDeliveryZonesResponse> =>
    $fetch<AdminDeliveryZonesResponse>(ADMIN_API_ROUTES.deliveryZones, {
      baseURL: apiBaseUrl,
      credentials: 'include',
    })

  const createDeliveryZone = (create: AdminDeliveryZoneCreate): Promise<AdminDeliveryZone> =>
    $fetch<AdminDeliveryZoneResponse>(ADMIN_API_ROUTES.deliveryZones, {
      baseURL: apiBaseUrl,
      body: create,
      credentials: 'include',
      method: 'POST',
    }).then((response) => response.zone)

  const updateDeliveryZone = (
    name: string,
    update: AdminDeliveryZoneUpdate,
  ): Promise<AdminDeliveryZone> =>
    $fetch<AdminDeliveryZoneResponse>(ADMIN_API_ROUTES.deliveryZone(name), {
      baseURL: apiBaseUrl,
      body: update,
      credentials: 'include',
      method: 'PATCH',
    }).then((response) => response.zone)

  return { createDeliveryZone, fetchDeliveryZones, updateDeliveryZone }
}
