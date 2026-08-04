<script setup lang="ts">
import {
  ADMIN_COPY,
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

const submit = (): void => {
  emit('submit')
}
</script>

<template>
  <form class="admin-product-form" @submit.prevent="submit">
    <div class="admin-detail-grid">
      <label class="admin-field">
        <span>{{ ADMIN_COPY.slugLabel }}</span>
        <input v-model.trim="draft.slug" type="text" required autocomplete="off">
      </label>

      <label class="admin-field">
        <span>{{ ADMIN_COPY.nameLabel }}</span>
        <input v-model.trim="draft.name" type="text" required autocomplete="off">
      </label>

      <div class="admin-field">
        <span>{{ ADMIN_COPY.productSectionLabel }}</span>
        <UiSelect
          v-model="draft.section"
          :options="PRODUCT_FORM_SECTION_OPTIONS"
          :label="ADMIN_COPY.productSectionLabel"
        />
      </div>

      <label class="admin-field">
        <span>{{ ADMIN_COPY.categoryLabel }}</span>
        <input v-model.trim="draft.category" type="text" required autocomplete="off">
      </label>

      <label class="admin-field">
        <span>{{ ADMIN_COPY.priceKoboLabel }}</span>
        <input v-model.number="draft.priceKobo" type="number" min="0" step="1" required>
      </label>

      <label class="admin-field">
        <span>{{ ADMIN_COPY.sortOrderLabel }}</span>
        <input v-model.number="draft.sortOrder" type="number" min="0" step="1" required>
      </label>

      <label class="admin-field">
        <span>{{ ADMIN_COPY.tagLabel }}</span>
        <input v-model.trim="draft.tag" type="text" autocomplete="off">
      </label>

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
    </div>

    <label class="admin-field admin-field-wide">
      <span>{{ ADMIN_COPY.descriptionLabel }}</span>
      <textarea v-model.trim="draft.description" rows="5" required />
    </label>

    <p v-if="errorMessage" class="admin-detail-error" role="alert">
      {{ errorMessage }}
    </p>
    <p v-if="successMessage" class="admin-detail-muted" role="status">
      {{ successMessage }}
    </p>

    <div class="admin-product-form-actions">
      <button type="submit" class="admin-inline-action" :disabled="isSaving">
        {{ isSaving ? ADMIN_COPY.savingProduct : ADMIN_COPY.saveProduct }}
      </button>
    </div>
  </form>
</template>
