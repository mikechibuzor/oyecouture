import { API_ROUTES, ASYNC_DATA_KEYS } from '~/constants/api'
import type { ApiHealthResponse } from '~/types/health'

export const useApiHealth = () => {

  return useFetch<ApiHealthResponse>(API_ROUTES.health, {
    baseURL: useApiBaseUrl(),
    key: ASYNC_DATA_KEYS.apiHealth,
    retry: 0,
  })
}
