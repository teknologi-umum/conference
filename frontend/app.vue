<script setup lang="ts">
useHead({
  titleTemplate: (t) => t ? `${t} - TeknumConf` : 'TeknumConf' 
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