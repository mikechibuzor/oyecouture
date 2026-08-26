import process from 'node:process'

import { DEFAULT_API_BASE_URL } from './constants/config'

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false },
  experimental: {
    appManifest: false,
    viewTransition: true,
  },
  css: ['~/assets/scss/main.scss'],
  runtimeConfig: {
    // Server-only: SSR reaches the API container directly inside the network.
    // Overridden in production by NUXT_API_INTERNAL_BASE.
    apiInternalBase: process.env.API_INTERNAL_BASE || DEFAULT_API_BASE_URL,
    public: {
      // Browser-facing: overridden in production by NUXT_PUBLIC_API_BASE.
      apiBase: process.env.API_BASE_URL || '',
      // Bespoke orders hand over to this number; set WHATSAPP_NUMBER per environment.
      whatsappNumber: process.env.WHATSAPP_NUMBER || '',
    },
  },
  // The ghost admin is private: it renders in the browser only, so no admin
  // markup is ever produced server-side, and it is never indexed.
  routeRules: {
    '/admin/**': {
      ssr: false,
      headers: { 'X-Robots-Tag': 'noindex, nofollow' },
    },
  },
  typescript: {
    strict: true,
    typeCheck: true,
  },
  vite: {
    css: {
      preprocessorOptions: {
        scss: {
          additionalData:
            '@use "~/assets/scss/design-system/breakpoints" as *;',
        },
      },
    },
  },
})
