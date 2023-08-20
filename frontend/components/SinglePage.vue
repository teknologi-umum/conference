<script lang="ts" setup>
const props = defineProps({
    title: String,
})
</script>

<template>
    <section class="single">
        <div class="container-sm mx-auto px-5 py-10">
            <div class="single-title mb-10">
                <h1 class="text-5xl">{{ title }}</h1>
            </div>
            <div class="single-content">
                <slot></slot>
            </div>
        </div>
    </section>
</template>