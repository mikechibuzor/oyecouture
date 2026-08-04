<script setup lang="ts">
const router = useRouter()
const { startLoading } = useAppLoading()

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
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
  <UiAppLoader />
  <UiNotifications />
</template>
