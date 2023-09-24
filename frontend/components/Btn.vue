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
    },
    withShadow: {
        type: Boolean,
        default: false
    }
})

const classes = computed(() => `btn btn-${props.color} btn-${props.size} inline-block`)
</script>
<template>
    <router-link :to="to!" :class="classes" v-if="isLink">
        <slot></slot>
    </router-link>
    <button :class="classes" v-else @click="$emit('click')">
        <slot></slot>
    </button>
</template>
<style>
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
.btn-sm {
    --radius: .2rem;
}
.btn-md {
    --radius: .3rem;
}
.btn-lg {
    --radius: .5rem;
    padding: 1rem;
    font-size: 1.2rem;
    min-width: min(100%, 300px);
}
.btn-xl {
    --radius: .75rem;
    padding: 2rem;
    font-size: 1.4rem;
    min-width: min(100%, 400px);
}
.btn-primary {
    background: var(--color-primary-gradient);
    color: white;
    border: none;
    box-shadow: 0 5px 30px #7474ff56;
}
.btn-primary:hover {
    box-shadow: 0 5px 30px rgba(96, 96, 255, 0.639);
    background: var(--color-primary-gradient-hover);
}
</style>