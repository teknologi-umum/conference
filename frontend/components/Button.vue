<script setup lang="ts">
const props = defineProps({
    isLink: {
        type: Boolean,
        default: false 
    },
    to: {
        type: String
    },
    color: {
        type: String,
        default: 'primary'
    },
    size: {
        type: String,
        default: 'md'
    }
})

const classes = computed(() => `btn btn-${props.color} btn-${props.size}`)
</script>
<template>
    <router-link :to="to!" :class="classes" v-if="isLink">
        <slot></slot>
    </router-link>
    <button :class="classes" v-else>
        <slot></slot>
    </button>
</template>
<style>
.btn-sm {
    --radius: .2rem;
}
.btn-md {
    --radius: .3rem;
}
.btn-lg {
    --radius: .75rem
}
.btn {
    padding: .5rem;
    cursor: pointer;
    transition: all .3s;
    border-radius: var(--radius);
}
.btn-md {
    padding: .75rem 1rem;
}
.btn-outline-white {
    border: 1px solid white;
    background-color: transparent;
    color: white;
}
.btn-outline-white:hover {
    background-color: white;
    color: black;
}
</style>