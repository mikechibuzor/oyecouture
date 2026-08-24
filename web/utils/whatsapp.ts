const WHATSAPP_BASE_URL = 'https://wa.me/'

export const buildWhatsAppUrl = (whatsAppNumber: string, message = ''): string => {
  const normalizedNumber = whatsAppNumber.replace(/\D/g, '')
  if (!normalizedNumber) {
    return ''
  }

  const encodedMessage = message ? `?text=${encodeURIComponent(message)}` : ''
  return `${WHATSAPP_BASE_URL}${normalizedNumber}${encodedMessage}`
}
