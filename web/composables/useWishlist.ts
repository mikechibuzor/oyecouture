import type { Product } from '~/types/product'

const WISHLIST_STORAGE_KEY = 'oyecouture-wishlist'

interface WishlistItem {
  readonly slug: string
  readonly product?: Product
}

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null

const isProductSnapshot = (value: unknown): value is Product =>
  isRecord(value)
  && typeof value.id === 'string'
  && typeof value.slug === 'string'
  && typeof value.name === 'string'
  && typeof value.category === 'string'
  && typeof value.priceKobo === 'number'
  && typeof value.description === 'string'
  && typeof value.tag === 'string'
  && typeof value.hasSizes === 'boolean'
  && typeof value.sortOrder === 'number'
  && typeof value.createdAt === 'string'
  && Array.isArray(value.images)

const toWishlistItem = (storedItem: unknown): WishlistItem | null => {
  if (typeof storedItem === 'string') {
    return { slug: storedItem }
  }

  if (!isRecord(storedItem)) {
    return null
  }

  if (isProductSnapshot(storedItem)) {
    return { slug: storedItem.slug, product: storedItem }
  }

  if (typeof storedItem.slug !== 'string') {
    return null
  }

  return {
    slug: storedItem.slug,
    product: isProductSnapshot(storedItem.product) ? storedItem.product : undefined,
  }
}

const readStoredWishlistItems = (): WishlistItem[] => {
  if (!import.meta.client) {
    return []
  }

  try {
    const storedValue = window.localStorage.getItem(WISHLIST_STORAGE_KEY)
    const parsedValue: unknown = storedValue ? JSON.parse(storedValue) : []

    return Array.isArray(parsedValue)
      ? parsedValue
          .map(toWishlistItem)
          .filter((item): item is WishlistItem => Boolean(item))
      : []
  } catch {
    return []
  }
}

export const useWishlist = () => {
  const items = useState<WishlistItem[]>('wishlist-items', () => [])
  const hasLoaded = useState<boolean>('wishlist-loaded', () => false)
  const hasBoundClientSync = useState<boolean>('wishlist-client-sync-bound', () => false)

  const loadStoredItems = (): void => {
    if (!import.meta.client || hasLoaded.value) {
      return
    }

    const storedItems = readStoredWishlistItems()

    if (items.value.length === 0 || storedItems.length > 0) {
      items.value = storedItems
    }

    hasLoaded.value = true
  }

  const syncStoredItems = (): void => {
    if (!import.meta.client) {
      return
    }

    items.value = readStoredWishlistItems()
    hasLoaded.value = true
  }

  if (import.meta.client && !hasBoundClientSync.value) {
    hasBoundClientSync.value = true

    watch(items, (nextItems) => {
      window.localStorage.setItem(WISHLIST_STORAGE_KEY, JSON.stringify(nextItems))
    }, { deep: true })

    window.addEventListener('storage', (event) => {
      if (event.key === WISHLIST_STORAGE_KEY) {
        syncStoredItems()
      }
    })

    window.addEventListener('focus', syncStoredItems)
  }

  if (import.meta.client) {
    onMounted(loadStoredItems)
  }

  const slugs = computed(() => items.value.map((item) => item.slug))
  const count = computed(() => items.value.length)
  const isEmpty = computed(() => items.value.length === 0)
  const isLoaded = computed(() => hasLoaded.value)
  const slugSet = computed(() => new Set(slugs.value))
  const savedProductsBySlug = computed(() =>
    new Map(
      items.value
        .filter((item): item is WishlistItem & { product: Product } =>
          Boolean(item.product),
        )
        .map((item) => [item.slug, item.product]),
    ),
  )

  const has = (productSlug: string): boolean => slugSet.value.has(productSlug)

  const add = (product: Product | string): void => {
    const productSlug = typeof product === 'string' ? product : product.slug

    if (has(productSlug)) {
      if (typeof product !== 'string') {
        updateProduct(product)
      }

      return
    }

    items.value = [
      ...items.value,
      typeof product === 'string'
        ? { slug: productSlug }
        : { slug: productSlug, product },
    ]
  }

  const remove = (productSlug: string): void => {
    items.value = items.value.filter((item) => item.slug !== productSlug)
  }

  const toggle = (product: Product | string): void => {
    const productSlug = typeof product === 'string' ? product : product.slug

    if (has(productSlug)) {
      remove(productSlug)
      return
    }

    add(product)
  }

  const clear = (): void => {
    items.value = []
  }

  const getSavedProduct = (productSlug: string): Product | undefined =>
    savedProductsBySlug.value.get(productSlug)

  function updateProduct(product: Product): void {
    if (!has(product.slug)) {
      return
    }

    items.value = items.value.map((item) =>
      item.slug === product.slug ? { slug: product.slug, product } : item,
    )
  }

  return {
    add,
    clear,
    count,
    getSavedProduct,
    has,
    isEmpty,
    isLoaded,
    items,
    remove,
    slugs,
    toggle,
    updateProduct,
  }
}
