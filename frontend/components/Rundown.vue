<script setup>
const events = ref([
    {
        name: 'Doors open / Opening Remarks',
        fromHour: '13.00',
        toHour: '13.05',
        recorded: false,
        youtube: ""
    },
    {
        name: 'The Business Mind in a Pragmatic Engineer: Building Boring Businesses by M. Gilang Januar',
        fromHour: '13.05',
        toHour: '13.45',
        recorded: true,
        youtube: "https://youtu.be/k_baN_tPcts?si=-6BteKYlA5FmodFR"
    },
    {
        name: 'Loving Rust not because of the Performance by Mustafa Zaki Assagaf',
        fromHour: '13.50',
        toHour: '14.30',
        recorded: true,
        youtube: "https://youtu.be/_3pr4xzjX8Q?si=Lc0Gz-YOC5zRa_dd"
    },
    {
        name: 'Career Talk: How to strive in Indonesia tech industry',
        fromHour: '14.35',
        toHour: '15.35',
        recorded: true,
        youtube: "https://youtu.be/dhB6dmDcir0?si=rla3sTicDfI2wwkx"
    },
    {
        name: 'Break Ashar / Networking Session',
        fromHour: '15.35',
        toHour: '16.00',
        recorded: false,
        youtube: ""
    },
    {
        name: 'Minimalism on Product Development: An Algobash Story by Elfino Sitompul',
        fromHour: '16.05',
        toHour: '16.45',
        recorded: false,
        youtube: "https://youtu.be/0TUUhXC6OrE?si=5eHbaGcw9fJk9xUv"
    },
    {
        name: 'Melakukan TDD meski situasi tidak kondusif by Riza Ramadan',
        fromHour: '16.50',
        toHour: '17.30',
        recorded: true,
        youtube: "https://youtu.be/A3W-kOHFhuo?si=L7FDaFh32TZMro1s"
    },
    {
        name: 'Break Maghrib / Networking Session',
        fromHour: '17.30',
        toHour: '18.15',
        recorded: false,
        youtube: ""
    },
    {
        name: 'Teks Adalah Mitos by Didiet Noor',
        fromHour: '18.15',
        toHour: '18.45',
        recorded: true,
        youtube: "https://youtu.be/0TUUhXC6OrE?si=5eHbaGcw9fJk9xUv"
    },
    {
        name: 'Closing Remarks',
        fromHour: '18.45',
        toHour: '19.00',
        recorded: false,
        youtube: ""
    },
])
</script>
<template>
    <div class="rundown">
        <div class="rundown-event" v-for="event in events">
            <div class="rundown-name mb-2">{{ event.name }}</div>
            <div class="rundown-hour text-sm text-gray-200">
              {{ event.fromHour }} - {{ event.toHour }}

              <span v-if="event.recorded">
                |
                <a :href="event.youtube" target="_blank">Watch Here!</a>
              </span>
            </div>
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