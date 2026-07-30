// Versioned so a future shape change can be detected and discarded rather than
// crashing on a bag saved by an older build.
export const CART_STORAGE_KEY = 'oyecouture.cart.v1' as const

export const CART_MINIMUM_QUANTITY = 1
export const CART_MAXIMUM_QUANTITY = 10

export const CART_COPY = {
  title: 'Your Bag',
  close: 'Close bag',
  remove: 'Remove from bag',
  decrease: 'Decrease quantity',
  increase: 'Increase quantity',
  quantity: 'Quantity',
  size: 'Size',
  subtotal: 'Subtotal',
  empty: 'Your bag is empty.',
  emptyAction: 'Continue shopping',
  bespokeNote:
    'This order will be confirmed with us on WhatsApp — sizing, fabric & delivery details are finalised there.',
  readyToWearNote:
    'Delivery is calculated at checkout.',
  confirmOnWhatsApp: 'Confirm order on WhatsApp',
  checkout: 'Proceed to checkout',
  openBag: 'Open bag',
} as const
