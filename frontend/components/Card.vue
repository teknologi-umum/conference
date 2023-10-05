<script lang="ts" setup>

defineProps({
    image_url: String,
    title: String,
    description: String
})

</script>
<template>
    <div class="card rounded-2xl bg-white bg-opacity-10 backdrop-blur-lg">
        <div class="card-image h-100" v-if="image_url">
            <img :src="image_url" :alt="`Photo of ${title}`" class="w-full h-full object-cover rounded-t-2xl mb-0">
        </div>
        <div class="card-body  rounded-b-2xl p-5">
            <h3 class="mb-4 text-3xl text-white">{{ title }}</h3>
            <slot></slot>
            <p class="text-lg">{{ description }}</p>
            
        </div>
    </div>
</template>
<style>
.card-image {
    position: relative
}
.card-image img {
}
.card-image:after {
    content: '';
    z-index: 3;
    position: absolute;
    inset: 0;
}
</style>