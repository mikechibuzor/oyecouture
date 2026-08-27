<script setup lang="ts">
import type { IconNode } from 'lucide'
import { Camera, Music2 } from 'lucide'
import {
  STOREFRONT_COPY,
  STOREFRONT_FOOTER_GROUPS,
  STOREFRONT_SOCIAL_LINKS,
  isInternalStorefrontRoute,
} from '~/constants/storefront'
import {
  StorefrontRoute,
  StorefrontSocialPlatform,
} from '~/types/storefront'

const { appointmentWhatsAppUrl } = await useAppointmentWhatsApp()

const resolveFooterLink = (storefrontRoute: string): string =>
  storefrontRoute === StorefrontRoute.BookFitting
    ? appointmentWhatsAppUrl.value
    : storefrontRoute

const socialIconByPlatform: Readonly<
  Partial<Record<StorefrontSocialPlatform, IconNode>>
> = {
  [StorefrontSocialPlatform.Instagram]: Camera,
  [StorefrontSocialPlatform.TikTok]: Music2,
}

const footerSocialLinks = STOREFRONT_SOCIAL_LINKS.map((socialLink) => ({
  ...socialLink,
  icon: socialIconByPlatform[socialLink.platform],
}))
</script>

<template>
  <footer class="storefront-footer">
    <div class="storefront-footer-content">
      <section class="footer-brand-column">
        <h2>{{ STOREFRONT_COPY.footer.heading }}</h2>
        <StorefrontBrandLogo />
        <p>{{ STOREFRONT_COPY.footer.description }}</p>

        <ul
          class="footer-social-list"
          :aria-label="STOREFRONT_COPY.footer.socialNavigationLabel"
        >
          <li v-for="socialLink in footerSocialLinks" :key="socialLink.platform">
            <a
              class="footer-social-icon"
              :href="socialLink.href"
              target="_blank"
              rel="noopener noreferrer"
              :aria-label="socialLink.label"
            >
              <StorefrontIcon
                v-if="socialLink.icon"
                :icon-node="socialLink.icon"
                :size="17"
                :stroke-width="1.6"
                aria-hidden="true"
              />
              <span v-else class="footer-social-letter" aria-hidden="true">X</span>
            </a>
          </li>
        </ul>
      </section>

      <section
        v-for="footerGroup in STOREFRONT_FOOTER_GROUPS"
        :key="footerGroup.heading"
        class="footer-link-column"
      >
        <h2>{{ footerGroup.heading }}</h2>
        <ul>
          <li v-for="footerLink in footerGroup.links" :key="footerLink.to">
            <NuxtLink
              :to="resolveFooterLink(footerLink.to)"
              :external="footerLink.to === StorefrontRoute.BookFitting || !isInternalStorefrontRoute(footerLink.to)"
            >
              {{ footerLink.label }}
            </NuxtLink>
          </li>
        </ul>
      </section>
    </div>

    <div class="storefront-footer-bottom">
      <p>{{ STOREFRONT_COPY.footer.copyright }}</p>
    </div>
  </footer>
</template>
