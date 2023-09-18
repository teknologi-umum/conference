<script setup lang="ts">

const photo = ref<File>();
const email = ref<string>("");
const config = useRuntimeConfig();
const alert = reactive({
    type: "",
    msg: ""
})

const submit = async () => {
    // TODO: change endpoint to the correct one
    const response = await useFetch(`${config.public.backendBaseUrl}/pay`, { 
        method: "POST", 
        body: {
            email: email.value,
            photo: photo.value,
            contentType: photo.value?.type ?? ""
        }
    });
    alert.type = 'success'
    alert.msg = "Upload Payment success. You can check the status via the dashboard within 7 business days."
    if (response.error.value != null) {
        // TODO: handle on upload response error
        alert.type = 'danger'

        alert.msg = response.error.value?.data?.message
        
        if(response.error.value?.statusCode == 400) {
            alert.msg = "Please check your input"
        }

        return;
    }
    

    email.value = ''
    photo.value = undefined
}

function onFileChange(e: Event) {
    const inputElement = e.target as HTMLInputElement;
    if (inputElement.files && inputElement.files.length > 0) {
        photo.value = inputElement.files[0];
    }
}
</script>

<template>
    <div id="page">
        <SinglePage title="Upload your payment proof!">
            <p class="desc">Upload now to secure your spot</p>
            <div :class="`alert alert-${alert.type} mb-5`" v-if="alert.type !== ''">{{ alert.msg }}</div>

            <form @submit.prevent="submit" action="" class="max-w-[500px] mb-24">
                <div class="form-group mb-5">
                    <label for="image">Payment proof</label>
                    <input type="file" accept="image/*" id='image' class="form-control-lg" @change="onFileChange">
                </div>
                <div class="form-group mb-8">
                    <label for="email-address">Email address</label>
                    <input type="email" id='email-address' class="form-control-lg" placeholder="juned@company.com" v-model="email">
                </div>
                <Btn size="lg">Save my spot</Btn>
            </form>

            <h2 class="mb-5">Important Notice:</h2>
            <p class="mb-5!">TeknumConf team will not contact you and ask for payment from any other medium than email. You can validate it by:</p>
            <ul class="mb-5 pl-5">
                <li>
                    Make sure the email is from conference@teknologiumum.com. To make it not to be on your spam folder, you can add it to your mail contact first.
                </li>
                <li>
                    Scan the QR code on the bottom of the email, you should be directed to https://conf.teknologiumum.com
                </li>
            </ul>
            <p>If you don't receive any email within 5 days, please contact <a href="mailto:opensource@teknologiumum.com">opensource@teknologiumum.com</a>.</p>
        </SinglePage>
    </div>
</template>
<style>
.form-control-lg {
    @apply w-full p-4 mt-2 text-lg text-white rounded-md 
        focus:outline-none 
        border-gray-500 border-1 border-solid;
    background-color: #c1c1c10e;
    transition: all .2s;;
}
.form-control-lg:focus {
    outline: none;
    background-color: #c1c1c128;
    border-color: #eeeeee;
}
.alert {
    padding: 1rem;
    border-radius: .4rem;
}
.alert-danger {
    background-color: rgba(193, 60, 60, 0.859);
}
.alert-success {
    background-color: rgb(42 141 45);
}
</style>
