import type { MaybeRefOrGetter } from 'vue'
import { PRODUCT_FALLBACK_IMAGES } from '~/constants/products'
import type { Product, ProductMediaItem } from '~/types/product'

// Resolves API images first, then approved design-source fallbacks for missing photography.
export const useProductMedia = (product: MaybeRefOrGetter<Product>) => {
  // Builds a stable media list without manufacturing carousel positions.
  const mediaItems = computed<readonly ProductMediaItem[]>(() => {
    const resolvedProduct = toValue(product)
    const apiMediaItems = resolvedProduct.images.map((productImage) => ({
      id: productImage.id,
      source: productImage.url,
      alt: resolvedProduct.name,
    }))
    const fallbackMediaItems = (
      PRODUCT_FALLBACK_IMAGES[resolvedProduct.slug] ?? []
    ).map((imageSource, imageIndex) => ({
      id: `${resolvedProduct.slug}-fallback-${imageIndex}`,
      source: imageSource,
      alt: resolvedProduct.name,
    }))
    const resolvedMediaItems = apiMediaItems.length > 0
      ? apiMediaItems
      : fallbackMediaItems

    return resolvedMediaItems.length > 0
      ? resolvedMediaItems
      : [{
          id: `${resolvedProduct.slug}-placeholder`,
          source: null,
          alt: resolvedProduct.name,
        }]
  })

  return { mediaItems }
}
