<script setup lang="ts">
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

// The client's own marks, replacing the lucide stand-ins: Instagram was a
// generic camera and TikTok a music note, neither of which read as the brand
// they were standing for, and X had no icon at all.
const socialIconByPlatform: Readonly<
  Partial<Record<StorefrontSocialPlatform, string>>
> = {
  [StorefrontSocialPlatform.Instagram]: '/images/social/instagram-line.png',
  [StorefrontSocialPlatform.TikTok]: '/images/social/tiktok-fill.png',
  [StorefrontSocialPlatform.X]: '/images/social/twitter-x-line.png',
}

const footerSocialLinks = STOREFRONT_SOCIAL_LINKS.map((socialLink) => ({
  ...socialLink,
  iconSrc: socialIconByPlatform[socialLink.platform],
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
              <img
                v-if="socialLink.iconSrc"
                class="footer-social-mark"
                :src="socialLink.iconSrc"
                :alt="''"
                width="17"
                height="17"
                loading="lazy"
              >
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
