<script setup lang="ts">
import { ADMIN_COPY, ADMIN_ROUTES } from '~/constants/admin'

definePageMeta({ layout: 'admin' })

useHead({ title: `${ADMIN_COPY.panel} · ${ADMIN_COPY.brand}` })

const { session, refresh, logIn } = useAdminSession()
const route = useRoute()
const router = useRouter()

const email = ref('')
const password = ref('')
const errorMessage = ref('')
const isSubmitting = ref(false)

const isSafeAdminDestination = (path: string): boolean =>
  path.startsWith('/admin') &&
  !path.startsWith('//') &&
  path !== ADMIN_ROUTES.login &&
  !path.startsWith(`${ADMIN_ROUTES.login}/`) &&
  !path.startsWith(`${ADMIN_ROUTES.login}?`) &&
  !path.startsWith(`${ADMIN_ROUTES.login}#`)

// Where to land after signing in. Only in-app paths are honoured, so a crafted
// ?next= cannot bounce the owner to another site or back to the sign-in page.
const destination = computed(() => {
  const next = route.query.next
  const path = typeof next === 'string' ? next : ''

  return isSafeAdminDestination(path)
    ? path
    : ADMIN_ROUTES.orders
})

// An owner who still has a valid session should not be asked to sign in again.
onMounted(async () => {
  if (session.value === undefined) {
    await refresh()
  }
  if (session.value) {
    await router.replace(destination.value)
  }
})

// Maps the API's response to copy the owner can act on.
const messageForError = (submitError: unknown): string => {
  const status = (submitError as { status?: number } | null)?.status

  if (status === 429) {
    return ADMIN_COPY.loginThrottled
  }

  return status === 401 ? ADMIN_COPY.loginFailed : ADMIN_COPY.loginUnreachable
}

const handleSubmit = async (): Promise<void> => {
  if (isSubmitting.value) {
    return
  }

  isSubmitting.value = true
  errorMessage.value = ''

  try {
    await logIn(email.value, password.value)
    await router.replace(destination.value)
  } catch (submitError) {
    errorMessage.value = messageForError(submitError)
    password.value = ''
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="admin-login">
    <form class="admin-login-card" novalidate @submit.prevent="handleSubmit">
      <h1>{{ ADMIN_COPY.loginHeading }}</h1>
      <p class="admin-login-intro">{{ ADMIN_COPY.loginIntro }}</p>

      <label class="admin-field">
        <span>{{ ADMIN_COPY.emailLabel }}</span>
        <input
          v-model.trim="email"
          type="email"
          name="email"
          autocomplete="username"
          required
        >
      </label>

      <label class="admin-field">
        <span>{{ ADMIN_COPY.passwordLabel }}</span>
        <input
          v-model="password"
          type="password"
          name="password"
          autocomplete="current-password"
          required
        >
      </label>

      <p v-if="errorMessage" class="admin-login-error" role="alert">
        {{ errorMessage }}
      </p>

      <button type="submit" class="admin-login-submit" :disabled="isSubmitting">
        {{ isSubmitting ? ADMIN_COPY.signingIn : ADMIN_COPY.signIn }}
      </button>
    </form>
  </div>
</template>
