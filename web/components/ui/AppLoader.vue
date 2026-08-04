<script setup lang="ts">
const { isAppLoading, loadingLabel } = useAppLoading()
</script>

<template>
  <Transition name="app-loader">
    <div v-if="isAppLoading" class="app-loader" role="status" aria-live="polite">
      <span class="app-loader-bar" aria-hidden="true" />
      <span class="app-loader-label">
        {{ loadingLabel || 'Working...' }}
      </span>
    </div>
  </Transition>
</template>
