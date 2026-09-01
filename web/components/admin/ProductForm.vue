<script setup lang="ts">
import {
  ADMIN_COPY,
  PRODUCT_CATEGORY_OPTIONS_BY_SECTION,
  PRODUCT_FORM_SECTION_OPTIONS,
} from '~/constants/admin'
import type { AdminProductDraft } from '~/types/admin'

interface AdminProductFormProps {
  isSaving: boolean
  errorMessage: string
  successMessage: string
}

defineProps<AdminProductFormProps>()

const draft = defineModel<AdminProductDraft>({ required: true })

const emit = defineEmits<{ submit: [] }>()

const categoryOptions = computed(() => {
  const sectionOptions = PRODUCT_CATEGORY_OPTIONS_BY_SECTION[draft.value.section]
  const hasCurrentCategory = sectionOptions.some(
    (option) => option.value === draft.value.category,
  )

  if (!draft.value.category || hasCurrentCategory) {
    return sectionOptions
  }

  return [
    { value: draft.value.category, label: draft.value.category },
    ...sectionOptions,
  ]
})

watch(
  () => draft.value.section,
  () => {
    const hasCategory = categoryOptions.value.some(
      (option) => option.value === draft.value.category,
    )

    if (!hasCategory) {
      draft.value.category = categoryOptions.value[0]?.value ?? ''
    }
  },
  { immediate: true },
)

const submit = (): void => {
  emit('submit')
}
</script>

<template>
  <form class="admin-product-form" @submit.prevent="submit">
    <section class="admin-form-section" aria-labelledby="admin-product-identity-heading">
      <header class="admin-form-section-head">
        <h2 id="admin-product-identity-heading">{{ ADMIN_COPY.productIdentityHeading }}</h2>
        <p>{{ ADMIN_COPY.productIdentityIntro }}</p>
      </header>

      <div class="admin-detail-grid">
        <label class="admin-field">
          <span>{{ ADMIN_COPY.nameLabel }}</span>
          <input v-model.trim="draft.name" type="text" required autocomplete="off">
        </label>

        <label class="admin-field">
          <span>{{ ADMIN_COPY.slugLabel }}</span>
          <input v-model.trim="draft.slug" type="text" required autocomplete="off">
        </label>
      </div>

      <label class="admin-field admin-field-wide">
        <span>{{ ADMIN_COPY.descriptionLabel }}</span>
        <textarea v-model.trim="draft.description" rows="5" required />
      </label>
    </section>

    <section class="admin-form-section" aria-labelledby="admin-product-catalogue-heading">
      <header class="admin-form-section-head">
        <h2 id="admin-product-catalogue-heading">{{ ADMIN_COPY.productCatalogueHeading }}</h2>
        <p>{{ ADMIN_COPY.productCatalogueIntro }}</p>
      </header>

      <div class="admin-detail-grid">
        <div class="admin-field">
          <span>{{ ADMIN_COPY.productSectionLabel }}</span>
          <UiSelect
            v-model="draft.section"
            :options="PRODUCT_FORM_SECTION_OPTIONS"
            :label="ADMIN_COPY.productSectionLabel"
          />
        </div>

        <div class="admin-field">
          <span>{{ ADMIN_COPY.categoryLabel }}</span>
          <UiSelect
            v-model="draft.category"
            :options="categoryOptions"
            :label="ADMIN_COPY.categoryLabel"
          />
        </div>

        <label class="admin-field">
          <span>{{ ADMIN_COPY.priceNairaLabel }}</span>
          <input v-model.number="draft.priceNaira" type="number" min="0" step="1" required>
        </label>

        <label class="admin-field">
          <span>{{ ADMIN_COPY.sortOrderLabel }}</span>
          <input v-model.number="draft.sortOrder" type="number" min="0" step="1" required>
        </label>

        <label class="admin-field">
          <span>{{ ADMIN_COPY.tagLabel }}</span>
          <input v-model.trim="draft.tag" type="text" autocomplete="off">
        </label>
      </div>
    </section>

    <section class="admin-form-section" aria-labelledby="admin-product-availability-heading">
      <header class="admin-form-section-head">
        <h2 id="admin-product-availability-heading">
          {{ ADMIN_COPY.productAvailabilityHeading }}
        </h2>
        <p>{{ ADMIN_COPY.productAvailabilityIntro }}</p>
      </header>

      <div class="admin-product-toggles">
        <label>
          <input v-model="draft.hasSizes" type="checkbox">
          <span>{{ ADMIN_COPY.hasSizesLabel }}</span>
        </label>

        <label>
          <input v-model="draft.active" type="checkbox">
          <span>{{ ADMIN_COPY.activeLabel }}</span>
        </label>
      </div>
    </section>

    <p v-if="errorMessage" class="admin-detail-error" role="alert">
      {{ errorMessage }}
    </p>
    <p v-if="successMessage" class="admin-detail-muted" role="status">
      {{ successMessage }}
    </p>

    <div class="admin-product-form-actions admin-form-actions">
      <button type="submit" class="admin-inline-action" :disabled="isSaving">
        {{ isSaving ? ADMIN_COPY.savingProduct : ADMIN_COPY.saveProduct }}
      </button>
    </div>
  </form>
</template>
