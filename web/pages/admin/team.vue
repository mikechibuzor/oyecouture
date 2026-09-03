<script setup lang="ts">
import { ADMIN_COPY, ADMIN_ROUTES } from '~/constants/admin'
import { AdminRole } from '~/types/team'
import type { TeamMember } from '~/types/team'
import { formatOrderDate } from '~/utils/admin'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const router = useRouter()
const { session } = useAdminSession()
const { fetchTeam, createMember, updateMember, deleteMember, changePassword } = useAdminTeam()
const { notifyError, notifySuccess } = useAppNotifications()

const members = ref<readonly TeamMember[]>([])
const isLoading = ref(true)
const hasFailed = ref(false)
const savingMemberId = ref('')

const inviteForm = reactive({ email: '', password: '', role: AdminRole.Staff })
const isCreating = ref(false)
const inviteError = ref('')

const passwordForm = reactive({ current: '', next: '', confirm: '' })
const isChangingPassword = ref(false)
const passwordError = ref('')

useHead({ title: `${ADMIN_COPY.teamHeading} · ${ADMIN_COPY.brand}` })

const roleOptions = [
  { value: AdminRole.Owner, label: ADMIN_COPY.teamRoleOwner },
  { value: AdminRole.Staff, label: ADMIN_COPY.teamRoleStaff },
]

const roleLabel = (role: AdminRole): string =>
  role === AdminRole.Owner ? ADMIN_COPY.teamRoleOwner : ADMIN_COPY.teamRoleStaff

// The signed-in admin is matched by email: the session carries no id, and the
// API rejects self-targeted changes regardless, so this only decides whether
// the controls are offered.
const isSelf = (member: TeamMember): boolean =>
  member.email.toLowerCase() === (session.value?.email ?? '').toLowerCase()

const loadTeam = async (): Promise<void> => {
  isLoading.value = true
  hasFailed.value = false

  try {
    members.value = await fetchTeam()
  } catch (loadError) {
    if ((loadError as { status?: number } | null)?.status === 401) {
      session.value = null
      await router.push(ADMIN_ROUTES.login)
      return
    }

    hasFailed.value = true
  } finally {
    isLoading.value = false
  }
}

// The API's message is shown as-is: it explains exactly which rule was hit,
// such as the last owner being protected, which a generic failure would lose.
const messageFrom = (thrown: unknown, fallback: string): string => {
  const message = (thrown as { data?: { error?: string } } | null)?.data?.error

  return message || fallback
}

const submitInvite = async (): Promise<void> => {
  isCreating.value = true
  inviteError.value = ''

  try {
    await createMember({ ...inviteForm })
    inviteForm.email = ''
    inviteForm.password = ''
    inviteForm.role = AdminRole.Staff
    notifySuccess(ADMIN_COPY.teamCreatedNotice)
    await loadTeam()
  } catch (createError) {
    inviteError.value = messageFrom(createError, ADMIN_COPY.teamLoadFailed)
    notifyError(inviteError.value)
  } finally {
    isCreating.value = false
  }
}

const saveMember = async (member: TeamMember, role: AdminRole, active: boolean): Promise<void> => {
  savingMemberId.value = member.id

  try {
    await updateMember(member.id, { role, active })
    notifySuccess(ADMIN_COPY.teamSavedNotice)
    await loadTeam()
  } catch (updateError) {
    notifyError(messageFrom(updateError, ADMIN_COPY.teamLoadFailed))
    await loadTeam()
  } finally {
    savingMemberId.value = ''
  }
}

const removeMember = async (member: TeamMember): Promise<void> => {
  // eslint-disable-next-line no-alert
  if (!window.confirm(ADMIN_COPY.teamDeleteConfirm)) {
    return
  }

  savingMemberId.value = member.id

  try {
    await deleteMember(member.id)
    notifySuccess(ADMIN_COPY.teamDeletedNotice)
    await loadTeam()
  } catch (deleteError) {
    notifyError(messageFrom(deleteError, ADMIN_COPY.teamLoadFailed))
  } finally {
    savingMemberId.value = ''
  }
}

const submitPasswordChange = async (): Promise<void> => {
  passwordError.value = ''

  // Checked here as well as by the API so the mistake is caught before a
  // round trip, and before the current password leaves the browser again.
  if (passwordForm.next !== passwordForm.confirm) {
    passwordError.value = ADMIN_COPY.passwordMismatch
    return
  }

  isChangingPassword.value = true

  try {
    await changePassword({
      currentPassword: passwordForm.current,
      newPassword: passwordForm.next,
    })
    passwordForm.current = ''
    passwordForm.next = ''
    passwordForm.confirm = ''
    notifySuccess(ADMIN_COPY.passwordChangedNotice)
  } catch (changeError) {
    passwordError.value = messageFrom(changeError, ADMIN_COPY.teamLoadFailed)
    notifyError(passwordError.value)
  } finally {
    isChangingPassword.value = false
  }
}

onMounted(loadTeam)
</script>

<template>
  <div class="admin-page">
    <header class="admin-page-head">
      <div>
        <h1>{{ ADMIN_COPY.teamHeading }}</h1>
        <p class="admin-page-subheading">{{ ADMIN_COPY.teamSubheading }}</p>
      </div>
    </header>

    <div class="admin-settings-form">
      <section class="admin-form-section" aria-labelledby="admin-password-heading">
        <header class="admin-form-section-head">
          <h2 id="admin-password-heading">{{ ADMIN_COPY.passwordHeading }}</h2>
          <p>{{ ADMIN_COPY.passwordIntro }}</p>
        </header>

        <form class="admin-team-fields" @submit.prevent="submitPasswordChange">
          <label class="admin-field">
            <span>{{ ADMIN_COPY.currentPasswordLabel }}</span>
            <input
              v-model="passwordForm.current"
              type="password"
              autocomplete="current-password"
              required
            >
          </label>
          <label class="admin-field">
            <span>{{ ADMIN_COPY.newPasswordLabel }}</span>
            <input
              v-model="passwordForm.next"
              type="password"
              autocomplete="new-password"
              minlength="12"
              required
            >
          </label>
          <label class="admin-field">
            <span>{{ ADMIN_COPY.confirmPasswordLabel }}</span>
            <input
              v-model="passwordForm.confirm"
              type="password"
              autocomplete="new-password"
              minlength="12"
              required
            >
          </label>

          <p v-if="passwordError" class="admin-detail-error admin-team-span" role="alert">
            {{ passwordError }}
          </p>

          <div class="admin-form-actions admin-team-span">
            <button
              type="submit"
              class="admin-inline-action is-primary"
              :disabled="isChangingPassword"
            >
              {{ isChangingPassword ? ADMIN_COPY.changingPassword : ADMIN_COPY.changePassword }}
            </button>
          </div>
        </form>
      </section>

      <section class="admin-form-section" aria-labelledby="admin-invite-heading">
        <header class="admin-form-section-head">
          <h2 id="admin-invite-heading">{{ ADMIN_COPY.teamInviteHeading }}</h2>
          <p>{{ ADMIN_COPY.teamInviteIntro }}</p>
        </header>

        <form class="admin-team-fields" @submit.prevent="submitInvite">
          <label class="admin-field">
            <span>{{ ADMIN_COPY.teamEmailLabel }}</span>
            <input v-model.trim="inviteForm.email" type="email" autocomplete="off" required>
          </label>

          <div class="admin-field">
            <span>{{ ADMIN_COPY.teamRoleLabel }}</span>
            <UiSelect
              v-model="inviteForm.role"
              :options="roleOptions"
              :label="ADMIN_COPY.teamRoleLabel"
            />
            <small>
              {{
                inviteForm.role === AdminRole.Owner
                  ? ADMIN_COPY.teamRoleOwnerHint
                  : ADMIN_COPY.teamRoleStaffHint
              }}
            </small>
          </div>

          <label class="admin-field admin-team-span">
            <span>{{ ADMIN_COPY.teamTemporaryPasswordLabel }}</span>
            <input v-model="inviteForm.password" type="text" minlength="12" required>
            <small>{{ ADMIN_COPY.teamTemporaryPasswordHint }}</small>
          </label>

          <p v-if="inviteError" class="admin-detail-error admin-team-span" role="alert">
            {{ inviteError }}
          </p>

          <div class="admin-form-actions admin-team-span">
            <button type="submit" class="admin-inline-action is-primary" :disabled="isCreating">
              {{ isCreating ? ADMIN_COPY.teamCreating : ADMIN_COPY.teamCreate }}
            </button>
          </div>
        </form>
      </section>

      <section class="admin-form-section" aria-labelledby="admin-team-list-heading">
        <header class="admin-form-section-head">
          <h2 id="admin-team-list-heading">{{ ADMIN_COPY.teamListHeading }}</h2>
          <p>{{ ADMIN_COPY.teamListIntro }}</p>
        </header>

        <p v-if="isLoading" class="admin-orders-message">{{ ADMIN_COPY.teamLoading }}</p>
        <p v-else-if="hasFailed" class="admin-detail-error" role="alert">
          {{ ADMIN_COPY.teamLoadFailed }}
        </p>

        <ul v-else class="admin-team-list">
          <li v-for="member in members" :key="member.id" class="admin-team-member">
            <div class="admin-team-member-identity">
              <p class="admin-team-member-email">
                {{ member.email }}
                <span v-if="isSelf(member)" class="admin-team-you">{{ ADMIN_COPY.teamYou }}</span>
              </p>
              <p class="admin-team-member-meta">
                {{ formatOrderDate(member.createdAt) }}
              </p>
            </div>

            <span
              class="admin-team-status"
              :class="member.active ? 'is-active' : 'is-revoked'"
            >
              {{ member.active ? ADMIN_COPY.teamActiveStatus : ADMIN_COPY.teamRevokedStatus }}
            </span>

            <div class="admin-team-member-role">
              <UiSelect
                v-if="!isSelf(member)"
                :model-value="member.role"
                :options="roleOptions"
                :label="ADMIN_COPY.teamRoleLabel"
                variant="bare"
                :disabled="savingMemberId === member.id"
                @update:model-value="(role) => saveMember(member, role, member.active)"
              />
              <span v-else class="admin-team-member-meta">{{ roleLabel(member.role) }}</span>
            </div>

            <!-- Self-targeted controls are omitted rather than disabled: the API
                 refuses them, and an inert button invites the attempt. -->
            <div v-if="!isSelf(member)" class="admin-team-member-actions">
              <button
                type="button"
                class="admin-inline-action"
                :disabled="savingMemberId === member.id"
                @click="saveMember(member, member.role, !member.active)"
              >
                {{ member.active ? ADMIN_COPY.teamRevoke : ADMIN_COPY.teamRestore }}
              </button>
              <button
                type="button"
                class="admin-inline-action is-danger"
                :disabled="savingMemberId === member.id"
                @click="removeMember(member)"
              >
                {{ ADMIN_COPY.teamDelete }}
              </button>
            </div>
            <p v-else class="admin-team-member-meta">{{ ADMIN_COPY.teamSelfHint }}</p>
          </li>
        </ul>
      </section>
    </div>
  </div>
</template>
