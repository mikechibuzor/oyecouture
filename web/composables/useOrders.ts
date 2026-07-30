import { API_ROUTES } from '~/constants/api'
import type { CreateOrderRequest, CreateOrderResponse } from '~/types/order'

// The API base must be resolved during setup: Nuxt composables like
// useRuntimeConfig() throw when called later from an async event handler.
export const useOrders = () => {
  const apiBaseUrl = useApiBaseUrl()

  const submitOrder = (
    request: CreateOrderRequest,
  ): Promise<CreateOrderResponse> =>
    $fetch<CreateOrderResponse>(API_ROUTES.orders, {
      baseURL: apiBaseUrl,
      method: 'POST',
      body: request,
    })

  return { submitOrder }
}
