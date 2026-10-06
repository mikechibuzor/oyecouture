import { ADMIN_API_ROUTES } from '~/constants/admin'
import type { SiteSettings } from '~/types/settings'

/**
 * Decides whether this visitor sees the holding page instead of the shop.
 *
 * Two things settle it: the owner's "Hold the shop" switch, and whether the request
 * carries an admin session. A signed-in admin always sees the real shop, so the owner can
 * keep reviewing it at its own address while visitors see the holding page — no second
 * link and no doubt about which version is being looked at.
 *
 * The API enforces the same rule on its own routes. This only decides what gets drawn; it
 * is not what keeps the shop closed.
 *
 * Everything Nuxt-aware is read before the first await, because a composable called past
 * one throws "[nuxt] instance unavailable" and would take every page down with it.
 */
export const useSiteHold = async () => {
  const apiBaseUrl = useApiBaseUrl()
  const forwardedCookies = import.meta.server ? useRequestHeaders(['cookie']) : undefined
  const isHeld = useState('site-held', () => false)
  const { data: settings } = await useSiteSettings()

  // Only ask who the visitor is when the answer changes what they see. The admin session
  // cookie is HttpOnly, so on the server it has to be forwarded by hand.
  if (settings.value?.held !== true) {
    isHeld.value = false

    return { settings, isHeld }
  }

  isHeld.value = await $fetch(ADMIN_API_ROUTES.me, {
    baseURL: apiBaseUrl,
    headers: forwardedCookies,
    credentials: 'include',
    retry: 0,
  })
    .then(() => false)
    .catch(() => true)

  return { settings, isHeld }
}
