<script setup lang="ts">
import {
  ADMIN_COPY,
  ADMIN_METRIC_CARDS,
  METRICS_PERIOD_LABELS,
  METRICS_PERIOD_OPTIONS,
} from '~/constants/admin'
import type {
  AdminMetricCard,
  AdminOrderMetrics,
  MetricsPeriod,
  OrderStatus,
} from '~/types/admin'
import type { SelectOption } from '~/types/ui'
import { formatProductPrice } from '~/utils/product'

interface AdminMetricCardsProps {
  metrics: AdminOrderMetrics | null
  activeStatus: OrderStatus | ''
  period: MetricsPeriod
}

const props = defineProps<AdminMetricCardsProps>()

const emit = defineEmits<{
  select: [statusFilter: OrderStatus | null]
  'update:period': [period: MetricsPeriod]
}>()

const periodOptions = computed<SelectOption<MetricsPeriod>[]>(() =>
  METRICS_PERIOD_OPTIONS.map((option) => ({
    value: option,
    label: METRICS_PERIOD_LABELS[option],
  })),
)

// Revenue reads as money; every other card is a count.
const valueFor = (card: AdminMetricCard): string => {
  if (!props.metrics) {
    return '—'
  }

  const value = props.metrics[card.id as keyof AdminOrderMetrics]

  return card.isMoney
    ? formatProductPrice(value)
    : new Intl.NumberFormat('en-NG').format(value)
}

// The revenue card is a readout, not a filter, so it is not selectable.
const isSelectable = (card: AdminMetricCard): boolean => !card.isMoney

// "All orders" is current whenever no status filter is applied.
const isCurrent = (card: AdminMetricCard): boolean =>
  isSelectable(card) && (card.statusFilter ?? '') === props.activeStatus
</script>

<template>
  <div class="admin-metrics-period">
    <UiSelect
      :model-value="period"
      :options="periodOptions"
      :label="ADMIN_COPY.periodLabel"
      @update:model-value="emit('update:period', $event)"
    />
  </div>

  <ul class="admin-metrics">
    <li v-for="card in ADMIN_METRIC_CARDS" :key="card.id">
      <component
        :is="isSelectable(card) ? 'button' : 'div'"
        class="admin-metric"
        :class="{ 'is-current': isCurrent(card), 'is-static': !isSelectable(card) }"
        :type="isSelectable(card) ? 'button' : undefined"
        :aria-pressed="isSelectable(card) ? isCurrent(card) : undefined"
        @click="isSelectable(card) && emit('select', card.statusFilter)"
      >
        <span class="admin-metric-label">{{ card.label }}</span>
        <span class="admin-metric-value">{{ valueFor(card) }}</span>
      </component>
    </li>
  </ul>
</template>
