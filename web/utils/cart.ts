import type { CartLine } from '~/types/cart'
import type { MeasurementsRequest } from '~/types/order'
import { measurementFieldLabels } from '~/utils/measurements'
import { formatProductPrice } from '~/utils/product'
import { buildWhatsAppUrl } from '~/utils/whatsapp'

export interface WhatsAppOrderDetails {
  readonly customerName?: string
  readonly email?: string
  readonly phone?: string
  readonly address?: string
  readonly deliveryZone?: string
  readonly measurements?: MeasurementsRequest | null
}

const compactLines = (lines: readonly string[]): string[] => lines.filter(Boolean)

const buildDetailSection = (heading: string, lines: readonly string[]): string[] => {
  const populatedLines = compactLines(lines)

  return populatedLines.length > 0 ? ['', `${heading}:`, ...populatedLines] : []
}

const formatMeasurementValue = (value: number, unit: string): string =>
  `${Number.isInteger(value) ? value : value.toFixed(2).replace(/0+$/, '').replace(/\.$/, '')} ${unit}`

const buildMeasurementLines = (measurements?: MeasurementsRequest | null): string[] => {
  if (!measurements) {
    return []
  }

  const unit = measurements.unit ?? 'in'
  const valueLines = Object.entries(measurements.values ?? {}).map(([fieldId, value]) => {
    const label = measurementFieldLabels[fieldId] ?? fieldId.replaceAll('-', ' ')
    return `${label}: ${formatMeasurementValue(value, unit)}`
  })
  const notes = measurements.notes.trim()

  return buildDetailSection('Measurements', [
    ...valueLines,
    notes ? `Additional notes: ${notes}` : '',
  ])
}

export const buildWhatsAppOrderMessage = (
  lines: readonly CartLine[],
  subtotalKobo: number,
  orderNo = '',
  details?: WhatsAppOrderDetails,
): string => {
  const orderLines = lines.map((line) => {
    const sizeSuffix = line.size ? ` (Size ${line.size})` : ''
    const lineTotal = formatProductPrice(line.priceKobo * line.quantity)

    return `• ${line.name}${sizeSuffix} x${line.quantity} — ${lineTotal}`
  })
  const opening = orderNo
    ? `Hi Oyecouture! I'd like to confirm order ${orderNo}:`
    : "Hi Oyecouture! I'd like to order the following:"
  const customerSection = buildDetailSection('Customer details', [
    details?.customerName ? `Name: ${details.customerName}` : '',
    details?.email ? `Email: ${details.email}` : '',
    details?.phone ? `Phone: ${details.phone}` : '',
  ])
  const deliverySection = buildDetailSection('Delivery details', [
    details?.deliveryZone ? `Area: ${details.deliveryZone}` : '',
    details?.address ? `Address: ${details.address}` : '',
  ])

  return [
    opening,
    '',
    ...orderLines,
    '',
    `Subtotal: ${formatProductPrice(subtotalKobo)}`,
    ...customerSection,
    ...deliverySection,
    ...buildMeasurementLines(details?.measurements),
    '',
    'Please confirm availability, sizing and delivery details.',
  ].join('\n')
}

// Builds the wa.me deep link that opens WhatsApp with the bag already written
// out. Phase 1 uses a plain deep link — no WhatsApp Business API, no server call.
export const buildWhatsAppOrderUrl = (
  whatsAppNumber: string,
  lines: readonly CartLine[],
  subtotalKobo: number,
  orderNo = '',
  details?: WhatsAppOrderDetails,
): string => {
  // wa.me expects digits only — strip any spacing or leading plus the number
  // may have been configured with.
  return buildWhatsAppUrl(
    whatsAppNumber,
    buildWhatsAppOrderMessage(lines, subtotalKobo, orderNo, details),
  )
}
