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
  <section class="admin-page">
    <header class="admin-page-header">
      <h1>{{ ADMIN_COPY.teamHeading }}</h1>
      <p>{{ ADMIN_COPY.teamSubheading }}</p>
    </header>

    <div class="admin-panel">
      <h2 class="admin-panel-heading">{{ ADMIN_COPY.passwordHeading }}</h2>
      <p class="admin-panel-intro">{{ ADMIN_COPY.passwordIntro }}</p>

      <form class="admin-form-grid" @submit.prevent="submitPasswordChange">
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

        <p v-if="passwordError" class="admin-detail-error" role="alert">{{ passwordError }}</p>

        <div class="admin-form-actions">
          <button type="submit" class="admin-inline-action is-primary" :disabled="isChangingPassword">
            {{ isChangingPassword ? ADMIN_COPY.changingPassword : ADMIN_COPY.changePassword }}
          </button>
        </div>
      </form>
    </div>

    <div class="admin-panel">
      <h2 class="admin-panel-heading">{{ ADMIN_COPY.teamInviteHeading }}</h2>

      <form class="admin-form-grid" @submit.prevent="submitInvite">
        <label class="admin-field">
          <span>{{ ADMIN_COPY.teamEmailLabel }}</span>
          <input v-model.trim="inviteForm.email" type="email" required>
        </label>
        <label class="admin-field">
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
        </label>
        <label class="admin-field admin-field-wide">
          <span>{{ ADMIN_COPY.teamTemporaryPasswordLabel }}</span>
          <input v-model="inviteForm.password" type="text" minlength="12" required>
          <small>{{ ADMIN_COPY.teamTemporaryPasswordHint }}</small>
        </label>

        <p v-if="inviteError" class="admin-detail-error" role="alert">{{ inviteError }}</p>

        <div class="admin-form-actions">
          <button type="submit" class="admin-inline-action is-primary" :disabled="isCreating">
            {{ isCreating ? ADMIN_COPY.teamCreating : ADMIN_COPY.teamCreate }}
          </button>
        </div>
      </form>
    </div>

    <p v-if="isLoading" class="admin-empty">{{ ADMIN_COPY.teamLoading }}</p>
    <p v-else-if="hasFailed" class="admin-detail-error" role="alert">
      {{ ADMIN_COPY.teamLoadFailed }}
    </p>

    <div v-else class="admin-table-wrapper">
      <table class="admin-table">
        <thead>
          <tr>
            <th>{{ ADMIN_COPY.teamEmailLabel }}</th>
            <th>{{ ADMIN_COPY.teamRoleLabel }}</th>
            <th>{{ ADMIN_COPY.teamStatusLabel }}</th>
            <th>{{ ADMIN_COPY.teamAddedLabel }}</th>
            <th aria-hidden="true" />
          </tr>
        </thead>
        <tbody>
          <tr v-for="member in members" :key="member.id">
            <td>
              {{ member.email }}
              <small v-if="isSelf(member)" class="admin-table-note">{{ ADMIN_COPY.teamYou }}</small>
            </td>
            <td>
              <UiSelect
                v-if="!isSelf(member)"
                :model-value="member.role"
                :options="roleOptions"
                :label="ADMIN_COPY.teamRoleLabel"
                variant="bare"
                :disabled="savingMemberId === member.id"
                @update:model-value="(role) => saveMember(member, role, member.active)"
              />
              <span v-else>{{ roleLabel(member.role) }}</span>
            </td>
            <td>
              <span class="admin-status-pill" :class="member.active ? 'is-positive' : 'is-muted'">
                {{ member.active ? ADMIN_COPY.teamActiveStatus : ADMIN_COPY.teamRevokedStatus }}
              </span>
            </td>
            <td>{{ formatOrderDate(member.createdAt) }}</td>
            <td class="admin-table-actions">
              <!-- Self-targeted controls are omitted rather than disabled: the
                   API refuses them, and an inert button invites the attempt. -->
              <template v-if="!isSelf(member)">
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
              </template>
              <small v-else class="admin-table-note">{{ ADMIN_COPY.teamSelfHint }}</small>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>
