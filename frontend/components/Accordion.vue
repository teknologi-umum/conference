<script setup lang="ts">
const faq = [
    {
        question: 'Lorem ipsum dolor, sit amet consectetur adipisicing elit?',
        answer: 'Minima obcaecati tempore quis labore, nam repellat                 blanditiis? Alias, tempore, error est esse aut ratione sequi, quos eum deserunt aperiam nam provident.'
    },
    {
        question: 'Lorem ipsum dolor, sit amet consectetur adipisicing elit?',
        answer: 'Minima obcaecati tempore quis labore, nam repellat                 blanditiis? Alias, tempore, error est esse aut ratione sequi, quos eum deserunt aperiam nam provident.'
    },
    {
        question: 'Lorem ipsum dolor, sit amet consectetur adipisicing elit?',
        answer: 'Minima obcaecati tempore quis labore, nam repellat                 blanditiis? Alias, tempore, error est esse aut ratione sequi, quos eum deserunt aperiam nam provident.'
    },
]
</script>
<template>
    <div class="accordion rounded-md">
        <div class="accordion-item" v-for="(question, i) in faq">
            <input type="radio" name="radio" class="accordion-radio hidden" :id="`question-${i}`">
            <label :for="`question-${i}`" class="accordion-title cursor-pointer">{{ question.question }}</label>
            <p class="accordion-content leading-7">{{ question.answer }}</p>
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