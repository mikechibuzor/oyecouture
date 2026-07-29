<script setup lang="ts">
import type { IconNode } from 'lucide'
import { AtSign, Camera, MessageCircle, Music2 } from 'lucide'
import {
  STOREFRONT_COPY,
  STOREFRONT_FOOTER_GROUPS,
  STOREFRONT_LEGAL_LINKS,
  STOREFRONT_SOCIAL_LINKS,
} from '~/constants/storefront'
import { StorefrontSocialPlatform } from '~/types/storefront'

const SOCIAL_ICONS: Readonly<Record<StorefrontSocialPlatform, IconNode>> = {
  [StorefrontSocialPlatform.Instagram]: Camera,
  [StorefrontSocialPlatform.X]: AtSign,
  [StorefrontSocialPlatform.TikTok]: Music2,
  [StorefrontSocialPlatform.WhatsApp]: MessageCircle,
}
</script>

<template>
  <footer class="storefront-footer">
    <div class="storefront-footer-content">
      <section class="footer-brand-column">
        <h2>{{ STOREFRONT_COPY.footer.heading }}</h2>
        <StorefrontBrandLogo />
        <p>{{ STOREFRONT_COPY.footer.description }}</p>

        <div class="footer-social-list">
          <span
            v-for="socialLink in STOREFRONT_SOCIAL_LINKS"
            :key="socialLink.platform"
            class="footer-social-icon"
            role="img"
            :aria-label="socialLink.label"
            :title="socialLink.label"
          >
            <StorefrontIcon
              :icon-node="SOCIAL_ICONS[socialLink.platform]"
              :size="18"
              :stroke-width="1.5"
              aria-hidden="true"
            />
          </span>
        </div>
      </section>

      <section
        v-for="footerGroup in STOREFRONT_FOOTER_GROUPS"
        :key="footerGroup.heading"
        class="footer-link-column"
      >
        <h2>{{ footerGroup.heading }}</h2>
        <ul>
          <li v-for="footerLink in footerGroup.links" :key="footerLink.to">
            <NuxtLink :to="footerLink.to" external>
              {{ footerLink.label }}
            </NuxtLink>
          </li>
        </ul>
      </section>
    </div>

    <div class="storefront-footer-bottom">
      <p>{{ STOREFRONT_COPY.footer.copyright }}</p>
      <nav :aria-label="STOREFRONT_COPY.footer.legalNavigationLabel">
        <NuxtLink
          v-for="legalLink in STOREFRONT_LEGAL_LINKS"
          :key="legalLink.to"
          :to="legalLink.to"
          external
        >
          {{ legalLink.label }}
        </NuxtLink>
      </nav>
    </div>
  </footer>
</template>
