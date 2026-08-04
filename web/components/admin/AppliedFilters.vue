<script setup lang="ts">
import { X } from 'lucide'
import { ADMIN_COPY } from '~/constants/admin'
import {
  AdminFilterType,
  type AdminFilterKey,
  type AdminFilterOption,
  type AdminFilterValues,
} from '~/types/admin'

interface AdminAppliedFiltersProps {
  options: readonly AdminFilterOption[]
  values: AdminFilterValues
}

const props = defineProps<AdminAppliedFiltersProps>()

const emit = defineEmits<{
  remove: [key: AdminFilterKey]
  clear: []
}>()

// Only non-empty filter values should render as chips.
const appliedFilters = computed(() =>
  props.options.filter((option) => props.values[option.key] !== ''),
)

// Each chip resolves the stored value through its filter option for readable UI.
const valueLabel = (option: AdminFilterOption): string => {
  const value = props.values[option.key]
  if (option.type === AdminFilterType.Date) {
    return value
  }

  return option.items?.find((item) => item.value === value)?.label ?? value
}
</script>

<template>
  <div
    v-if="appliedFilters.length"
    class="admin-applied-filters"
    :aria-label="ADMIN_COPY.appliedFiltersLabel"
  >
    <span v-for="option in appliedFilters" :key="option.key" class="admin-applied-filter">
      {{ option.label }}: <strong>{{ valueLabel(option) }}</strong>
      <button
        type="button"
        :aria-label="`${ADMIN_COPY.clearFilter}: ${option.label}`"
        @click="emit('remove', option.key)"
      >
        <StorefrontIcon :icon-node="X" :size="13" :stroke-width="1.75" />
      </button>
    </span>

    <button type="button" class="admin-applied-clear" @click="emit('clear')">
      {{ ADMIN_COPY.clearFilters }}
    </button>
  </div>
</template>
