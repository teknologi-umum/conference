<script lang="ts" setup>
useHead({
    title: "Verify"
})

import {QrcodeStream} from 'vue-qrcode-reader'
const detectedUser = ref()
const detectedUsers = ref([])
const key = ref([])
const config = useRuntimeConfig();
const alertClass = ref<null|boolean>(null);
const paused = ref(false)
const invalidTicketReason = ref<string>("");
interface ScanResponse {
    message: string
    student: boolean
    name: string
    email: string
}
const onDetect = async (a: any) => {
    const response = await useFetch<ScanResponse>(`${config.public.backendBaseUrl}/scan-tiket`, { 
        method: "POST", 
        body: {
            code: a[0].rawValue,
            key: key.value
        }
    });
    paused.value = true

    if (response.error.value?.statusCode && [406, 403].includes(response.error.value?.statusCode)) {
        alertClass.value = false
        invalidTicketReason.value = response.error.value?.data.errors;
    } else {
        const body = response.data.value;
        alertClass.value = true
        detectedUser.value = {
            name: body?.name,
            student: body?.student,
            email: body?.email,
        }
    }
    setTimeout(() => {
        paused.value = false
    }, 100);
}
const scanNext = () => {
    detectedUsers.value.push(detectedUser.value)
    detectedUser.value = null
}

</script>
<template>
    <div id="page">
        <SinglePage title="Verify Guest">
            <div :class="[`alert mb-5`, alertClass ? 'alert-success' : 'alert-danger']" v-if="alertClass !== null">
                {{ alertClass ? 'User verified!' : invalidTicketReason }}
            </div>
            <template v-if="!detectedUser">
                <input type="text" class="form-control-lg mb-5" placeholder="Key" v-model="key">
                <qrcode-stream  @detect="onDetect" :paused="paused"></qrcode-stream>
            </template>
            <div class="success" v-else>
                <Card>
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