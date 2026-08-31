import {
  CART_MAXIMUM_QUANTITY,
  CART_MINIMUM_QUANTITY,
} from '~/constants/cart'
import type { CartLine } from '~/types/cart'
import { CartCheckoutRoute } from '~/types/cart'
import type { Product, ProductSize } from '~/types/product'

// Identifies a line. The same shirt in M and in L are two lines; the same shirt
// in M twice is one line with quantity 2.
export const buildCartLineKey = (
  productSlug: string,
  size: ProductSize | null,
): string => `${productSlug}::${size ?? ''}`

// Keeps ready-to-wear quantities bounded while made-to-measure lines stay singular.
const clampQuantity = (quantity: number, hasSizes: boolean): number => {
  if (!hasSizes) {
    return CART_MINIMUM_QUANTITY
  }

  return Math.min(
    Math.max(Math.trunc(quantity), CART_MINIMUM_QUANTITY),
    CART_MAXIMUM_QUANTITY,
  )
}

// Guest bag. State lives in useState so server and client share one instance;
// persistence to localStorage is handled by plugins/cart.client.ts.
export const useCart = () => {
  const lines = useState<CartLine[]>('cart-lines', () => [])
  const isOpen = useState<boolean>('cart-open', () => false)

  const count = computed(() =>
    lines.value.reduce((total, line) => total + line.quantity, 0),
  )
  const subtotalKobo = computed(() =>
    lines.value.reduce((total, line) => total + line.priceKobo * line.quantity, 0),
  )
  const isEmpty = computed(() => lines.value.length === 0)

  // Any bespoke piece routes the whole bag through confirmation, where complete
  // measurements unlock payment and an empty form keeps the WhatsApp path.
  const checkoutRoute = computed(() =>
    lines.value.some((line) => !line.hasSizes)
      ? CartCheckoutRoute.Bespoke
      : CartCheckoutRoute.ReadyToWear,
  )

  // Adds a product, merging into an existing line when slug and size match.
  // Returns false when a sized product arrives without a size, so the caller can
  // send the customer to the product page to choose one.
  const add = (product: Product, size: ProductSize | null = null): boolean => {
    if (product.hasSizes && size === null) {
      return false
    }

    const key = buildCartLineKey(product.slug, size)
    const existingLine = lines.value.find((line) => line.key === key)

    if (existingLine) {
      if (!existingLine.hasSizes) {
        return true
      }

      setQuantity(key, existingLine.quantity + 1)
      return true
    }

    lines.value = [
      ...lines.value,
      {
        key,
        slug: product.slug,
        name: product.name,
        section: product.section,
        category: product.category,
        priceKobo: product.priceKobo,
        hasSizes: product.hasSizes,
        size,
        imageUrl: product.images[0]?.url ?? null,
        quantity: CART_MINIMUM_QUANTITY,
      },
    ]

    return true
  }

  // Sets one line's quantity, removing the line when it drops below the minimum.
  const setQuantity = (key: string, quantity: number): void => {
    if (quantity < CART_MINIMUM_QUANTITY) {
      remove(key)
      return
    }

    lines.value = lines.value.map((line) =>
      line.key === key
        ? { ...line, quantity: clampQuantity(quantity, line.hasSizes) }
        : line,
    )
  }

  const remove = (key: string): void => {
    lines.value = lines.value.filter((line) => line.key !== key)
  }

  const clear = (): void => {
    lines.value = []
  }

  const open = (): void => {
    isOpen.value = true
  }

  const close = (): void => {
    isOpen.value = false
  }

  return {
    add,
    checkoutRoute,
    clear,
    close,
    count,
    isEmpty,
    isOpen,
    lines,
    open,
    remove,
    setQuantity,
    subtotalKobo,
  }
}
