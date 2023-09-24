<script lang="ts" setup>
useHead({
    title: "Verify"
})

import {QrcodeStream} from 'vue-qrcode-reader'
const detectedUser = ref()
const detectedUsers = ref([])
const onDetect = (a: any) => {
    detectedUser.value = {
        name: 'Aji',
        student: true,
        institution: "PT Mencari Cinta Sejati"
    }
    // Fetching..
}
const scanNext = () => {
    detectedUsers.value.push(detectedUser.value)
    detectedUser.value = null
}
</script>
<template>
    <div id="page">
        <SinglePage title="Verify Guest">
            <qrcode-stream v-if="!detectedUser" @detect="onDetect"></qrcode-stream>
            <div class="success" v-else>
                <Card>
                    <div class="alert alert-success mb-5">
                        User Verified!
                    </div>
                    <table class="table mb-5">
                        <tr class="text-center">
                            <td>Key</td>
                            <td>Value</td>
                        </tr>
                        <template v-for="key in (Object.keys(detectedUser) as Array<keyof typeof detectedUser>)">
                            <tr>
                                <td class="font-bold">{{ key.charAt(0).toUpperCase() + key.slice(1) }}</td>
                                <td>{{ detectedUser[key] }}</td>
                            </tr>
                        </template>
                    </table>
                    <div class="flex justify-end">

                        <Btn @click="scanNext">Scan Next &raquo;</Btn>
                    </div>
                </Card>
            </div>
        </SinglePage>
    </div>
</template>