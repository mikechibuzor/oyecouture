<script setup lang="ts">
import {
  STOREFRONT_COPY,
  STOREFRONT_FOOTER_GROUPS,
  isInternalStorefrontRoute,
} from '~/constants/storefront'
import { StorefrontRoute } from '~/types/storefront'

const { appointmentWhatsAppUrl } = await useAppointmentWhatsApp()

const resolveFooterLink = (storefrontRoute: string): string =>
  storefrontRoute === StorefrontRoute.BookFitting
    ? appointmentWhatsAppUrl.value
    : storefrontRoute
</script>

<template>
  <footer class="storefront-footer">
    <div class="storefront-footer-content">
      <section class="footer-brand-column">
        <h2>{{ STOREFRONT_COPY.footer.heading }}</h2>
        <StorefrontBrandLogo />
        <p>{{ STOREFRONT_COPY.footer.description }}</p>

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
