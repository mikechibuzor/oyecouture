<script setup lang="ts">
import { ADMIN_COPY, ADMIN_ROUTES } from '~/constants/admin'
import type { SiteSettings, SiteSettingsDraft } from '~/types/settings'

definePageMeta({ layout: 'admin', middleware: 'admin' })

useHead({ title: `${ADMIN_COPY.settingsHeading} · ${ADMIN_COPY.brand}` })

const { fetchSettings, updateSettings } = useAdminSettings()
const { session } = useAdminSession()
const router = useRouter()
const { notifySuccess } = useAppNotifications()
const { withLoading } = useAppLoading()

const draft = ref<SiteSettingsDraft>({
  announcement: '',
  whatsappNumber: '',
  held: false,
  holdHeadline: '',
  holdMessage: '',
})
const isLoading = ref(true)
const hasFailed = ref(false)
const isSaving = ref(false)
const errorMessage = ref('')

const isExpiredSession = async (requestError: unknown): Promise<boolean> => {
  if ((requestError as { status?: number } | null)?.status !== 401) {
    return false
  }

  session.value = null
  await router.push(ADMIN_ROUTES.login)

  return true
}

const syncDraft = (settings: SiteSettings): void => {
  draft.value = {
    announcement: settings.announcement,
    whatsappNumber: settings.whatsappNumber,
    held: settings.held,
    holdHeadline: settings.holdHeadline,
    holdMessage: settings.holdMessage,
  }
}

// The API refuses a hold with no headline; saying so here means the owner is not told by a
// failed save what a disabled button could have told them first.
const isMissingHoldHeadline = computed(
  () => draft.value.held && draft.value.holdHeadline.trim() === '',
)

const load = async (): Promise<void> => {
  isLoading.value = true
  hasFailed.value = false
  errorMessage.value = ''

  try {
    const settings = await withLoading(fetchSettings, ADMIN_COPY.settingsLoading)
    syncDraft(settings)
  } catch (loadError) {
    if (!(await isExpiredSession(loadError))) {
      hasFailed.value = true
    }
  } finally {
    isLoading.value = false
  }
}

const submit = async (): Promise<void> => {
  if (isMissingHoldHeadline.value) {
    errorMessage.value = ADMIN_COPY.holdHeadlineRequired

    return
  }

  isSaving.value = true
  errorMessage.value = ''

  try {
    const settings = await withLoading(
      () => updateSettings(draft.value),
      ADMIN_COPY.savingSettings,
    )
    syncDraft(settings)
    notifySuccess(ADMIN_COPY.settingsSavedNotice, ADMIN_COPY.settingsSubheading)
  } catch (saveError) {
    if (!(await isExpiredSession(saveError))) {
      errorMessage.value = ADMIN_COPY.settingsSaveFailed
    }
  } finally {
    isSaving.value = false
  }
}

onMounted(() => {
  void load()
})
</script>

<template>
  <div class="admin-page">
    <header class="admin-page-head">
      <div>
        <h1>{{ ADMIN_COPY.settingsHeading }}</h1>
        <p class="admin-page-subheading">{{ ADMIN_COPY.settingsSubheading }}</p>
      </div>
    </header>

    <p v-if="isLoading" class="admin-orders-message">{{ ADMIN_COPY.settingsLoading }}</p>

    <div v-else-if="hasFailed" class="admin-orders-message" role="alert">
      <p>{{ ADMIN_COPY.settingsLoadFailed }}</p>
      <button type="button" class="admin-inline-action" @click="load">
        {{ ADMIN_COPY.retry }}
      </button>
    </div>

    <form v-else class="admin-settings-form" @submit.prevent="submit">
      <section class="admin-form-section" aria-labelledby="admin-settings-panel-heading">
        <header class="admin-form-section-head">
          <h2 id="admin-settings-panel-heading">{{ ADMIN_COPY.settingsPanelHeading }}</h2>
          <p>{{ ADMIN_COPY.settingsPanelIntro }}</p>
        </header>

        <label class="admin-field admin-field-wide">
          <span>{{ ADMIN_COPY.announcementLabel }}</span>
          <input
            v-model.trim="draft.announcement"
            type="text"
            maxlength="180"
            autocomplete="off"
          >
          <small>{{ ADMIN_COPY.announcementHint }}</small>
        </label>

        <label class="admin-field admin-field-wide">
          <span>{{ ADMIN_COPY.whatsappNumberLabel }}</span>
          <input
            v-model.trim="draft.whatsappNumber"
            type="tel"
            autocomplete="tel"
          >
          <small>{{ ADMIN_COPY.whatsappNumberHint }}</small>
        </label>
      </section>

      <section class="admin-form-section" aria-labelledby="admin-hold-panel-heading">
        <header class="admin-form-section-head">
          <h2 id="admin-hold-panel-heading">{{ ADMIN_COPY.holdPanelHeading }}</h2>
          <p>{{ ADMIN_COPY.holdPanelIntro }}</p>
        </header>

        <label class="admin-field admin-field-wide admin-field-toggle">
          <input v-model="draft.held" type="checkbox">
          <span>{{ ADMIN_COPY.holdToggleLabel }}</span>
          <small>{{ ADMIN_COPY.holdToggleHint }}</small>
        </label>

        <label class="admin-field admin-field-wide">
          <span>{{ ADMIN_COPY.holdHeadlineLabel }}</span>
          <input
            v-model.trim="draft.holdHeadline"
            type="text"
            maxlength="70"
            autocomplete="off"
          >
          <small>{{ ADMIN_COPY.holdHeadlineHint }}</small>
        </label>

        <label class="admin-field admin-field-wide">
          <span>{{ ADMIN_COPY.holdMessageLabel }}</span>
          <textarea
            v-model.trim="draft.holdMessage"
            rows="3"
            maxlength="220"
          />
          <small>{{ ADMIN_COPY.holdMessageHint }}</small>
        </label>
      </section>

      <p v-if="errorMessage" class="admin-detail-error" role="alert">
        {{ errorMessage }}
      </p>

      <div class="admin-settings-actions">
        <button
          type="submit"
          class="admin-inline-action"
          :disabled="isSaving || isMissingHoldHeadline"
        >
          {{ isSaving ? ADMIN_COPY.savingSettings : ADMIN_COPY.saveSettings }}
        </button>
      </div>
    </form>
  </div>
</template>
