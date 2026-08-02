import { ADMIN_ROUTES } from '~/constants/admin'

// Gates every admin page. The API is the real gate — this only saves the owner
// from staring at an empty screen before the first request comes back.
export default defineNuxtRouteMiddleware(async (to) => {
  const { session, refresh } = useAdminSession()

  if (session.value === undefined) {
    await refresh()
  }

  if (!session.value) {
    return navigateTo({
      path: ADMIN_ROUTES.login,
      query: { next: to.fullPath },
    })
  }
})
