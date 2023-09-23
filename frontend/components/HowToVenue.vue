<script setup lang="ts">

const showWheelchair = ref(false)
</script>
<template>
    <Tabs class="tab-venue mb-5">
        <Tab name="Train (preferred)">
            <p>Jakarta has a public train called KRL. Using this method is most preferred during the day of the event.
                Please note that this walk path is not wheelchair accessible.</p>
            <ol class="pl-5">
                <li>From any other station, go for <strong>"Bogor Line"</strong></li>
                <li>Stop at <strong>Pondok Cina station</strong>, don't use the normal exit, go upstairs (through stairs or
                    escalators) exit through gate A pointed to "Perumnas Mahata" and "Jl. Margonda Raya"</li>
                <li>Go down from Perumnas Mahata through elevator (from UG to GF) or stairs, walk to Depok Town Square (you
                    should be able to see its' parking building)</li>
                <li>Go up the stairs (black alumunium stairs), then go left at the interchange</li>
                <li>The venue is on your exact left</li>
            </ol>
        </Tab>
        <Tab name="Ridesharing (Gojek, Grab)">
            <ol class="pl-5">
                <li>Set your destination either directly to Depok Town Square or Margo City. If Margo City is cheaper for
                    you, you can walk through the pedestrian bridge that connects Margo City to Depok Town Square.</li>
                <li>From the lobby, go up 1 floor to the UG level</li>
                <li>Walk through the right side of the mall. Halfway, you'll see an exit on the right.</li>
                <li>Go out through the door (you'll be outdoor now), the venue is across the mall entrance door slightly to
                    the left</li>
            </ol>
        </Tab>
        <Tab name="Personal Vehicle (Car, Bikes)">
            <p>Parking is available for both cars and motorbikes. For cars, it's recommended to park near P3 level.</p>
            <p>If you are a person with wheelchair, you can drop exactly at the front of the venue. It is across the mall
                entrance door after you enter the cars parking gate.</p>
        </Tab>
    </Tabs>
    <div class="wheelchair-accessibility">
        <a href="#" aria-controls="wheelchair-card" class="flex items-center gap-1 mb-5"
            @click.prevent="showWheelchair = !showWheelchair">Wheelchair accessibility <svg xmlns="http://www.w3.org/2000/svg"
                width="20" height="20" viewBox="0 0 24 24">
                <g fill="none">
                    <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="1.5" />
                    <path stroke="currentColor" stroke-linecap="round" stroke-width="1.5"
                        d="M10.125 8.875a1.875 1.875 0 1 1 2.828 1.615c-.475.281-.953.708-.953 1.26V13" />
                    <circle cx="12" cy="16" r="1" fill="currentColor" />
                </g>
            </svg></a>
    <div class="card card-translucent mb-3" v-show="showWheelchair" id="wheelchair-card">
        <h3 class="mb-3">Wheelchair Accessibility</h3>
        <p>
            Kalau diantar, minta diturunin persis di depan venue. dari pintu masuk parkir (ambil tiket), naik ke atas
            ramp 1 kali, abis itu di ramp ke 2, ambil kanan (jangan naik ramp kedua), lokasi ada di sebelah kanan persis
        </p>
    </div>
</div></template>
<style>.tab-venue li {
    margin-bottom: 1rem;
}</style>