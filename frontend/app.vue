<script setup lang="ts">
useHead({
  titleTemplate: (t) => t ? `${t} - TeknumConf` : 'TeknumConf',
  script:[
    {
      src: 'https://analytics.teknologiumum.com/script.js',
      'data-website-id': 'b6a8d7a7-a90a-476c-8fd2-355afea623c4'
    },
  ],
  })
</script>
<template>
  <div>    
    <NuxtLayout name="app">
      <NuxtPage/>
    </NuxtLayout>
  </div>
</template>
<style>
</style>