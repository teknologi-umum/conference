<script setup>
const events = ref([
    {
        name: 'Doors open / Opening Remarks',
        fromHour: '13.00',
        toHour: '13.05'
    },
    {
        name: 'The Business Mind in a Pragmatic Engineer: Building Boring Businesses by M. Gilang Januar',
        fromHour: '13.05',
        toHour: '13.45'
    },
    {
        name: 'Loving Rust not because of the Performance by Mustafa Zaki Assagaf',
        fromHour: '13.50',
        toHour: '14.30'
    },
    {
        name: 'Career Talk: How to strive in Indonesia tech industry',
        fromHour: '14.35',
        toHour: '15.35'
    },
    {
        name: 'Break Ashar / Networking Session',
        fromHour: '15.35',
        toHour: '16.00'
    },
    {
        name: 'Surprise Workshop (soon)',
        fromHour: '16.05',
        toHour: '16.45'
    },
    {
        name: 'Melakukan TDD meski situasi tidak kondusif by Riza Ramadan',
        fromHour: '16.50',
        toHour: '17.30'
    },
    {
        name: 'Break Maghrib / Networking Session',
        fromHour: '17.30',
        toHour: '18.15'
    },
    {
        name: 'Tanda-Tanda Kamu Punya Skill Issue by Didiet Noor',
        fromHour: '18.15',
        toHour: '18.45'
    },
    {
        name: 'Closing Remarks',
        fromHour: '18.45',
        toHour: '19.00'
    },
])
</script>
<template>
    <div class="rundown">
        <div class="rundown-event" v-for="event in events">
            <div class="rundown-name mb-2">{{ event.name }}</div>
            <div class="rundown-hour text-sm text-gray-200">{{ event.fromHour }} - {{ event.toHour }}</div>
        </div>
    </div>
</template>
<style>
.rundown-event {
    position: relative;
    padding-left: 3rem;
    min-height: 100px;
}
.rundown-event:not(:last-child):after {
    position: absolute;
    content: '';
    top: .5rem;
    height: 100%;
    z-index: 1;
    background: var(--color-primary-gradient);
    width: 2px;
    left: 4px;
}
.rundown-event:before {
    position: absolute;
    content: '';
    z-index: 2;
    left: 0;
    top: .5rem;
    width: 8px;
    border-radius: 50%;
    box-shadow: 0 3px 20px 2px #cdcde275;
    height: 8px;
    border: 1px solid black;
    background-color: var(--color-primary);
}
</style>