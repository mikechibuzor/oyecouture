<script setup lang="ts">
interface MenuToggleIconProps {
  isOpen: boolean
}

defineProps<MenuToggleIconProps>()
</script>

<template>
  <svg
    class="menu-toggle-icon"
    :class="{ 'is-open': isOpen }"
    viewBox="0 0 28 28"
    fill="none"
    aria-hidden="true"
  >
    <circle
      class="menu-toggle-ring"
      cx="14"
      cy="14"
      r="12"
      vector-effect="non-scaling-stroke"
    />
    <path
      class="menu-toggle-line is-top"
      d="M7 9H21"
      vector-effect="non-scaling-stroke"
    />
    <path
      class="menu-toggle-line is-middle"
      d="M7 14H21"
      vector-effect="non-scaling-stroke"
    />
    <path
      class="menu-toggle-line is-bottom"
      d="M7 19H21"
      vector-effect="non-scaling-stroke"
    />
  </svg>
</template>
