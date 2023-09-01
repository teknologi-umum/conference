<script setup lang="ts">
const faq = [
    {
        question: 'What is Teknum?',
        answer: `Teknologi Umum (Common Technology in English) — familiar with the abbreviation Teknum, cause it's shorter — is a community on Telegram that revolves around technology related things, with programming as its' main focus. It was founded as a means of frustration and anxiety of its owner that's having a hard time finding a community to talk things considered as out-of-topic on general programming groups, now it's evolved to a community that create projects and contributed to Open Source.`
    },
    {
        question: 'Do Teknologi Umum have Open Source projects?',
        answer: `Yes! We have many of them. As well as scientific journal that's being made public and widespread for pushing mankind's knowledge, Teknologi Umum believes the same thing for software engineering and the like. You can find many of our public project and repositories at <a target='_blank' href='https://github.com/teknologi-umum'>https://github.com/teknologi-umum</a>`
    },
    {
        question: 'What is TeknumConf?',
        answer: `TeknumConf is a meetup event made by Teknologi Umum community, containing workshops, talkshows, and networking session. Unlike other programming or technology events, TeknumConf includes topic that's not related to technology.`
    },
    {
        question: 'When will TeknumConf be held?',
        answer: `We will be in Jakarta on <span class='text-orange-300 font-bold'>21 October 2023</span>, at the same time with Hacktoberfest.`
    },
    {
        question: 'Will there be a live streaming or recorded videos?',
        answer: `No, there won't be any live streaming for TeknumConf. However, we are planning to record every workshop and talkshow, to be uploaded later to YouTube.`
    },
    {
        question: 'Why is TeknumConf is a paid event?',
        answer: `Teknologi Umum is a very young community. We don't have the funds needed to move independently without community donations. All of the ticket sales will be allocated to pay for venue rent, merchandise and other things that will be documented and shown publicly.`
    },
]
</script>
<template>
    <div class="accordion rounded-md">
        <div class="accordion-item" v-for="(question, i) in faq">
            <input type="radio" name="radio" class="accordion-radio hidden" :id="`question-${i}`">
            <label :for="`question-${i}`" class="accordion-title cursor-pointer">{{ question.question }}</label>
            <p class="accordion-content leading-7" v-html="question.answer"></p>
        </div>
    </div>
</template>
<style>
.accordion-item, .accordion-title {
    display: block;
    transition: all .2s;
    background-color: rgba(141, 141, 141, 0.088);
}
.accordion-title, .accordion-content {
    padding: 1rem;
}
.accordion-item:not(:last-child) {
    border-bottom: 2px solid #3a3a3a;
}
.accordion-content {
    display: none;
}
.accordion-radio:checked ~ .accordion-content {
    display: block;
}
.accordion-radio:checked ~ .accordion-title {
    background-color: rgba(141, 141, 141, 0.251);
}
.accordion-item:first-child {
    border-top-left-radius: .5rem;
    border-top-right-radius: .5rem;
}
.accordion-item:last-child {
    border-bottom-left-radius: .5rem;
    border-bottom-right-radius: .5rem;
}
</style>