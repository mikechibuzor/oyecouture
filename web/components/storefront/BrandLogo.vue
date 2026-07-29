<script setup lang="ts">
import { STOREFRONT_COPY } from '~/constants/storefront'
import { StorefrontRoute } from '~/types/storefront'
</script>

<template>
  <NuxtLink class="brand-logo" :to="StorefrontRoute.Home">
    <img
      :src="STOREFRONT_COPY.brandLogoPath"
      :alt="STOREFRONT_COPY.brandAlt"
      width="636"
      height="122"
    >
  </NuxtLink>
</template>
