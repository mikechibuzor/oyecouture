import type { CartLine } from '~/types/cart'
import { formatProductPrice } from '~/utils/product'

const WHATSAPP_BASE_URL = 'https://wa.me/'

// Builds the wa.me deep link that opens WhatsApp with the bag already written
// out. Phase 1 uses a plain deep link — no WhatsApp Business API, no server call.
export const buildWhatsAppOrderUrl = (
  whatsAppNumber: string,
  lines: readonly CartLine[],
  subtotalKobo: number,
  orderNo = '',
): string => {
  // Wording follows the approved design: "Hi Oyecouture! …" with x-quantity
  // line items and a closing request to confirm availability and sizing.
  const orderLines = lines.map((line) => {
    const sizeSuffix = line.size ? ` (Size ${line.size})` : ''
    const lineTotal = formatProductPrice(line.priceKobo * line.quantity)

    return `• ${line.name}${sizeSuffix} x${line.quantity} — ${lineTotal}`
  })

  // Once the order is saved it carries a reference, so the conversation and the
  // order record stay tied together.
  const opening = orderNo
    ? `Hi Oyecouture! I'd like to confirm order ${orderNo}:`
    : "Hi Oyecouture! I'd like to order the following:"

  const message = [
    opening,
    '',
    ...orderLines,
    '',
    `Subtotal: ${formatProductPrice(subtotalKobo)}`,
    '',
    'Please confirm availability, sizing and delivery details.',
  ].join('\n')

  // wa.me expects digits only — strip any spacing or leading plus the number
  // may have been configured with.
  const normalizedNumber = whatsAppNumber.replace(/\D/g, '')

  return `${WHATSAPP_BASE_URL}${normalizedNumber}?text=${encodeURIComponent(message)}`
}
