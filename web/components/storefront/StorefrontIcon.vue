<script setup lang="ts">
import type { IconNode } from 'lucide'

interface StorefrontIconProps {
  iconNode: IconNode
  size: number
  strokeWidth: number
}

defineProps<StorefrontIconProps>()
</script>

<template>
  <svg
    class="lucide"
    xmlns="http://www.w3.org/2000/svg"
    :width="size"
    :height="size"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    :stroke-width="strokeWidth"
    stroke-linecap="round"
    stroke-linejoin="round"
  >
    <component
      :is="iconElementName"
      v-for="([iconElementName, iconElementAttributes], iconElementIndex) in iconNode"
      :key="iconElementIndex"
      v-bind="iconElementAttributes"
    />
  </svg>
</template>
