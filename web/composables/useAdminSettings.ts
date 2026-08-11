import { ADMIN_API_ROUTES } from '~/constants/admin'
import type { SiteSettings, SiteSettingsDraft } from '~/types/settings'

export const useAdminSettings = () => {
  const apiBaseUrl = useApiBaseUrl()

  const fetchSettings = (): Promise<SiteSettings> =>
    $fetch<SiteSettings>(ADMIN_API_ROUTES.settings, {
      baseURL: apiBaseUrl,
      credentials: 'include',
    })

  const updateSettings = (draft: SiteSettingsDraft): Promise<SiteSettings> =>
    $fetch<SiteSettings>(ADMIN_API_ROUTES.settings, {
      baseURL: apiBaseUrl,
      body: draft,
      credentials: 'include',
      method: 'PATCH',
    })

  return { fetchSettings, updateSettings }
}
