<script setup>
const events = ref([
    {
        name: 'Registration',
        fromHour: '08.00',
        toHour: '10.00'
    },
    {
        name: 'Introduction',
        fromHour: '08.00',
        toHour: '10.00'
    },
    {
        name: 'Talkshow',
        fromHour: '10.00',
        toHour: '14.00'
    },
    {
        name: 'Closing',
        fromHour: '08.00',
        toHour: '14.00'
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