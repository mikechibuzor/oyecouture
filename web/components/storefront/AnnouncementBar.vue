<script setup lang="ts">
import { STOREFRONT_COPY } from '~/constants/storefront'
</script>

<template>
  <aside
    class="announcement-bar"
    :aria-label="STOREFRONT_COPY.announcementLabel"
  >
    <p>{{ STOREFRONT_COPY.announcement }}</p>
  </aside>
</template>
