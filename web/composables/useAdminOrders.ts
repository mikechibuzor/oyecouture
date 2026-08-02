import { ADMIN_API_ROUTES } from '~/constants/admin'
import type {
  AdminOrderDetail,
  AdminOrderMetrics,
  AdminOrderPage,
  AdminOrderQuery,
  MetricsPeriod,
  OrderStatus,
} from '~/types/admin'

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
        ...(query.from ? { from: query.from } : {}),
        ...(query.to ? { to: query.to } : {}),
        ...(query.limit ? { limit: query.limit } : {}),
        ...(query.offset ? { offset: query.offset } : {}),
      },
    })

  const fetchMetrics = (period: MetricsPeriod): Promise<AdminOrderMetrics> =>
    $fetch<AdminOrderMetrics>(ADMIN_API_ROUTES.orderMetrics, {
      baseURL: apiBaseUrl,
      credentials: 'include',
      query: { period },
    })

  const fetchOrder = (orderNo: string): Promise<AdminOrderDetail> =>
    $fetch<AdminOrderDetail>(ADMIN_API_ROUTES.orderDetail(orderNo), {
      baseURL: apiBaseUrl,
      credentials: 'include',
    })

  const updateOrderStatus = (
    orderNo: string,
    status: OrderStatus,
  ): Promise<AdminOrderDetail> =>
    $fetch<AdminOrderDetail>(ADMIN_API_ROUTES.orderStatus(orderNo), {
      baseURL: apiBaseUrl,
      credentials: 'include',
      method: 'PATCH',
      body: { status },
    })

  return { fetchOrders, fetchMetrics, fetchOrder, updateOrderStatus }
}
