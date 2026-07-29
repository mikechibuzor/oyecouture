import { PRODUCT_PRICE_FORMAT } from '~/constants/products'
import type { ProductFetchError } from '~/types/product'

// Formats integer kobo as Nigerian naira without introducing floating-point storage.
export const formatProductPrice = (priceKobo: number): string =>
  new Intl.NumberFormat(PRODUCT_PRICE_FORMAT.locale, {
    style: 'currency',
    currency: PRODUCT_PRICE_FORMAT.currency,
    maximumFractionDigits: 0,
  }).format(priceKobo / 100)

// Reads a stable API message when available and otherwise returns safe page copy.
export const getProductErrorMessage = (
  requestError: unknown,
  fallbackMessage: string,
): string => {
  const productFetchError = requestError as ProductFetchError | null

  return productFetchError?.data?.error.message ?? fallbackMessage
}
