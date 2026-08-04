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
  const sessionRequestVersion = useState<number>('admin-session-request-version', () => 0)

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
    const requestVersion = ++sessionRequestVersion.value

    try {
      const refreshedSession = await request<AdminSession>(ADMIN_API_ROUTES.me)

      if (requestVersion === sessionRequestVersion.value) {
        session.value = refreshedSession
      }
    } catch {
      if (requestVersion === sessionRequestVersion.value) {
        session.value = null
      }
    }

    return session.value ?? null
  }

  const logIn = async (email: string, password: string): Promise<void> => {
    const loggedInSession = await request<AdminSession>(ADMIN_API_ROUTES.login, {
      method: 'POST',
      body: { email, password },
    })

    sessionRequestVersion.value += 1
    session.value = loggedInSession
  }

  // Clearing local state matters even if the request fails: the cookie may
  // already be gone, and the UI must not keep showing a signed-in shell.
  const logOut = async (): Promise<void> => {
    const requestVersion = ++sessionRequestVersion.value

    try {
      await request(ADMIN_API_ROUTES.logout, { method: 'POST' })
    } finally {
      if (requestVersion === sessionRequestVersion.value) {
        session.value = null
      }
    }
  }

  return { session, refresh, logIn, logOut }
}
