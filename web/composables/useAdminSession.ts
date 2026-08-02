import { ADMIN_API_ROUTES } from '~/constants/admin'
import type { AdminSession } from '~/types/admin'

// The signed-in admin, shared across the admin pages. `null` means signed out;
// `undefined` means we have not asked the server yet.
export const useAdminSession = () => {
  const apiBaseUrl = useApiBaseUrl()
  const session = useState<AdminSession | null | undefined>(
    'admin-session',
    () => undefined,
  )

  // Every admin call carries the session cookie. In development the storefront
  // and API sit on different ports, so the request must be credentialed.
  const request = <ResponseBody>(
    path: string,
    options: Parameters<typeof $fetch<ResponseBody>>[1] = {},
  ): Promise<ResponseBody> =>
    $fetch<ResponseBody>(path, {
      baseURL: apiBaseUrl,
      credentials: 'include',
      ...options,
    })

  // Asks the server who we are. A 401 is an answer, not a failure.
  const refresh = async (): Promise<AdminSession | null> => {
    try {
      session.value = await request<AdminSession>(ADMIN_API_ROUTES.me)
    } catch {
      session.value = null
    }

    return session.value
  }

  const logIn = async (email: string, password: string): Promise<void> => {
    session.value = await request<AdminSession>(ADMIN_API_ROUTES.login, {
      method: 'POST',
      body: { email, password },
    })
  }

  // Clearing local state matters even if the request fails: the cookie may
  // already be gone, and the UI must not keep showing a signed-in shell.
  const logOut = async (): Promise<void> => {
    try {
      await request(ADMIN_API_ROUTES.logout, { method: 'POST' })
    } finally {
      session.value = null
    }
  }

  return { session, refresh, logIn, logOut }
}
