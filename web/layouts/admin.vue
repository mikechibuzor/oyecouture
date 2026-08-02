<script setup lang="ts">
import { ADMIN_ROUTES } from '~/constants/admin'

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
  <div class="admin-shell" :class="{ 'is-signed-in': Boolean(session) }">
    <AdminSideNav v-if="session" :session="session" @sign-out="handleSignOut" />

    <main class="admin-main">
      <slot />
    </main>
  </div>
</template>
