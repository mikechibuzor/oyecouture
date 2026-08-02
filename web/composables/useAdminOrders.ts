import { ADMIN_API_ROUTES } from '~/constants/admin'
import type { AdminOrderPage, AdminOrderQuery } from '~/types/admin'

// Reads the admin orders list. The API clamps paging and ignores unknown
// filter values, so the query is passed through as typed.
export const useAdminOrders = () => {
  const apiBaseUrl = useApiBaseUrl()

  const fetchOrders = (query: AdminOrderQuery): Promise<AdminOrderPage> =>
    $fetch<AdminOrderPage>(ADMIN_API_ROUTES.orders, {
      baseURL: apiBaseUrl,
      credentials: 'include',
      // Empty filters are dropped so the URL stays readable.
      query: {
        ...(query.status ? { status: query.status } : {}),
        ...(query.kind ? { kind: query.kind } : {}),
        ...(query.q ? { q: query.q } : {}),
        ...(query.limit ? { limit: query.limit } : {}),
        ...(query.offset ? { offset: query.offset } : {}),
      },
    })

  return { fetchOrders }
}
