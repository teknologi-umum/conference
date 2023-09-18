<script lang="ts" setup>

const isMobile = ref(false)
const showNavbar = ref(true)
function checkWindowSize() {
    const screenWidth = window.innerWidth
    isMobile.value = screenWidth < 1024
    showNavbar.value = !isMobile.value
}
onMounted(() => {
    console.log("rstrst")
    // Get screen width
    checkWindowSize()
    window.addEventListener('resize', checkWindowSize)    
})
</script>
<template>
    <header id="page-header" class="relative">
        <div class="container | px-5 py-5 mx-auto flex justify-between items-center">
            <div class="header-left">
                <router-link to="/" class="font-bold">
                    <img src="@/assets/images/logo_copy.png" alt="TeknumConf" height="60">
                </router-link>
            </div>
            <Btn color="outline-white" v-show="isMobile" @click="showNavbar = !showNavbar">&equiv;</Btn>
            <div class="header-right z-10">
                <nav class="container mx-auto" v-show="showNavbar">
                    <ul class="nav-menu list-none">
                        <li class="nav-item">
                            <router-link class="nav-link white" to="/#hero">Home</router-link>
                        </li>
                        <li class="nav-item">
                            <router-link class="nav-link white" to="/#whats-good">About</router-link>
                        </li>
                        <li class="nav-item">
                            <router-link class="nav-link white" to="/#faq">FAQ</router-link>
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
    backdrop-filter: blur(100px);
}
.nav-link {
    display: block;
    padding: 1rem 1.25rem;
}

@media screen and (min-width: 1024px) {
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