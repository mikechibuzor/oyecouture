<script setup lang="ts">
const { isAppLoading, loadingLabel } = useAppLoading()
const isLoaderVisible = ref(false)
const APP_LOADER_DELAY_MS = 160
let visibilityTimer: number | undefined

watch(isAppLoading, (isLoading) => {
  if (!import.meta.client) {
    return
  }

  window.clearTimeout(visibilityTimer)

  if (!isLoading) {
    isLoaderVisible.value = false
    return
  }

  visibilityTimer = window.setTimeout(() => {
    isLoaderVisible.value = true
  }, APP_LOADER_DELAY_MS)
}, { immediate: true })

onBeforeUnmount(() => {
  if (import.meta.client) {
    window.clearTimeout(visibilityTimer)
  }
})
</script>

<template>
  <Transition name="app-loader">
    <div v-if="isLoaderVisible" class="app-loader" role="status" aria-live="polite">
      <span class="app-loader-bar" aria-hidden="true" />
      <span class="visually-hidden">
        {{ loadingLabel || 'Working...' }}
      </span>
    </div>
  </Transition>
</template>
