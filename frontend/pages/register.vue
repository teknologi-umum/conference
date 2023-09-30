<script setup lang="ts">
const fullName = ref<string>("");
const email = ref<string>("");
const config = useRuntimeConfig();
const alert = reactive({
    type: "",
    msg: ""
})

const submit = async () => {
    const response = await useFetch(`${config.public.backendBaseUrl}/users`, { 
        method: "POST", 
        body: {
            email: email.value,
            name: fullName.value
        }
    });
    alert.type = 'success'
    alert.msg = "Registration success. You will receive an invitation via email within 7 business days."
    if (response.error.value != null) {
        alert.type = 'danger'
        alert.msg = response.error.value?.data?.message
        
        if (response.error.value?.statusCode == 400) {
            alert.msg = "Please check your input"
        } else if (response.error?.value.statusCode == 406) {
            alert.msg = "Sorry, registration has been closed. We're no longer accepting any attendee registration."
        }

        return;
    }
    

    email.value = ''
    fullName.value = ''
}
</script>

<template>
    <div id="page">
        <SinglePage title="Save your spot!">
            <p class="desc">We only have a limit for 50 participants. Reserve yours now. </p>
            <div :class="`alert alert-${alert.type} mb-5`" v-if="alert.type !== ''">{{ alert.msg }}</div>
            <form @submit.prevent="submit" action="" class="max-w-[500px] mb-24">
                <div class="form-group mb-5">
                    <label for="full-name">Full name</label>
                    <input type="text" id='full-name' class="form-control-lg" placeholder="Juned" v-model="fullName">
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
</style>
