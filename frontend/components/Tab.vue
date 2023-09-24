<script setup>
const props = defineProps({
    name: String
})
</script>
<template>
    <div class="tab p-5 text-left" :id="name">
        <slot></slot>
    </div>
</template>