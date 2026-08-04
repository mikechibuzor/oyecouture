const WISHLIST_STORAGE_KEY = 'oyecouture-wishlist'

const readStoredWishlistSlugs = (): string[] => {
  if (!import.meta.client) {
    return []
  }

  try {
    const storedValue = window.localStorage.getItem(WISHLIST_STORAGE_KEY)
    const parsedValue: unknown = storedValue ? JSON.parse(storedValue) : []

    return Array.isArray(parsedValue)
      ? parsedValue.filter((value): value is string => typeof value === 'string')
      : []
  } catch {
    return []
  }
}

export const useWishlist = () => {
  const slugs = useState<string[]>('wishlist-slugs', () => [])
  const hasLoaded = useState<boolean>('wishlist-loaded', () => false)

  if (import.meta.client) {
    onMounted(() => {
      if (hasLoaded.value) {
        return
      }

      slugs.value = readStoredWishlistSlugs()
      hasLoaded.value = true
    })

    watch(
      slugs,
      (nextSlugs) => {
        window.localStorage.setItem(WISHLIST_STORAGE_KEY, JSON.stringify(nextSlugs))
      },
      { deep: true },
    )
  }

  const count = computed(() => slugs.value.length)
  const isEmpty = computed(() => slugs.value.length === 0)
  const isLoaded = computed(() => hasLoaded.value)
  const slugSet = computed(() => new Set(slugs.value))

  const has = (productSlug: string): boolean => slugSet.value.has(productSlug)

  const add = (productSlug: string): void => {
    if (has(productSlug)) {
      return
    }

    slugs.value = [...slugs.value, productSlug]
  }

  const remove = (productSlug: string): void => {
    slugs.value = slugs.value.filter((slug) => slug !== productSlug)
  }

  const toggle = (productSlug: string): void => {
    if (has(productSlug)) {
      remove(productSlug)
      return
    }

    add(productSlug)
  }

  const clear = (): void => {
    slugs.value = []
  }

  return {
    add,
    clear,
    count,
    has,
    isEmpty,
    isLoaded,
    remove,
    slugs,
    toggle,
  }
}
