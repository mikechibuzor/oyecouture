import { API_ROUTES } from '~/constants/api'
import { STOREFRONT_COPY } from '~/constants/storefront'
import type { SiteSettings } from '~/types/settings'

export const defaultSiteSettings = (): SiteSettings => {
  const runtimeConfig = useRuntimeConfig()

  return {
    announcement: STOREFRONT_COPY.announcement,
    whatsappNumber: String(runtimeConfig.public.whatsappNumber ?? ''),
  }
}

export const useSiteSettings = () => {
  const fallbackSettings = defaultSiteSettings()

  return useFetch<SiteSettings>(API_ROUTES.siteSettings, {
    baseURL: useApiBaseUrl(),
    default: () => fallbackSettings,
    key: 'site-settings',
    retry: 0,
  })
}
