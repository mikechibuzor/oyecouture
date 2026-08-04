import { ADMIN_API_ROUTES } from '~/constants/admin'
import type {
  AdminProduct,
  AdminProductDetailResponse,
  AdminProductDraft,
  AdminProductPage,
  AdminProductQuery,
} from '~/types/admin'

// Reads and mutates owner-facing catalogue data through authenticated admin routes.
export const useAdminProducts = () => {
  const apiBaseUrl = useApiBaseUrl()

  const fetchProducts = (query: AdminProductQuery): Promise<AdminProductPage> =>
    $fetch<AdminProductPage>(ADMIN_API_ROUTES.products, {
      baseURL: apiBaseUrl,
      credentials: 'include',
      query: {
        ...(query.section ? { section: query.section } : {}),
        ...(query.active ? { active: query.active } : {}),
        ...(query.q ? { q: query.q } : {}),
        ...(query.limit ? { limit: query.limit } : {}),
        ...(query.offset ? { offset: query.offset } : {}),
      },
    })

  const fetchProduct = (slug: string): Promise<AdminProduct> =>
    $fetch<AdminProductDetailResponse>(ADMIN_API_ROUTES.productDetail(slug), {
      baseURL: apiBaseUrl,
      credentials: 'include',
    }).then((response) => response.product)

  const createProduct = (draft: AdminProductDraft): Promise<AdminProduct> =>
    $fetch<AdminProductDetailResponse>(ADMIN_API_ROUTES.products, {
      baseURL: apiBaseUrl,
      credentials: 'include',
      method: 'POST',
      body: draft,
    }).then((response) => response.product)

  const updateProduct = (
    slug: string,
    draft: AdminProductDraft,
  ): Promise<AdminProduct> =>
    $fetch<AdminProductDetailResponse>(ADMIN_API_ROUTES.productDetail(slug), {
      baseURL: apiBaseUrl,
      credentials: 'include',
      method: 'PUT',
      body: draft,
    }).then((response) => response.product)

  const updateProductActive = (slug: string, active: boolean): Promise<void> =>
    $fetch<void>(ADMIN_API_ROUTES.productActive(slug), {
      baseURL: apiBaseUrl,
      credentials: 'include',
      method: 'PATCH',
      body: { active },
    })

  return { fetchProducts, fetchProduct, createProduct, updateProduct, updateProductActive }
}
