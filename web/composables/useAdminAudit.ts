import { ADMIN_API_ROUTES } from '~/constants/admin'
import type { AuditPage } from '~/types/team'

export const useAdminAudit = () => {
  const apiBaseUrl = useApiBaseUrl()

  const fetchAudit = (limit: number, offset: number): Promise<AuditPage> =>
    $fetch<AuditPage>(ADMIN_API_ROUTES.audit, {
      baseURL: apiBaseUrl,
      credentials: 'include',
      query: { limit, offset },
    })

  return { fetchAudit }
}
