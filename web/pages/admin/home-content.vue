<script setup lang="ts">
import { ADMIN_COPY, ADMIN_ROUTES } from '~/constants/admin'
import type { HomeContent, HomeContentDraft } from '~/types/settings'

definePageMeta({ layout: 'admin', middleware: 'admin' })

useHead({ title: `${ADMIN_COPY.homeContentHeading} · ${ADMIN_COPY.brand}` })

const { fetchHomeContent, updateHomeContent } = useAdminHomeContent()
const { session } = useAdminSession()
const router = useRouter()
const { notifySuccess } = useAppNotifications()
const { withLoading } = useAppLoading()

const draft = ref<HomeContentDraft>(toHomeContentDraft(defaultHomeContent()))
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

const syncDraft = (homeContent: HomeContent): void => {
  draft.value = toHomeContentDraft(homeContent)
}

const load = async (): Promise<void> => {
  isLoading.value = true
  hasFailed.value = false
  errorMessage.value = ''

  try {
    const homeContent = await withLoading(fetchHomeContent, ADMIN_COPY.homeContentLoading)
    syncDraft(homeContent)
  } catch (loadError) {
    if (!(await isExpiredSession(loadError))) {
      hasFailed.value = true
    }
  } finally {
    isLoading.value = false
  }
}

const submit = async (): Promise<void> => {
  isSaving.value = true
  errorMessage.value = ''

  try {
    const homeContent = await withLoading(
      () => updateHomeContent(draft.value),
      ADMIN_COPY.savingHomeContent,
    )
    syncDraft(homeContent)
    notifySuccess(ADMIN_COPY.homeContentSavedNotice, ADMIN_COPY.homeContentSubheading)
  } catch (saveError) {
    if (!(await isExpiredSession(saveError))) {
      errorMessage.value = ADMIN_COPY.homeContentSaveFailed
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
        <h1>{{ ADMIN_COPY.homeContentHeading }}</h1>
        <p class="admin-page-subheading">{{ ADMIN_COPY.homeContentSubheading }}</p>
      </div>
    </header>

    <p v-if="isLoading" class="admin-orders-message">{{ ADMIN_COPY.homeContentLoading }}</p>

    <div v-else-if="hasFailed" class="admin-orders-message" role="alert">
      <p>{{ ADMIN_COPY.homeContentLoadFailed }}</p>
      <button type="button" class="admin-inline-action" @click="load">
        {{ ADMIN_COPY.retry }}
      </button>
    </div>

    <form v-else class="admin-settings-form admin-home-content-form" @submit.prevent="submit">
      <section class="admin-detail-panel">
        <header class="admin-form-section-head">
          <h2>{{ ADMIN_COPY.homeHeroPanel }}</h2>
          <p>{{ ADMIN_COPY.homeHeroPanelIntro }}</p>
        </header>
        <div class="admin-detail-grid">
          <label class="admin-field">
            <span>{{ ADMIN_COPY.homeHeroLeadLabel }}</span>
            <input v-model.trim="draft.hero.headingLead" type="text" required autocomplete="off">
          </label>

          <label class="admin-field">
            <span>{{ ADMIN_COPY.homeHeroTailLabel }}</span>
            <input v-model.trim="draft.hero.headingTail" type="text" required autocomplete="off">
          </label>

          <label class="admin-field">
            <span>{{ ADMIN_COPY.homeHeroActionLabel }}</span>
            <input v-model.trim="draft.hero.actionLabel" type="text" required autocomplete="off">
          </label>

          <label class="admin-field">
            <span>{{ ADMIN_COPY.homeHeroActionToLabel }}</span>
            <input v-model.trim="draft.hero.actionTo" type="text" required autocomplete="off">
          </label>
        </div>

        <label class="admin-field admin-field-wide">
          <span>{{ ADMIN_COPY.homeHeroDescriptionLabel }}</span>
          <textarea v-model.trim="draft.hero.description" rows="4" required />
        </label>
      </section>

      <section class="admin-detail-panel">
        <header class="admin-form-section-head">
          <h2>{{ ADMIN_COPY.homeCategoriesPanel }}</h2>
          <p>{{ ADMIN_COPY.homeCategoriesPanelIntro }}</p>
        </header>
        <div class="admin-detail-grid">
          <label class="admin-field">
            <span>{{ ADMIN_COPY.homeCategoryEyebrowLabel }}</span>
            <input v-model.trim="draft.categories.eyebrow" type="text" required autocomplete="off">
          </label>

          <label class="admin-field">
            <span>{{ ADMIN_COPY.homeCategoryHeadingLabel }}</span>
            <input v-model.trim="draft.categories.heading" type="text" required autocomplete="off">
          </label>
        </div>

        <div class="admin-home-tiles">
          <fieldset
            v-for="tile in draft.categories.tiles"
            :key="tile.id"
            class="admin-home-tile"
          >
            <legend>{{ tile.heading || tile.id }}</legend>
            <div class="admin-detail-grid">
              <label class="admin-field">
                <span>{{ ADMIN_COPY.homeTileLabelLabel }}</span>
                <input v-model.trim="tile.label" type="text" required autocomplete="off">
              </label>

              <label class="admin-field">
                <span>{{ ADMIN_COPY.homeTileHeadingLabel }}</span>
                <input v-model.trim="tile.heading" type="text" required autocomplete="off">
              </label>

              <label class="admin-field">
                <span>{{ ADMIN_COPY.homeTileImagePathLabel }}</span>
                <input v-model.trim="tile.imagePath" type="text" required autocomplete="off">
              </label>

              <label class="admin-field">
                <span>{{ ADMIN_COPY.homeTileToLabel }}</span>
                <input v-model.trim="tile.to" type="text" required autocomplete="off">
              </label>
            </div>

            <label class="admin-field admin-field-wide">
              <span>{{ ADMIN_COPY.homeTileImageAltLabel }}</span>
              <input v-model.trim="tile.imageAlt" type="text" required autocomplete="off">
            </label>
          </fieldset>
        </div>
      </section>

      <p v-if="errorMessage" class="admin-detail-error" role="alert">
        {{ errorMessage }}
      </p>

      <div class="admin-settings-actions">
        <button type="submit" class="admin-inline-action" :disabled="isSaving">
          {{ isSaving ? ADMIN_COPY.savingHomeContent : ADMIN_COPY.saveHomeContent }}
        </button>
      </div>
    </form>
  </div>
</template>
