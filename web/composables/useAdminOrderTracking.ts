import { ADMIN_API_ROUTES } from '~/constants/admin'
import type {
  AdminOrderTracking,
  FulfillmentTransitionInput,
} from '~/types/tracking'

export const useAdminOrderTracking = () => {
  const apiBaseUrl = useApiBaseUrl()

  const fetchOrderTracking = (orderNo: string): Promise<AdminOrderTracking> =>
    $fetch<AdminOrderTracking>(ADMIN_API_ROUTES.orderTracking(orderNo), {
      baseURL: apiBaseUrl,
      credentials: 'include',
    })

  const advanceOrderTracking = (
    orderNo: string,
    transition: FulfillmentTransitionInput,
  ): Promise<AdminOrderTracking> =>
    $fetch<AdminOrderTracking>(ADMIN_API_ROUTES.orderTracking(orderNo), {
      baseURL: apiBaseUrl,
      credentials: 'include',
      method: 'PATCH',
      body: transition,
    })

  return { advanceOrderTracking, fetchOrderTracking }
}
