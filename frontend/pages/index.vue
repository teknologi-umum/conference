<script setup>
import personImage from "@/assets/images/person.jpg"
const config = useRuntimeConfig()
useHead({
    title: "Home"
})

const speakers = ref([
    {
        name: "Didiet Noor",
        title: "Tanda-Tanda Kamu Punya Skill Issue",
        description: "Developer bisa coding doang itu banyak. Banyak juga yang dilayoff. Saingan kamu banyak kalo cuma bisa coding.",
    },
    {
        name: "M. Gilang Januar",
        title: "The Business Mind in a Pragmatic Engineer: Building Boring Businesses",
        description: "Berbicara tentang memiliki side hustle sebagai seorang software engineer. Mulai dari ideation, business and marketing strategy, contoh studi kasus, dan hal menarik lainnya. Harapannya peserta dapat melihat business value dari kacamata seoarang pragmatic engineer.",
    },
    {
        name: "Riza Ramadan",
        title: "Melakukan TDD meski situasi tidak kondusif",
        description: "Often, 'You're doing TDD wrong' is heard, but usually, TDD is just applied to unsuitable problems. Many do it right, just not for the right reasons. In this talk, we will discuss when it makes sense to use TDD and when it doesn't, regardless of our understanding of how to do TDD."
    },
    {
        name: "Mustafa Zaki A.",
        title: "Loving Rust not because of the performance",
        description: "Rust dikenal dengan bahasa pemrograman yang cepat. Namun sebenarnya ada sisi lain yang ditawarkan di Rust. Talk kali ini akan membahas apa saja hal yang ditawarkan oleh Rust selain kecepatan.",
    }
]) 
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
                            <Btn :isLink="true" to="register" class="mt-5 lg:mt-15 text-center btn-save-my-spot" size="xl">Save my spot!</Btn>
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
        <section id="speakers" class="section my-32">
            <div class="container | mx-auto px-5">
                <h2 class="section-title | mb-24 text-center w-full leading-tight">Speakers</h2>
                <div class="section-body">
                    <div class="speakers | grid grid-cols-1 lg:grid-cols-4 gap-10">
                        <Card v-for="speaker in speakers" :key="speaker.name" :image_url="personImage" :title="speaker.name" :description="speaker.description"></Card>
                    </div>
                </div>
            </div>
        </section>
        <section id="rundown" class="section">
            <div class="container | mx-auto px-5 lg:flex">
                <div class="section-header text-center mb-16 lg:w-1/2">
                    <h2 class="section-title rundown-title">Rundown</h2>
                </div>
                <div class="section-content flex justify-center">
                    <Rundown></Rundown>
                </div>
            </div>
        </section>
        <section id="venue" class="section relative">
            <div class="container | mx-auto px-5 lg:flex ">
                <div class="section-left text-center mb-16 lg:w-1/2">
                    <h2 class="section-title mb-8">How to get to venue</h2>
                    <Tabs>
                        <Tab name="Train (preferred)">
                          <p>Jakarta has a public train called KRL. Using this method is most preferred during the day of the event. Please note that this walk path is not wheelchair accessible.</p>
                          <ol>
                            <li>From any other station, go for <strong>"Bogor Line"</strong></li>
                            <li>Stop at <strong>Pondok Cina station</strong>, don't use the normal exit, go upstairs (through stairs or escalators) exit through gate A pointed to "Perumnas Mahata" and "Jl. Margonda Raya"</li>
                            <li>Go down from Perumnas Mahata through elevator (from UG to GF) or stairs, walk to Depok Town Square (you should be able to see its' parking building)</li>
                            <li>Go up the stairs (black alumunium stairs), then go left at the interchange</li>
                            <li>The venue is on your exact left</li>
                          </ol>
                        </Tab>
                        <Tab name="Ridesharing (Gojek, Grab, Maxim)">
                          <ol>
                            <li>Set your destination either directly to Depok Town Square or Margo City. If Margo City is cheaper for you, you can walk through the pedestrian bridge that connects Margo City to Depok Town Square.</li>
                            <li>From the lobby, go up 1 floor to the UG level</li>
                            <li>Walk through the right side of the mall. Halfway, you'll see an exit on the right.</li>
                            <li>Go out through the door (you'll be outdoor now), the venue is across the mall entrance door slightly to the left</li>
                          </ol>
                        </Tab>
                        <Tab name="Personal Vehicle (Car, Bikes)">
                          <p>Parking is available for both cars and motorbikes. For cars, it's recommended to park near P3 level.</p>
                          <p>If you are a person with wheelchair, you can drop exactly at the front of the venue. It is across the mall entrance door after you enter the cars parking gate.</p>
                        </Tab>
                    </Tabs>
                </div>
                <div class="section-right flex justify-center">
                    <img src="/detos.jpg" id="img-venue">
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
    background: linear-gradient(120deg, rgba(255, 255, 255, 0.051), rgba(255, 255, 255, 0.131), rgba(93, 93, 93, 0.24), rgba(255, 255, 255, 0.2));
    padding: 3rem 0;
    backdrop-filter: blur(1rem);
    min-height: 500px;
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