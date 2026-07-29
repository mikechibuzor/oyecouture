<script setup lang="ts">
import type { IconNode } from 'lucide'
import { MessageCircle } from 'lucide'
import { PRODUCT_COPY } from '~/constants/products'

const BESPOKE_NOTE_ICON: IconNode = MessageCircle
</script>

<template>
  <aside class="bespoke-note">
    <StorefrontIcon
      :icon-node="BESPOKE_NOTE_ICON"
      :size="20"
      :stroke-width="1.5"
      aria-hidden="true"
    />
    <p>{{ PRODUCT_COPY.bespokeNote }}</p>
  </aside>
</template>
