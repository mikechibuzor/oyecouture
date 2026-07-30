import type { CartLine } from '~/types/cart'
import { formatProductPrice } from '~/utils/product'

const WHATSAPP_BASE_URL = 'https://wa.me/'

// Builds the wa.me deep link that opens WhatsApp with the bag already written
// out. Phase 1 uses a plain deep link — no WhatsApp Business API, no server call.
export const buildWhatsAppOrderUrl = (
  whatsAppNumber: string,
  lines: readonly CartLine[],
  subtotalKobo: number,
): string => {
  const orderLines = lines.map((line) => {
    const sizeSuffix = line.size ? ` (Size ${line.size})` : ''
    const lineTotal = formatProductPrice(line.priceKobo * line.quantity)

    return `• ${line.name}${sizeSuffix} × ${line.quantity} — ${lineTotal}`
  })

  const message = [
    'Hello Oyecouture, I would like to confirm this order:',
    '',
    ...orderLines,
    '',
    `Subtotal: ${formatProductPrice(subtotalKobo)}`,
  ].join('\n')

  // wa.me expects digits only — strip any spacing or leading plus the number
  // may have been configured with.
  const normalizedNumber = whatsAppNumber.replace(/\D/g, '')

  return `${WHATSAPP_BASE_URL}${normalizedNumber}?text=${encodeURIComponent(message)}`
}
