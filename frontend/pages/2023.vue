<script setup>
import didietImage from "@/assets/images/speaker_didiet.png"
import elfinoImage from "@/assets/images/speaker_elfino.png";
import rizaImage from "@/assets/images/speaker_riza.png";
import gilangImage from "@/assets/images/speaker_gilang.png";
import mustafaImage from "@/assets/images/speaker_mustafa.png";

const config = useRuntimeConfig()
useHead({
  title: "TeknumConf 2023",
})

const speakers = ref([
  {
    name: "Didiet Noor",
    title: "Teks adalah Mitos",
    description: "Semua yang teman-teman ketahui tentang 'teks' hanyalah sebuah mitos. Akan berbicara dan membongkar khurafat mengenai teks, ascii, unicode, emoji, dan glif.",
    image: didietImage,
  },
  {
    name: "M. Gilang Januar",
    title: "The Business Mind in a Pragmatic Engineer: Building Boring Businesses",
    description: "Berbicara tentang memiliki side hustle sebagai seorang software engineer. Mulai dari ideation, business and marketing strategy, contoh studi kasus, dan hal menarik lainnya. Harapannya peserta dapat melihat business value dari kacamata seorang pragmatic engineer.",
    image: gilangImage,
  },
  {
    name: "Riza Ramadan",
    title: "Melakukan TDD meski situasi tidak kondusif",
    description: "Often, 'You're doing TDD wrong' is heard, but usually, TDD is just applied to unsuitable problems. Many do it right, just not for the right reasons. In this talk, we will discuss when it makes sense to use TDD and when it doesn't, regardless of our understanding of how to do TDD.",
    image: rizaImage,
  },
  {
    name: "Mustafa Zaki A.",
    title: "Loving Rust not because of the performance",
    description: "Rust dikenal dengan bahasa pemrograman yang cepat. Namun sebenarnya ada sisi lain yang ditawarkan di Rust. Talk kali ini akan membahas apa saja hal yang ditawarkan oleh Rust selain kecepatan.",
    image: mustafaImage,
  },
  {
    name: "Elfino Sitompul",
    title: "Minimalism on Product Development: An Algobash Story",
    description: "Menjadi developer yang produktif, cepat, efisien, dan dicintai oleh team. Menghilangkan hal-hal yang memperlambat proses product development, apalagi di Algobash.",
    image: elfinoImage,
  }
])

const showVideo = ref(false)
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
                <span>Kode Creative Hub, Depok Town Square<br/>Depok, Indonesia</span>
              </h3>
              <Btn v-umami="{name: 'go-to-2024-page-btn'}" :is-link="true" to="/" class="mt-5 lg:mt-15 text-center btn-save-my-spot" size="xl">Back to Current Event</Btn>
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
          <div class="speakers | grid grid-cols-1 lg:grid-cols-3 gap-10">
            <Card v-for="speaker in speakers" :key="speaker.name" :image_url="speaker.image" :title="speaker.name" :description="speaker.description"></Card>
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
      <div class="container | mx-auto px-5 lg:flex  lg:w-1/2">
        <div class="section-left text-center mb-16">
          <h2 class="section-title leading-normal mb-8">How to get to venue</h2>

          <!-- Toggle video -->
          <Btn color="outline" size="lg2" class="mb-10" @click="showVideo=!showVideo">
            <div class="flex items-center justify-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24"><path fill="currentColor" d="M10 16.5v-9l6 4.5M12 2A10 10 0 0 0 2 12a10 10 0 0 0 10 10a10 10 0 0 0 10-10A10 10 0 0 0 12 2Z"/></svg>
              Show video
            </div>
          </Btn>
          <div v-show="showVideo">
            <video src="@/assets/how_to_venue.mp4" controls class="h-80" ></video>
          </div>

          <HowToVenue/>
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