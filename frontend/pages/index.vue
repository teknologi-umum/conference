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
                            <h1 class="hero-title | mb-8  leading-snug">TeknumConf<br/> 2024</h1>
                            <h3 class="text-xl lg:text-3xl mb-5 flex items-center gap-3">
                              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24"><path fill="currentColor" d="M19 19H5V8h14m-3-7v2H8V1H6v2H5c-1.11 0-2 .89-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V5a2 2 0 0 0-2-2h-1V1"/></svg>
                              <span>COMING SOON</span>
                            </h3>
                            <h3 class="text-xl flex items-center gap-3">
                                <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                  <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
                                  <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
                                </svg>
                                <span>We are still preparing for TeknumConf 2024</span>
                            </h3>
                            <Btn v-umami="{name: 'go-to-2023-page-btn'}" :isLink="true" to="/2023" class="mt-5 lg:mt-15 text-center btn-save-my-spot" size="xl">Discover What You Missed in 2023!</Btn>
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
                    <ul class="list list-none flex-grow-1 mt-5">
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


@media screen and (max-width: 768px) {
  .card-left h1 {
    font-size: 32px;
  }
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
@media screen and (min-width: 1024px) {
    .rundown-title {
        writing-mode: vertical-lr;
        text-orientation:mixed;
        font-size: clamp(6rem, 10vw, 9rem);
        text-align: right;
        transform: rotateZ(180deg);
    }
    #rundown .section-header {
        justify-content: flex-end;
        display: flex;
        padding: 0 3rem;
    }
}
#venue {
    background: linear-gradient(120deg, rgba(1, 2, 15, 0.883),
                    rgba(0, 0, 0, 0.95), 
                    rgba(5, 12, 34, 0.879)),
                url('@/assets/images/venue.jpg') no-repeat;
    background-size: cover;
    padding: 3rem 0;
    backdrop-filter: blur(1rem);
    min-height: 500px;
    line-height: 2rem;
}
#img-venue {
    object-fit: cover;
    position: absolute;
    top: 0;
    right: 0;
    width: 40%;
    height: 100%;
}
</style>