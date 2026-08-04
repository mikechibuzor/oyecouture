<script setup lang="ts">
import { ChevronRight, SlidersHorizontal } from 'lucide'
import { ADMIN_COPY } from '~/constants/admin'
import {
  AdminFilterType,
  type AdminFilterKey,
  type AdminFilterOption,
  type AdminFilterValues,
} from '~/types/admin'

interface AdminTableFilterProps {
  options: readonly AdminFilterOption[]
  values: AdminFilterValues
}

const props = defineProps<AdminTableFilterProps>()

const emit = defineEmits<{
  update: [key: AdminFilterKey, value: string]
}>()

const DESKTOP_VIEWPORT_MARGIN_PX = 16
const MOBILE_VIEWPORT_MARGIN_PX = 12
const MOBILE_POPOVER_BREAKPOINT_PX = 640
const MAX_POPOVER_WIDTH_PX = 544

const root = ref<HTMLElement | null>(null)
const trigger = ref<HTMLElement | null>(null)
const isOpen = ref(false)
const activeKey = ref<AdminFilterKey>(props.options[0]?.key ?? 'status')
const popoverStyle = ref<Record<string, string>>({})

// The right-hand pane mirrors whichever filter category is active.
const activeOption = computed(() =>
  props.options.find((option) => option.key === activeKey.value) ?? props.options[0],
)

// Date filters cannot target future orders.
const today = computed(() => new Date().toISOString().slice(0, 10))

// The popover is fixed so it can clamp to the viewport instead of its table column.
const positionPopover = (): void => {
  if (typeof window === 'undefined' || !trigger.value) {
    return
  }

  const triggerRect = trigger.value.getBoundingClientRect()
  const viewportWidth = window.innerWidth
  const viewportMargin =
    viewportWidth <= MOBILE_POPOVER_BREAKPOINT_PX
      ? MOBILE_VIEWPORT_MARGIN_PX
      : DESKTOP_VIEWPORT_MARGIN_PX
  const availableWidth = Math.max(0, viewportWidth - viewportMargin * 2)
  const popoverWidth = Math.min(MAX_POPOVER_WIDTH_PX, availableWidth)
  const preferredLeft =
    viewportWidth <= MOBILE_POPOVER_BREAKPOINT_PX ? viewportMargin : triggerRect.left
  const maxLeft = viewportWidth - popoverWidth - viewportMargin
  const clampedLeft = Math.min(Math.max(viewportMargin, preferredLeft), maxLeft)

  popoverStyle.value = {
    '--admin-filter-popover-left': `${Math.max(viewportMargin, clampedLeft)}px`,
    '--admin-filter-popover-top': `${triggerRect.bottom + 8}px`,
    '--admin-filter-popover-width': `${popoverWidth}px`,
  }
}

// Opening the menu lands on the first active filter, or the first category.
const open = (): void => {
  const activeFilter = props.options.find((option) => props.values[option.key])
  activeKey.value = activeFilter?.key ?? props.options[0]?.key ?? 'status'
  isOpen.value = true
  void nextTick(positionPopover)
}

// Closing leaves the last active category intact for the next open.
const close = (): void => {
  isOpen.value = false
}

// The trigger behaves like VBAN's filter button: one button owns the popover.
const toggle = (): void => {
  if (isOpen.value) {
    close()
    return
  }

  open()
}

// Pushes one selected value back to the page-level query state.
const updateFilter = (key: AdminFilterKey, value: string): void => {
  emit('update', key, value)
}

// Reads the native date input without widening the event contract in template.
const handleDateInput = (key: AdminFilterKey, inputEvent: Event): void => {
  updateFilter(key, (inputEvent.target as HTMLInputElement).value)
}

// Pointerdown closes before focus shifts to whatever the owner clicked next.
const handlePointerDown = (pointerEvent: PointerEvent): void => {
  if (isOpen.value && !root.value?.contains(pointerEvent.target as Node)) {
    close()
  }
}

// Keep an open menu anchored if the user scrolls or resizes with it visible.
const handleViewportChange = (): void => {
  if (isOpen.value) {
    positionPopover()
  }
}

// Escape gives keyboard users a quick way out of the open popover.
const handleKeydown = (keyboardEvent: KeyboardEvent): void => {
  if (keyboardEvent.key === 'Escape') {
    close()
  }
}

onMounted(() => {
  document.addEventListener('pointerdown', handlePointerDown)
  document.addEventListener('scroll', handleViewportChange, true)
  window.addEventListener('resize', handleViewportChange)
})

onBeforeUnmount(() => {
  document.removeEventListener('pointerdown', handlePointerDown)
  document.removeEventListener('scroll', handleViewportChange, true)
  window.removeEventListener('resize', handleViewportChange)
})
</script>

<template>
  <div ref="root" class="admin-filter" @keydown="handleKeydown">
    <button
      ref="trigger"
      type="button"
      class="admin-filter-trigger"
      :aria-expanded="isOpen"
      aria-haspopup="menu"
      @click="toggle"
    >
      <StorefrontIcon :icon-node="SlidersHorizontal" :size="16" :stroke-width="1.75" />
      {{ ADMIN_COPY.filter }}
    </button>

    <div v-show="isOpen" class="admin-filter-popover" role="menu" :style="popoverStyle">
      <div class="admin-filter-categories">
        <button
          v-for="option in options"
          :key="option.key"
          type="button"
          class="admin-filter-option"
          :class="{ 'is-active': option.key === activeKey, 'has-value': values[option.key] }"
          role="menuitem"
          @click="activeKey = option.key"
        >
          <span>{{ option.label }}</span>
          <StorefrontIcon :icon-node="ChevronRight" :size="14" :stroke-width="1.75" />
        </button>
      </div>

      <div v-if="activeOption" class="admin-filter-values">
        <template v-if="activeOption.type === AdminFilterType.Select">
          <button
            type="button"
            class="admin-filter-value"
            :class="{ 'is-selected': values[activeOption.key] === '' }"
            @click="updateFilter(activeOption.key, '')"
          >
            {{ ADMIN_COPY.filterEmpty }}
          </button>
          <button
            v-for="item in activeOption.items"
            :key="item.value"
            type="button"
            class="admin-filter-value"
            :class="{ 'is-selected': values[activeOption.key] === item.value }"
            @click="updateFilter(activeOption.key, item.value)"
          >
            {{ item.label }}
          </button>
        </template>

        <template v-else>
          <label class="admin-filter-date">
            <span>{{ activeOption.label }}</span>
            <input
              type="date"
              :value="values[activeOption.key]"
              :max="today"
              @input="handleDateInput(activeOption.key, $event)"
            >
          </label>
          <button
            v-if="values[activeOption.key]"
            type="button"
            class="admin-filter-value"
            @click="updateFilter(activeOption.key, '')"
          >
            {{ ADMIN_COPY.clearFilter }}
          </button>
        </template>
      </div>
    </div>
  </div>
</template>
