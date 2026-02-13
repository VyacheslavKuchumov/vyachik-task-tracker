// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: process.env.NUXT_DEVTOOLS === 'true' },
  modules: [
    '@nuxt/eslint',
    '@nuxt/ui',
    '@pinia/nuxt',
    'pinia-plugin-persistedstate'
  ],
  runtimeConfig: {
    backendUrl: process.env.BACKEND_URL || 'http://localhost:8000'
  },
  css: ['~/assets/css/main.css']
})
