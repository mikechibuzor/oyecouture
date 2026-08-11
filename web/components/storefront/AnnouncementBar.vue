<script setup lang="ts">
import { STOREFRONT_COPY } from '~/constants/storefront'

const { data: siteSettings } = await useSiteSettings()

const announcement = computed(() =>
  siteSettings.value?.announcement || STOREFRONT_COPY.announcement,
)
</script>

<template>
  <aside
    class="announcement-bar"
    :aria-label="STOREFRONT_COPY.announcementLabel"
  >
    <p>{{ announcement }}</p>
  </aside>
</template>
