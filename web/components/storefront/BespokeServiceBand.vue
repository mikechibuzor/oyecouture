<script setup lang="ts">
import { BESPOKE_SERVICE_COPY } from '~/constants/storefront'

// The About page runs its own heading and adds a line of copy; every other page
// uses the shared wording.
interface BespokeServiceBandProps {
  heading?: string
  description?: string
}

withDefaults(defineProps<BespokeServiceBandProps>(), {
  heading: BESPOKE_SERVICE_COPY.heading,
  description: '',
})

const { appointmentWhatsAppUrl, opensWhatsApp } = await useAppointmentWhatsApp()
</script>

<template>
  <section id="bespoke" class="bespoke-service-band" aria-labelledby="bespoke-service-title">
    <img
      class="bespoke-service-mark"
      src="/images/brand/oyecouture-icon-gold.svg"
      alt=""
      width="36"
      height="36"
      aria-hidden="true"
    >

    <p class="eyebrow">{{ BESPOKE_SERVICE_COPY.eyebrow }}</p>
    <h2 id="bespoke-service-title">{{ heading }}</h2>
    <p v-if="description" class="bespoke-service-description">{{ description }}</p>

    <a
      class="bespoke-service-action"
      :href="appointmentWhatsAppUrl"
      :target="opensWhatsApp ? '_blank' : undefined"
      :rel="opensWhatsApp ? 'noopener' : undefined"
    >
      {{ BESPOKE_SERVICE_COPY.action }}
    </a>
  </section>
</template>
