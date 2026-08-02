<script setup lang="ts" generic="TValue extends string">
import { ChevronDown } from 'lucide'
import type { SelectOption } from '~/types/ui'

// A styled dropdown replacing the native control, which cannot be themed
// consistently across browsers. It follows the ARIA combobox pattern: focus
// stays on the trigger and the active option is tracked with
// aria-activedescendant, so screen readers announce the choice without the
// focus moving into the list.
interface UiSelectProps {
  options: readonly SelectOption<TValue>[]
  label: string
  // Shown when the model is empty. Selecting it is only possible when the
  // field is optional.
  placeholder?: string
  disabled?: boolean
  // Renders the trigger without a border, for toolbar use.
  variant?: 'field' | 'bare'
}

const props = withDefaults(defineProps<UiSelectProps>(), {
  placeholder: '',
  disabled: false,
  variant: 'field',
})

const selected = defineModel<TValue>({ required: true })

const listboxId = useId()
const root = ref<HTMLElement | null>(null)
const listbox = ref<HTMLElement | null>(null)
const isOpen = ref(false)
const activeIndex = ref(-1)

const selectableIndexes = computed(() =>
  props.options.reduce<number[]>((indexes, option, index) => {
    if (!option.disabled) {
      indexes.push(index)
    }

    return indexes
  }, []),
)

const selectedIndex = computed(() =>
  props.options.findIndex((option) => option.value === selected.value),
)

const triggerLabel = computed(
  () => props.options[selectedIndex.value]?.label ?? props.placeholder,
)

const optionId = (index: number): string => `${listboxId}-option-${index}`

// Opening lands on the current choice so arrow keys continue from there.
const open = (): void => {
  if (props.disabled) {
    return
  }

  isOpen.value = true
  activeIndex.value =
    selectedIndex.value >= 0 ? selectedIndex.value : (selectableIndexes.value[0] ?? -1)
}

const close = (): void => {
  isOpen.value = false
  activeIndex.value = -1
}

const choose = (index: number): void => {
  const option = props.options[index]
  if (!option || option.disabled) {
    return
  }

  selected.value = option.value
  close()
}

// Steps through the selectable options, skipping any that are disabled.
const moveActive = (step: number): void => {
  const selectable = selectableIndexes.value
  if (selectable.length === 0) {
    return
  }

  const currentPosition = selectable.indexOf(activeIndex.value)
  const nextPosition =
    currentPosition === -1
      ? 0
      : (currentPosition + step + selectable.length) % selectable.length

  activeIndex.value = selectable[nextPosition] as number
}

// Type-ahead: typing "de" jumps to Delivered, the way a native select does.
let typedPrefix = ''
let typedTimer: ReturnType<typeof setTimeout> | undefined

const jumpToPrefix = (character: string): void => {
  clearTimeout(typedTimer)
  typedPrefix += character.toLowerCase()
  typedTimer = setTimeout(() => {
    typedPrefix = ''
  }, 600)

  const match = selectableIndexes.value.find((index) =>
    props.options[index]?.label.toLowerCase().startsWith(typedPrefix),
  )

  if (match !== undefined) {
    activeIndex.value = match
    if (!isOpen.value) {
      choose(match)
    }
  }
}

const handleKeydown = (keyboardEvent: KeyboardEvent): void => {
  const { key } = keyboardEvent

  if (key === 'Escape') {
    close()
    return
  }

  if (key === 'Tab') {
    close()
    return
  }

  if (key === 'ArrowDown' || key === 'ArrowUp') {
    keyboardEvent.preventDefault()
    if (!isOpen.value) {
      open()
      return
    }
    moveActive(key === 'ArrowDown' ? 1 : -1)
    return
  }

  if (key === 'Home' || key === 'End') {
    if (!isOpen.value) {
      return
    }
    keyboardEvent.preventDefault()
    const selectable = selectableIndexes.value
    activeIndex.value = (
      key === 'Home' ? selectable[0] : selectable[selectable.length - 1]
    ) as number
    return
  }

  if (key === 'Enter' || key === ' ') {
    keyboardEvent.preventDefault()
    if (isOpen.value) {
      choose(activeIndex.value)
    } else {
      open()
    }
    return
  }

  // Single printable characters drive type-ahead; modifier combos must not.
  if (key.length === 1 && !keyboardEvent.metaKey && !keyboardEvent.ctrlKey) {
    jumpToPrefix(key)
  }
}

// Keep the active option visible when arrowing through a long list.
watch(activeIndex, async (index) => {
  if (!isOpen.value || index < 0) {
    return
  }

  await nextTick()
  listbox.value
    ?.querySelector(`#${CSS.escape(optionId(index))}`)
    ?.scrollIntoView({ block: 'nearest' })
})

// Clicking away closes the list. Pointerdown rather than click, so the list is
// gone before a click lands on whatever is underneath.
const handlePointerDown = (pointerEvent: PointerEvent): void => {
  if (isOpen.value && !root.value?.contains(pointerEvent.target as Node)) {
    close()
  }
}

onMounted(() => document.addEventListener('pointerdown', handlePointerDown))
onBeforeUnmount(() => {
  document.removeEventListener('pointerdown', handlePointerDown)
  clearTimeout(typedTimer)
})
</script>

<template>
  <div ref="root" class="ui-select" :class="[`is-${variant}`, { 'is-open': isOpen }]">
    <button
      type="button"
      class="ui-select-trigger"
      role="combobox"
      :aria-expanded="isOpen"
      :aria-controls="listboxId"
      aria-haspopup="listbox"
      :aria-activedescendant="isOpen && activeIndex >= 0 ? optionId(activeIndex) : undefined"
      :aria-label="label"
      :disabled="disabled"
      @click="isOpen ? close() : open()"
      @keydown="handleKeydown"
    >
      <span class="ui-select-value" :class="{ 'is-placeholder': selectedIndex < 0 }">
        {{ triggerLabel }}
      </span>
      <StorefrontIcon
        class="ui-select-chevron"
        :icon-node="ChevronDown"
        :size="16"
        :stroke-width="1.75"
      />
    </button>

    <ul
      v-show="isOpen"
      :id="listboxId"
      ref="listbox"
      class="ui-select-list"
      role="listbox"
      :aria-label="label"
    >
      <li
        v-for="(option, index) in options"
        :id="optionId(index)"
        :key="option.value"
        class="ui-select-option"
        :class="{
          'is-active': index === activeIndex,
          'is-selected': option.value === selected,
          'is-disabled': option.disabled,
        }"
        role="option"
        :aria-selected="option.value === selected"
        :aria-disabled="option.disabled"
        @click="choose(index)"
        @pointermove="option.disabled || (activeIndex = index)"
      >
        {{ option.label }}
      </li>
    </ul>
  </div>
</template>
