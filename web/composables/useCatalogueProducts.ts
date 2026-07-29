import type { MaybeRefOrGetter } from 'vue'
import {
  CASUAL_PRODUCT_SECTIONS,
  PRODUCT_COPY,
  PRODUCT_FILTER_CATEGORIES,
} from '~/constants/products'
import type { Product } from '~/types/product'
import {
  CatalogueFilterId,
  ProductSection,
  ProductSortId,
} from '~/types/product'

const casualProductSectionSet: ReadonlySet<ProductSection> = new Set(
  CASUAL_PRODUCT_SECTIONS,
)

// Coordinates client-side catalogue filtering, sorting, and result counts.
export const useCatalogueProducts = (
  products: MaybeRefOrGetter<readonly Product[]>,
) => {
  const activeFilterId = ref(CatalogueFilterId.All)
  const activeSortId = ref(ProductSortId.Featured)

  // Reports whether a product belongs in the selected catalogue filter.
  const productMatchesFilter = (product: Product): boolean => {
    switch (activeFilterId.value) {
      case CatalogueFilterId.Suits:
        return product.section === ProductSection.Suits
      case CatalogueFilterId.AsoOke:
        return product.category === PRODUCT_FILTER_CATEGORIES.asoOke
      case CatalogueFilterId.Casuals:
        return casualProductSectionSet.has(product.section)
      case CatalogueFilterId.Kaftans:
        return product.category === PRODUCT_FILTER_CATEGORIES.kaftans
      case CatalogueFilterId.All:
      default:
        return true
    }
  }

  // Orders two products according to the selected catalogue sort.
  const compareProducts = (firstProduct: Product, secondProduct: Product): number => {
    switch (activeSortId.value) {
      case ProductSortId.PriceAscending:
        return firstProduct.priceKobo - secondProduct.priceKobo
      case ProductSortId.PriceDescending:
        return secondProduct.priceKobo - firstProduct.priceKobo
      case ProductSortId.Newest:
        return Date.parse(secondProduct.createdAt) - Date.parse(firstProduct.createdAt)
      case ProductSortId.Featured:
      default:
        return firstProduct.sortOrder - secondProduct.sortOrder
    }
  }

  // Returns the visible products after applying the active filter and sort.
  const visibleProducts = computed(() =>
    [...toValue(products)]
      .filter(productMatchesFilter)
      .sort(compareProducts),
  )

  // Formats the visible result count with the correct noun.
  const pieceCountLabel = computed(() => {
    const productCount = visibleProducts.value.length
    const pieceLabel = productCount === 1
      ? PRODUCT_COPY.catalogue.pieceSingular
      : PRODUCT_COPY.catalogue.piecePlural

    return `${productCount} ${pieceLabel}`
  })

  return {
    activeFilterId,
    activeSortId,
    pieceCountLabel,
    visibleProducts,
  }
}
