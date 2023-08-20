// https://nuxt.com/docs/api/configuration/nuxt-config
import transformerDirectives from '@unocss/transformer-directives'
import presetUno from '@unocss/preset-uno'


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
    '@nuxtjs/google-fonts'
  ],
  devtools: { enabled: true },
  
})
