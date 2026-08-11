import { ADMIN_API_ROUTES } from '~/constants/admin'
import type {
  AdminDeliveryZone,
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

  return { fetchDeliveryZones, updateDeliveryZone }
}
