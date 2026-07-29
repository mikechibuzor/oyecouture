import process from 'node:process'

import { DEFAULT_API_BASE_URL } from './constants/config'

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  css: ['~/assets/scss/main.scss'],
  runtimeConfig: {
    public: {
      apiBase: process.env.API_BASE_URL || DEFAULT_API_BASE_URL,
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
