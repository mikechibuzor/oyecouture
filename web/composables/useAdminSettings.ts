import { ADMIN_API_ROUTES } from '~/constants/admin'
import type { AdminSettings, SiteSettingsDraft } from '~/types/settings'

export const useAdminSettings = () => {
  const apiBaseUrl = useApiBaseUrl()

  // The admin endpoint, which returns the owner-only fields as well. The public
  // /site-settings read deliberately does not carry them.
  const fetchSettings = (): Promise<AdminSettings> =>
    $fetch<AdminSettings>(ADMIN_API_ROUTES.settings, {
      baseURL: apiBaseUrl,
      credentials: 'include',
    })

  const updateSettings = (draft: SiteSettingsDraft): Promise<AdminSettings> =>
    $fetch<AdminSettings>(ADMIN_API_ROUTES.settings, {
      baseURL: apiBaseUrl,
      body: draft,
      credentials: 'include',
      method: 'PATCH',
    })

  return { fetchSettings, updateSettings }
}
