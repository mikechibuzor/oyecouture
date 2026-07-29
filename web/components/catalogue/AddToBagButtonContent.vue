<script setup lang="ts">
import type { IconNode } from 'lucide'
import { ShoppingBag } from 'lucide'
import {
  PRODUCT_BUTTON_FIZZ_PARTICLE_COUNT,
  PRODUCT_COPY,
} from '~/constants/products'

interface AddToBagButtonContentProps {
  iconSize?: number
}

withDefaults(defineProps<AddToBagButtonContentProps>(), {
  iconSize: 18,
})

const ADD_TO_BAG_ICON: IconNode = ShoppingBag
</script>

<template>
  <span class="add-to-bag-button-content">
    <span class="add-to-bag-fizz" aria-hidden="true">
      <span
        v-for="particleIndex in PRODUCT_BUTTON_FIZZ_PARTICLE_COUNT"
        :key="particleIndex"
        class="add-to-bag-fizz-particle"
      />
    </span>
    <StorefrontIcon
      class="add-to-bag-button-icon"
      :icon-node="ADD_TO_BAG_ICON"
      :size="iconSize"
      :stroke-width="1.5"
      aria-hidden="true"
    />
    <span class="add-to-bag-button-label">
      {{ PRODUCT_COPY.actions.addToBag }}
    </span>
  </span>
</template>
