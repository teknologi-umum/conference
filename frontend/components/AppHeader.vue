<script lang="ts" setup>

const isMobile = ref(false)
const showNavbar = ref(true)

onMounted(() => {
    // Get screen width
    const screenWidth = window.innerWidth
    isMobile.value = screenWidth < 992

    if(isMobile.value) showNavbar.value = false
    
})
</script>
<template>
    <header id="page-header" class="relative">
        <div class="container | px-5 py-5 mx-auto flex justify-between items-center">
            <div class="header-left">
                <router-link to="/" class="font-bold">TeknumConf</router-link>
            </div>
            <Button color="outline-white" v-show="isMobile" @click="showNavbar = !showNavbar">&equiv;</Button>
            <div class="header-right">
                <nav class="container mx-auto" v-show="showNavbar">
                    <ul class="nav-menu">
                        <li class="nav-item">
                            <router-link class="nav-link" to="#hero">Home</router-link>
                        </li>
                        <li class="nav-item">
                            <router-link class="nav-link" to="#about">About</router-link>
                        </li>
                        <li class="nav-item">
                            <router-link class="nav-link" to="#faq">FAQ</router-link>
                        </li>
                    </ul>
                </nav>
            </div>
        </div>
    </header>
</template>
<style>
.header-right {
    position: absolute;
    width: 100%;
    background-color: rgba(248, 248, 248, 0.121);
    left: 0;
    top: 0;
    right: 0;
    top: 100%;
}
.nav-menu {
}
.nav-link {
    display: block;
    padding: 1rem 1.25rem;
}

@media screen and (min-width: 992px) {
    .header-right {
        position: static;
        background-color: transparent;
    }
    .header-right .nav-menu {
        display: flex;
        justify-content: flex-end;
    }
}
</style>