import { API_ROUTES } from '~/constants/api'
import {
  PRODUCT_API_QUERY_KEYS,
  PRODUCT_ASYNC_DATA_KEYS,
} from '~/constants/products'
import type {
  ProductDetailResponse,
  ProductListResponse,
} from '~/types/product'
import type { ProductSection } from '~/types/product'

// Loads the public product catalogue with an optional validated section filter.
export const useProductList = (productSection?: ProductSection) => {
  const requestKey = productSection
    ? `${PRODUCT_ASYNC_DATA_KEYS.sectionPrefix}-${productSection}`
    : PRODUCT_ASYNC_DATA_KEYS.all

  return useFetch<ProductListResponse>(API_ROUTES.products, {
    baseURL: useApiBaseUrl(),
    key: requestKey,
    query: productSection
      ? { [PRODUCT_API_QUERY_KEYS.section]: productSection }
      : undefined,
    retry: 0,
  })
}

// Loads one public product using its canonical slug.
export const useProductBySlug = (productSlug: string) => {
  const encodedProductSlug = encodeURIComponent(productSlug)

  return useFetch<ProductDetailResponse>(
    `${API_ROUTES.products}/${encodedProductSlug}`,
    {
      baseURL: useApiBaseUrl(),
      key: `${PRODUCT_ASYNC_DATA_KEYS.detailPrefix}-${encodedProductSlug}`,
      retry: 0,
    },
  )
}
