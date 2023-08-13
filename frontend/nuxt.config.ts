// https://nuxt.com/docs/api/configuration/nuxt-config
import transformerDirectives from '@unocss/transformer-directives'
import presetUno from '@unocss/preset-uno'


export default defineNuxtConfig({
  googleFonts: {
    families: {
      Rubik: true
    }
  },
  presets: [
    presetUno(),
  ],
  app: {
    head: {
      titleTemplate: (t) => t ? `${t} - TeknumConf` : 'TeknumConf' 
    }
  },
  css: ['assets/css/style.css'],
  modules: [
    '@unocss/nuxt',
    '@nuxtjs/google-fonts'
  ],
  devtools: { enabled: true },
  
})
