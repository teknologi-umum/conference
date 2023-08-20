<script setup>
const config = useRuntimeConfig()
useHead({
    title: "Home"
})
</script>
<template>
    <div class="page">
        <section id="hero" class="section">
            <div class="container | mx-auto px-5">
                <div class="hero-card | p-12 rounded-2xl blur-bg">
                    <div class="flex between">
                        <div class="card-left">
                            <h1 class="hero-title | mb-8  leading-snug">TeknumConf<br/> 2023</h1>
                            <h3 v-if="config.public.announcementDate" class="text-xl lg:text-3xl mb-5 flex items-center gap-3">
                                <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24"><path fill="currentColor" d="M19 19H5V8h14m-3-7v2H8V1H6v2H5c-1.11 0-2 .89-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V5a2 2 0 0 0-2-2h-1V1"/></svg>
                                <span>{{config.public.announcementDate}}</span>
                            </h3>
                            <h3 class="text-xl lg:text-3xl flex items-center gap-3">
                                <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24"><path fill="currentColor" d="M12 6.5A2.5 2.5 0 0 1 14.5 9a2.5 2.5 0 0 1-2.5 2.5A2.5 2.5 0 0 1 9.5 9A2.5 2.5 0 0 1 12 6.5M12 2a7 7 0 0 1 7 7c0 5.25-7 13-7 13S5 14.25 5 9a7 7 0 0 1 7-7m0 2a5 5 0 0 0-5 5c0 1 0 3 5 9.71C17 12 17 10 17 9a5 5 0 0 0-5-5Z"/></svg>
                                <span>Depok, Indonesia</span>
                            </h3>
                            <Btn class="mt-15" size="xl">Save my spot!</Btn>
                        </div>
                        <div class="card-right flex flex-grow-1 justify-end items-center">
                            <img src="@/assets/images/logo-white.png" alt="TeknumConf logo" class="w-100 hidden lg:block rounded-full">
                        </div>
                    </div>
                </div>
            </div>
        </section>
        <section id="whats-good" class="section my-32">
            <div class="container | mx-auto px-5">
                <div class="flex items-end flex-wrap">
                    <h2 class="section-title | w-full lg:w-1/2 leading-tight">What's Good on TeknumConf?</h2>
                    <ul class="list flex-grow-1 mt-5">
                        <li>Get to attend technology related workshop</li>
                        <li>Also attend non-technology related workshop, refreshen your mind</li>
                        <li>Career advice talk show from wide-ranging backgrounds</li>
                        <li>Hangout and network yourself with others</li>
                        <li v-if="config.public.fifaChampionship">Local FIFA championship, just to have some fun!</li>
                    </ul>
                </div>
            </div>
        </section>
        <section id="faq" class="section">
            <div class="container | mx-auto px-5">
                <div class="section-header text-center mb-16">
                    <h2 class="section-title">Frequently Asked Questions</h2>
                </div>
                <div class="section-content">
                    <Accordion></Accordion>
                </div>
            </div>
        </section>
    </div>
</template>
<style>
.hero-card {
    background-color: rgba(255, 255, 255, 0.076);
    outline: 2px solid rgba(255, 255, 255, 0.245);
    transition: all .3s;
}
.hero-card:hover {
    background-color: rgba(255, 255, 255, 0.115);
    outline: 2px solid white;
}
.hero-title {
    font-size: clamp(3rem, 7vw, 10rem);
    font-weight: 400;
}
.section-title {
    font-size: clamp(3rem, 5vw, 4rem);
}
.list li {
    padding: 2rem 1rem;
    font-size: 1.3rem;
    border-bottom: 1px solid #ffffff5b;
}
</style>