import { ADMIN_API_ROUTES } from '~/constants/admin'
import type {
  AdminProduct,
  AdminProductDetailResponse,
  AdminProductDraft,
  AdminProductImagesResponse,
  AdminProductPage,
  AdminProductPayload,
  AdminProductQuery,
} from '~/types/admin'
import { nairaToKobo } from '~/utils/product'

// Reads and mutates owner-facing catalogue data through authenticated admin routes.
export const useAdminProducts = () => {
  const apiBaseUrl = useApiBaseUrl()

  const toProductPayload = (draft: AdminProductDraft): AdminProductPayload => ({
    slug: draft.slug,
    name: draft.name,
    section: draft.section,
    category: draft.category,
    priceKobo: nairaToKobo(draft.priceNaira),
    description: draft.description,
    tag: draft.tag,
    hasSizes: draft.hasSizes,
    active: draft.active,
    sortOrder: draft.sortOrder,
  })

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
      body: toProductPayload(draft),
    }).then((response) => response.product)

  const updateProduct = (
    slug: string,
    draft: AdminProductDraft,
  ): Promise<AdminProduct> =>
    $fetch<AdminProductDetailResponse>(ADMIN_API_ROUTES.productDetail(slug), {
      baseURL: apiBaseUrl,
      credentials: 'include',
      method: 'PUT',
      body: toProductPayload(draft),
    }).then((response) => response.product)

  const updateProductActive = (slug: string, active: boolean): Promise<void> =>
    $fetch<void>(ADMIN_API_ROUTES.productActive(slug), {
      baseURL: apiBaseUrl,
      credentials: 'include',
      method: 'PATCH',
      body: { active },
    })

  const fetchProductImages = (slug: string): Promise<AdminProductImagesResponse> =>
    $fetch<AdminProductImagesResponse>(ADMIN_API_ROUTES.productImages(slug), {
      baseURL: apiBaseUrl,
      credentials: 'include',
    })

  const uploadProductImage = (slug: string, imageFile: File): Promise<AdminProduct> => {
    const formData = new FormData()
    formData.append('image', imageFile)

    return $fetch<AdminProductDetailResponse>(ADMIN_API_ROUTES.productImages(slug), {
      baseURL: apiBaseUrl,
      body: formData,
      credentials: 'include',
      method: 'POST',
    }).then((response) => response.product)
  }

  const reorderProductImages = (
    slug: string,
    imageIds: readonly string[],
  ): Promise<AdminProduct> =>
    $fetch<AdminProductDetailResponse>(ADMIN_API_ROUTES.productImagesReorder(slug), {
      baseURL: apiBaseUrl,
      body: { imageIds },
      credentials: 'include',
      method: 'PATCH',
    }).then((response) => response.product)

  const deleteProductImage = (slug: string, imageId: string): Promise<AdminProduct> =>
    $fetch<AdminProductDetailResponse>(ADMIN_API_ROUTES.productImage(slug, imageId), {
      baseURL: apiBaseUrl,
      credentials: 'include',
      method: 'DELETE',
    }).then((response) => response.product)

  return {
    fetchProducts,
    fetchProduct,
    fetchProductImages,
    createProduct,
    updateProduct,
    updateProductActive,
    uploadProductImage,
    reorderProductImages,
    deleteProductImage,
  }
}
