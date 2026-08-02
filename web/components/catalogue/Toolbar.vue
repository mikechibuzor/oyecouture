<script setup lang="ts">
import {
  CATALOGUE_FILTER_OPTIONS,
  PRODUCT_COPY,
  PRODUCT_SORT_OPTIONS,
} from '~/constants/products'
import { CatalogueFilterId, ProductSortId } from '~/types/product'
import type { SelectOption } from '~/types/ui'

interface CatalogueToolbarProps {
  pieceCountLabel: string
  showFilters?: boolean
}

withDefaults(defineProps<CatalogueToolbarProps>(), {
  showFilters: true,
})

const activeFilterId = defineModel<CatalogueFilterId>('activeFilterId', {
  required: true,
})
const activeSortId = defineModel<ProductSortId>('activeSortId', {
  required: true,
})

const sortSelectOptions = computed<SelectOption<ProductSortId>[]>(() =>
  PRODUCT_SORT_OPTIONS.map((sortOption) => ({
    value: sortOption.id,
    label: sortOption.label,
  })),
)

// Activates one catalogue filter chip.
const selectFilter = (catalogueFilterId: CatalogueFilterId): void => {
  activeFilterId.value = catalogueFilterId
}
</script>

<template>
  <div class="catalogue-toolbar">
    <div
      v-if="showFilters"
      class="catalogue-filter-list"
      role="group"
      :aria-label="PRODUCT_COPY.catalogue.filterLabel"
    >
      <button
        v-for="filterOption in CATALOGUE_FILTER_OPTIONS"
        :key="filterOption.id"
        type="button"
        :class="{ 'is-active': activeFilterId === filterOption.id }"
        :aria-pressed="activeFilterId === filterOption.id"
        @click="selectFilter(filterOption.id)"
      >
        {{ filterOption.label }}
      </button>
    </div>

    <div class="catalogue-sort-group">
      <span>{{ pieceCountLabel }}</span>
      <UiSelect
        v-model="activeSortId"
        :options="sortSelectOptions"
        :label="PRODUCT_COPY.catalogue.sortLabel"
        variant="bare"
      />
    </div>
  </div>
</template>
