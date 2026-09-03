<script setup lang="ts">
import { ADMIN_COPY, ADMIN_PILL_TONE, ADMIN_ROUTES } from '~/constants/admin'
import { HTTP_STATUS } from '~/types/http'
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

// One drawer, two purposes, matching how products edit: the list stays the
// page, and anything that needs a form opens over it.
type TeamDrawer = 'invite' | 'password'

const openDrawer = ref<TeamDrawer | null>(null)
const drawerElement = ref<HTMLElement | null>(null)
const drawerCloseButton = ref<HTMLButtonElement | null>(null)
const isDrawerOpen = computed(() => openDrawer.value !== null)

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

const drawerTitle = computed(() =>
  openDrawer.value === 'invite' ? ADMIN_COPY.teamInviteHeading : ADMIN_COPY.passwordHeading,
)

// The signed-in admin is matched by email: the session carries no id, and the
// API rejects self-targeted changes regardless, so this only decides whether
// the controls are offered.
const isSelf = (member: TeamMember): boolean =>
  member.email.toLowerCase() === (session.value?.email ?? '').toLowerCase()

const closeDrawer = (): void => {
  openDrawer.value = null
  inviteError.value = ''
  passwordError.value = ''
}

const openInvite = (): void => {
  inviteForm.email = ''
  inviteForm.password = ''
  inviteForm.role = AdminRole.Staff
  inviteError.value = ''
  openDrawer.value = 'invite'
}

const openPasswordChange = (): void => {
  passwordForm.current = ''
  passwordForm.next = ''
  passwordForm.confirm = ''
  passwordError.value = ''
  openDrawer.value = 'password'
}

useModalFocus({
  isOpen: isDrawerOpen,
  container: drawerElement,
  initialFocus: drawerCloseButton,
  onEscape: closeDrawer,
  lockScroll: false,
})

const loadTeam = async (): Promise<void> => {
  isLoading.value = true
  hasFailed.value = false

  try {
    members.value = await fetchTeam()
  } catch (loadError) {
    if ((loadError as { status?: number } | null)?.status === HTTP_STATUS.unauthorized) {
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
    notifySuccess(ADMIN_COPY.teamCreatedNotice)
    closeDrawer()
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
    notifySuccess(ADMIN_COPY.passwordChangedNotice)
    closeDrawer()
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

      <div class="admin-page-actions">
        <button type="button" class="admin-inline-action" @click="openPasswordChange">
          {{ ADMIN_COPY.changePassword }}
        </button>
        <button type="button" class="admin-inline-action" @click="openInvite">
          {{ ADMIN_COPY.teamCreate }}
        </button>
      </div>
    </header>

    <p v-if="isLoading" class="admin-orders-message">{{ ADMIN_COPY.teamLoading }}</p>
    <p v-else-if="hasFailed" class="admin-orders-message" role="alert">
      {{ ADMIN_COPY.teamLoadFailed }}
    </p>

    <div v-else class="admin-table-scroll">
      <table class="admin-table">
        <thead>
          <tr>
            <th scope="col">{{ ADMIN_COPY.teamEmailLabel }}</th>
            <th scope="col">{{ ADMIN_COPY.teamRoleLabel }}</th>
            <th scope="col">{{ ADMIN_COPY.teamStatusLabel }}</th>
            <th scope="col">{{ ADMIN_COPY.teamAddedLabel }}</th>
            <th scope="col"><span class="visually-hidden">Actions</span></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="member in members" :key="member.id">
            <td :data-label="ADMIN_COPY.teamEmailLabel">
              {{ member.email }}
              <span v-if="isSelf(member)" class="admin-team-you">{{ ADMIN_COPY.teamYou }}</span>
            </td>

            <td :data-label="ADMIN_COPY.teamRoleLabel">
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

            <td :data-label="ADMIN_COPY.teamStatusLabel">
              <span
                class="admin-team-status"
                :class="member.active ? ADMIN_PILL_TONE.positive : ADMIN_PILL_TONE.muted"
              >
                {{ member.active ? ADMIN_COPY.teamActiveStatus : ADMIN_COPY.teamRevokedStatus }}
              </span>
            </td>

            <td :data-label="ADMIN_COPY.teamAddedLabel">{{ formatOrderDate(member.createdAt) }}</td>

            <!-- Self-targeted controls are omitted rather than disabled: the API
                 refuses them, and an inert button invites the attempt. -->
            <td class="admin-table-actions">
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
              <span v-else class="admin-detail-muted">{{ ADMIN_COPY.teamSelfHint }}</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <Teleport to="body">
      <div
        v-if="isDrawerOpen"
        class="admin-drawer-backdrop"
        role="presentation"
        @click.self="closeDrawer"
      >
        <aside
          ref="drawerElement"
          class="admin-drawer"
          role="dialog"
          aria-modal="true"
          aria-labelledby="admin-team-drawer-title"
          tabindex="-1"
        >
          <header class="admin-drawer-head">
            <div>
              <p class="admin-detail-eyebrow">{{ ADMIN_COPY.teamHeading }}</p>
              <h2 id="admin-team-drawer-title">{{ drawerTitle }}</h2>
            </div>

            <div class="admin-drawer-actions">
              <button
                ref="drawerCloseButton"
                type="button"
                class="admin-inline-action"
                @click="closeDrawer"
              >
                {{ ADMIN_COPY.closeProductEditor }}
              </button>
            </div>
          </header>

          <form
            v-if="openDrawer === 'password'"
            class="admin-settings-form"
            @submit.prevent="submitPasswordChange"
          >
            <section class="admin-form-section">
              <header class="admin-form-section-head">
                <h3>{{ ADMIN_COPY.passwordHeading }}</h3>
                <p>{{ ADMIN_COPY.passwordIntro }}</p>
              </header>

              <label class="admin-field admin-field-wide">
                <span>{{ ADMIN_COPY.currentPasswordLabel }}</span>
                <input
                  v-model="passwordForm.current"
                  type="password"
                  autocomplete="current-password"
                  required
                >
              </label>
              <label class="admin-field admin-field-wide">
                <span>{{ ADMIN_COPY.newPasswordLabel }}</span>
                <input
                  v-model="passwordForm.next"
                  type="password"
                  autocomplete="new-password"
                  minlength="12"
                  required
                >
              </label>
              <label class="admin-field admin-field-wide">
                <span>{{ ADMIN_COPY.confirmPasswordLabel }}</span>
                <input
                  v-model="passwordForm.confirm"
                  type="password"
                  autocomplete="new-password"
                  minlength="12"
                  required
                >
              </label>
            </section>

            <p v-if="passwordError" class="admin-detail-error" role="alert">{{ passwordError }}</p>

            <div class="admin-form-actions">
              <button
                type="submit"
                class="admin-inline-action is-primary"
                :disabled="isChangingPassword"
              >
                {{ isChangingPassword ? ADMIN_COPY.changingPassword : ADMIN_COPY.changePassword }}
              </button>
            </div>
          </form>

          <form v-else class="admin-settings-form" @submit.prevent="submitInvite">
            <section class="admin-form-section">
              <header class="admin-form-section-head">
                <h3>{{ ADMIN_COPY.teamInviteHeading }}</h3>
                <p>{{ ADMIN_COPY.teamInviteIntro }}</p>
              </header>

              <label class="admin-field admin-field-wide">
                <span>{{ ADMIN_COPY.teamEmailLabel }}</span>
                <input v-model.trim="inviteForm.email" type="email" autocomplete="off" required>
              </label>

              <div class="admin-field admin-field-wide">
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

              <label class="admin-field admin-field-wide">
                <span>{{ ADMIN_COPY.teamTemporaryPasswordLabel }}</span>
                <input v-model="inviteForm.password" type="text" minlength="12" required>
                <small>{{ ADMIN_COPY.teamTemporaryPasswordHint }}</small>
              </label>
            </section>

            <p v-if="inviteError" class="admin-detail-error" role="alert">{{ inviteError }}</p>

            <div class="admin-form-actions">
              <button type="submit" class="admin-inline-action is-primary" :disabled="isCreating">
                {{ isCreating ? ADMIN_COPY.teamCreating : ADMIN_COPY.teamCreate }}
              </button>
            </div>
          </form>
        </aside>
      </div>
    </Teleport>
  </div>
</template>
