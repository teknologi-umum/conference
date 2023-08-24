<script setup lang="ts">
const fullName = ref<string>("");
const email = ref<string>("");
const config = useRuntimeConfig();
const submit = async () => {
    const response = await useFetch(`${config.public.backendBaseUrl}/users`, { method: "POST" });
    if (response.error != null) {
        // TODO: handle on registration response error
        return;
    }

    // TODO: handle success response
}
</script>

<template>
    <div id="page">
        <SinglePage title="Save your spot!">
            <p class="desc">Lorem ipsum dolor sit amet, consectetur adipisicing elit. Consectetur illo dicta, dolor commodi voluptates sit. </p>

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
                    Make sure the email is from conference@teknologiumum.com
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
</style>