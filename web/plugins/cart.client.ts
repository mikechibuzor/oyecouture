import { CART_MAXIMUM_QUANTITY, CART_STORAGE_KEY } from '~/constants/cart'
import type { CartLine } from '~/types/cart'
import { ProductSection } from '~/types/product'

// Accepts a stored line only if every field is the right shape. A bag written by
// an older build, hand-edited, or corrupted is discarded rather than crashing the
// storefront or producing an order with a bad price.
const isValidCartLine = (value: unknown): value is CartLine => {
  if (typeof value !== 'object' || value === null) {
    return false
  }

  const line = value as Record<string, unknown>
  const validSections = Object.values(ProductSection)

  return (
    typeof line.key === 'string'
    && typeof line.slug === 'string'
    && typeof line.name === 'string'
    && (line.section === undefined || validSections.includes(line.section as ProductSection))
    && typeof line.category === 'string'
    && typeof line.priceKobo === 'number'
    && Number.isInteger(line.priceKobo)
    && line.priceKobo >= 0
    && typeof line.hasSizes === 'boolean'
    && (typeof line.size === 'string' || line.size === null)
    && (typeof line.imageUrl === 'string' || line.imageUrl === null)
    && typeof line.quantity === 'number'
    && Number.isInteger(line.quantity)
    && line.quantity >= 1
    && line.quantity <= CART_MAXIMUM_QUANTITY
  )
}

// Restores the bag on load and writes it back on every change. Client-only:
// localStorage does not exist during server rendering.
export default defineNuxtPlugin(() => {
  const { lines } = useCart()

  // Restore first, so the badge and drawer render the real bag immediately.
  try {
    const storedCart = localStorage.getItem(CART_STORAGE_KEY)
    if (storedCart) {
      const parsedCart: unknown = JSON.parse(storedCart)
      const validLines = Array.isArray(parsedCart)
        ? parsedCart.filter(isValidCartLine)
        : []

      lines.value = validLines

      // Rewrite immediately when anything was rejected, so bad data does not
      // sit in storage until the customer happens to change the bag.
      if (!Array.isArray(parsedCart) || validLines.length !== parsedCart.length) {
        localStorage.setItem(CART_STORAGE_KEY, JSON.stringify(validLines))
      }
    }
  } catch {
    // Unreadable bag: start empty rather than block the storefront.
    lines.value = []
    localStorage.removeItem(CART_STORAGE_KEY)
  }

  // Persist every subsequent change.
  watch(
    lines,
    (currentLines) => {
      try {
        localStorage.setItem(CART_STORAGE_KEY, JSON.stringify(currentLines))
      } catch {
        // The bag still works for this visit in private browsing or at full quota.
      }
    },
    { deep: true },
  )
})
