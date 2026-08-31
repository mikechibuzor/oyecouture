import { API_ROUTES } from '~/constants/api'
import type { PublicOrderTracking } from '~/types/tracking'

export const useOrderTracking = () => {
  const apiBaseUrl = useApiBaseUrl()

  const fetchOrderTracking = (trackingToken: string): Promise<PublicOrderTracking> =>
    $fetch<PublicOrderTracking>(API_ROUTES.tracking(trackingToken), {
      baseURL: apiBaseUrl,
    })

  return { fetchOrderTracking }
}
