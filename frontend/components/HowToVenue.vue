<script setup lang="ts">

const showWheelchair = ref(false)
</script>
<template>
    <Tabs class="tab-venue mb-5">
        <Tab name="Go-Jek/Ojol">
            <ol class="pl-5">
                <li>
                    Langsung aja ke Depok Town Square, tapi bandingin harganya dengan ke Margo City. kalau harganya lumayan beda jauh, turun di Margo City aja,
                </li>
                <li>
                    Kalau dari Depok Town Square:
                    <ol class="pl-5" type="a">
                        <li>kalo turun di lobby, masuk aja, naik ke UG, nanti ambil ke arah kanan, ada keluaran terang seakan-akan keluar detos, masuk kesana.</li>
                        <li>kalo mau parkir, kalo bawa motor, parkiran dibawah. kalo bawa mobil, disarankan parkir di UG (P2)</li>
                    </ol>
                </li>
                <li>
                    Kalau dari Margo City: nyebrang lewat Jembatan Penyebrangan Orang
                </li>
            </ol>
        </Tab>
        <Tab name="KRL">
            <ol class="pl-5">
                <li>
                    Kalo kamu dari tengah kota, ambil rute "Bogor Line"
                </li>
                <li>
                    Berhenti di Stasiun Pondok Cina
                </li>
                <li>
                    Dari situ tinggal jalan kaki 10 menit, ada pathway langsung ke venue (tidak wheelchair accessible)
                </li>
            </ol>
        </Tab>
    </Tabs>
    <div class="wheelchair-accessibility">
        <a href="#" aria-controls="wheelchair-card" class="flex items-center gap-1 mb-5" @click.prevent="showWheelchair=!showWheelchair">Wheelchair accessibility <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"><g fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="1.5"/><path stroke="currentColor" stroke-linecap="round" stroke-width="1.5" d="M10.125 8.875a1.875 1.875 0 1 1 2.828 1.615c-.475.281-.953.708-.953 1.26V13"/><circle cx="12" cy="16" r="1" fill="currentColor"/></g></svg></a>
        <div class="card card-translucent mb-3" v-show="showWheelchair" id="wheelchair-card">
            <h3 class="mb-3">Wheelchair Accessibility</h3>
            <p>
                Kalau diantar, minta diturunin persis di depan venue. dari pintu masuk parkir (ambil tiket), naik ke atas ramp 1 kali, abis itu di ramp ke 2, ambil kanan (jangan naik ramp kedua), lokasi ada di sebelah kanan persis
            </p>
        </div>
    </div>
</template>
<style>
.tab-venue li {
    margin-bottom: 1rem;
}

</style>