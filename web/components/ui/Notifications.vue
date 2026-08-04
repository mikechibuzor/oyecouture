<script setup lang="ts">
const { notifications, dismiss } = useAppNotifications()
</script>

<template>
  <Teleport to="body">
    <TransitionGroup name="app-toast" tag="div" class="app-notifications">
      <article
        v-for="notification in notifications"
        :key="notification.id"
        class="app-toast"
        :data-tone="notification.tone"
        :role="notification.tone === 'error' ? 'alert' : 'status'"
      >
        <div class="app-toast-mark" aria-hidden="true" />
        <div class="app-toast-copy">
          <strong>{{ notification.title }}</strong>
          <p v-if="notification.message">{{ notification.message }}</p>
        </div>
        <button
          type="button"
          class="app-toast-close"
          aria-label="Dismiss notification"
          @click="dismiss(notification.id)"
        >
          x
        </button>
      </article>
    </TransitionGroup>
  </Teleport>
</template>
