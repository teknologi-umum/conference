// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  runtimeConfig: {
    public: {
      showAnnouncementDate: true,
      announcementDate: "21 October 2023",
      attendeeRegistration: true,
      speakerRegistration: true,
      eventSchedule: true,
      aggressiveIntroduction: true,
      fifaChampionship: true,
      backendBaseUrl: process.env.NODE_ENV === "development" ? "http://localhost:8080" : "https://conf.teknologiumum.com/api"
    }
  },
  googleFonts: {
    families: {
      Rubik: true
    }
  },
  app: {
    head: {
    }
  },
  css: ['assets/css/style.css'],
  modules: [
    '@unocss/nuxt',
    '@nuxtjs/google-fonts',
    // "@model-w/sentry"
  ],
  // sentry: {
  //   dsn: process.env.SENTRY_DSN,
  //   environment: process.env.SENTRY_ENVIRONMENT
  // },
  devtools: { enabled: true },
})
