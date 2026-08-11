import { ADMIN_API_ROUTES } from '~/constants/admin'
import type { HomeContent, HomeContentDraft } from '~/types/settings'

export const useAdminHomeContent = () => {
  const apiBaseUrl = useApiBaseUrl()

  const fetchHomeContent = (): Promise<HomeContent> =>
    $fetch<HomeContent>(ADMIN_API_ROUTES.homeContent, {
      baseURL: apiBaseUrl,
      credentials: 'include',
    })

  const updateHomeContent = (draft: HomeContentDraft): Promise<HomeContent> =>
    $fetch<HomeContent>(ADMIN_API_ROUTES.homeContent, {
      baseURL: apiBaseUrl,
      body: draft,
      credentials: 'include',
      method: 'PATCH',
    })

  return { fetchHomeContent, updateHomeContent }
}
