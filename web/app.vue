<script setup lang="ts">
import { ADMIN_ROUTES } from '~/constants/admin'

const router = useRouter()
const route = useRoute()
const { startLoading } = useAppLoading()

// The admin is how the hold is lifted, so it is never held itself.
const isAdminRoute = computed(
  () => route.path === ADMIN_ROUTES.base || route.path.startsWith(`${ADMIN_ROUTES.base}/`),
)

const { settings, isHeld } = await useSiteHold()

let finishRouteLoading: (() => void) | undefined
let removeBeforeEachGuard: (() => void) | undefined
let removeAfterEachGuard: (() => void) | undefined

onMounted(() => {
  removeBeforeEachGuard = router.beforeEach(() => {
    finishRouteLoading?.()
    finishRouteLoading = startLoading('Loading page')
  })

  removeAfterEachGuard = router.afterEach(() => {
    finishRouteLoading?.()
    finishRouteLoading = undefined
  })
})

onBeforeUnmount(() => {
  finishRouteLoading?.()
  removeBeforeEachGuard?.()
  removeAfterEachGuard?.()
})
</script>

<template>
  <StorefrontSiteHold v-if="isHeld && !isAdminRoute" :settings="settings" />

  <template v-else>
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <UiAppLoader />
    <UiNotifications />
  </template>
</template>
