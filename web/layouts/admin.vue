<script setup lang="ts">
import { ADMIN_COPY, ADMIN_ROUTES } from '~/constants/admin'

const { session, logOut } = useAdminSession()
const router = useRouter()

// The admin must never be indexed. This sits alongside the X-Robots-Tag header
// in nuxt.config and the public robots.txt rule.
useHead({
  meta: [{ name: 'robots', content: 'noindex, nofollow' }],
})

// Signing out always lands on the login page, even if the request failed.
const handleSignOut = async (): Promise<void> => {
  await logOut()
  await router.push(ADMIN_ROUTES.login)
}
</script>

<template>
  <div class="admin-shell">
    <header v-if="session" class="admin-bar">
      <p class="admin-bar-brand">
        {{ ADMIN_COPY.brand }}
        <span>{{ ADMIN_COPY.panel }}</span>
      </p>

      <div class="admin-bar-session">
        <span class="admin-bar-email">{{ session.email }}</span>
        <button type="button" class="admin-bar-signout" @click="handleSignOut">
          {{ ADMIN_COPY.signOut }}
        </button>
      </div>
    </header>

    <main class="admin-main">
      <slot />
    </main>
  </div>
</template>
