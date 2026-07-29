import type { MaybeRefOrGetter } from 'vue'
import {
  PRODUCT_FALLBACK_IMAGES,
  PRODUCT_MEDIA_PLACEHOLDER_COUNT,
} from '~/constants/products'
import type { Product, ProductMediaItem } from '~/types/product'

// Resolves API images first, then approved design-source fallbacks for missing photography.
export const useProductMedia = (
  product: MaybeRefOrGetter<Product>,
  fillGallerySlots = false,
) => {
  // Builds a stable media list and fills missing gallery positions with placeholders.
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

    if (!fillGallerySlots) {
      return resolvedMediaItems.length > 0
        ? resolvedMediaItems
        : [{
            id: `${resolvedProduct.slug}-placeholder`,
            source: null,
            alt: resolvedProduct.name,
          }]
    }

    const galleryMediaItems: ProductMediaItem[] = [...resolvedMediaItems]

    while (galleryMediaItems.length < PRODUCT_MEDIA_PLACEHOLDER_COUNT) {
      galleryMediaItems.push({
        id: `${resolvedProduct.slug}-placeholder-${galleryMediaItems.length}`,
        source: null,
        alt: resolvedProduct.name,
      })
    }

    return galleryMediaItems
  })

  return { mediaItems }
}
