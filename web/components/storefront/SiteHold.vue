<script setup lang="ts">
import { STOREFRONT_COPY } from '~/constants/storefront'
import { SITE_HOLD_COPY } from '~/constants/site-hold'
import type { SiteSettings } from '~/types/settings'
import { buildWhatsAppUrl } from '~/utils/whatsapp'

interface SiteHoldProps {
  settings: SiteSettings | null
}

const props = defineProps<SiteHoldProps>()

// The owner's own words, with a usable line kept behind each in case the shop is held
// before anything has been written.
const headline = computed(
  () => props.settings?.holdHeadline?.trim() || SITE_HOLD_COPY.fallbackHeadline,
)

const message = computed(() => props.settings?.holdMessage?.trim() ?? '')

// The number already kept in settings, so the holding page needs no second place to set it.
const whatsAppUrl = computed(() =>
  buildWhatsAppUrl(String(props.settings?.whatsappNumber ?? ''), SITE_HOLD_COPY.whatsAppMessage),
)

// A held shop that search engines index is what people find long after it reopens.
useHead({
  title: `${STOREFRONT_COPY.brandAlt} — ${SITE_HOLD_COPY.documentTitle}`,
  meta: [{ name: 'robots', content: 'noindex, nofollow' }],
})
</script>

<template>
  <main class="site-hold">
    <div class="site-hold-inner">
      <img
        class="site-hold-mark"
        :src="STOREFRONT_COPY.brandLogoPath"
        :alt="STOREFRONT_COPY.brandAlt"
        width="636"
        height="122"
      >

      <span class="site-hold-rule" aria-hidden="true" />

      <h1 class="site-hold-headline">{{ headline }}</h1>
      <p v-if="message" class="site-hold-message">{{ message }}</p>

      <a
        v-if="whatsAppUrl"
        class="site-hold-action"
        :href="whatsAppUrl"
        target="_blank"
        rel="noopener noreferrer"
      >
        {{ SITE_HOLD_COPY.whatsAppAction }}
      </a>

      <span class="site-hold-rule" aria-hidden="true" />
    </div>
  </main>
</template>

<style scoped lang="scss">
.site-hold {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  min-height: 100dvh;
  padding: 48px 24px;
  background: var(--color-canvas);
  color: var(--color-ink);
}

.site-hold-inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 28px;
  width: 100%;
  max-width: 560px;
  text-align: center;
}

.site-hold-mark {
  width: min(228px, 62vw);
  height: auto;
}

// The rules draw themselves in, which is the one piece of movement on an otherwise still page.
.site-hold-rule {
  width: 100%;
  max-width: 320px;
  height: 1px;
  background: var(--color-border-warm);
  transform-origin: center;
  animation: site-hold-rule 900ms ease both;
}

.site-hold-headline {
  margin: 0;
  font-family: var(--font-family-display);
  font-size: clamp(1.5rem, 1.1rem + 1.6vw, 2.125rem);
  font-weight: 400;
  line-height: 1.22;
  letter-spacing: var(--tracking-tight);
  text-wrap: balance;
}

.site-hold-message {
  margin: 0;
  max-width: 44ch;
  color: var(--color-body-copy);
  font-size: var(--text-sm);
  line-height: 1.7;
  text-wrap: pretty;
}

.site-hold-action {
  padding: 14px 30px;
  border: 1px solid var(--color-gold);
  color: var(--color-gold);
  font-size: var(--text-xs);
  letter-spacing: var(--tracking-widest);
  text-transform: uppercase;
  text-decoration: none;
  transition: background-color 220ms ease, color 220ms ease;

  &:is(:hover, :focus-visible) {
    background: var(--color-gold);
    color: var(--color-white);
  }
}

@keyframes site-hold-rule {
  from { transform: scaleX(0); }
  to { transform: scaleX(1); }
}

@media (prefers-reduced-motion: reduce) {
  .site-hold-rule { animation: none; }
}
</style>
