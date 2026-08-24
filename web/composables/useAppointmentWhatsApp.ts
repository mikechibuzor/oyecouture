import { StorefrontRoute } from '~/types/storefront'
import { buildWhatsAppUrl } from '~/utils/whatsapp'

const APPOINTMENT_MESSAGE
  = "Hi Oyecouture, I'd like to book a private fitting appointment."

export const useAppointmentWhatsApp = async () => {
  const { data: siteSettings } = await useSiteSettings()

  const appointmentWhatsAppUrl = computed(() =>
    buildWhatsAppUrl(
      String(siteSettings.value?.whatsappNumber ?? ''),
      APPOINTMENT_MESSAGE,
    ) || StorefrontRoute.Contact,
  )

  const opensWhatsApp = computed(() =>
    appointmentWhatsAppUrl.value.startsWith('https://wa.me/'),
  )

  return {
    appointmentWhatsAppUrl,
    opensWhatsApp,
  }
}
