<script>
	import { redirect } from "@sveltejs/kit";
	import Notification from "../../components/Notification.svelte";

    let tos = $state(false)
    let pp = $state(false)
    let buttonState = $state(false)
    let showNotif = $state(false);

    const onTOSChange = () => {
        tos = !tos
    }

    const onPPChange = () => {
        pp = !pp
    }

    const onButtonClicked = () => {
        if (!tos && !pp) {
            showNotif = true;
        }

        window.location.href = "http://localhost:8080/auth/google";
    }

    const toggleOffShowNotif = () => {
        showNotif = false;
    }
</script>

<div class="h-screen w-screen bg-[#111111] flex justify-center items-center">
    {#key showNotif}
        {#if showNotif}
            <Notification
                isSuccess={false}
                status="failed"
                message="please agree on terms & service and Privacy Policy"
                callback={toggleOffShowNotif}
            />
        {/if}
    {/key}
    <div class="log-box w-2/6 h-3/6 bg-[#000003]  rounded-lg text-white p-5">
        <div class="head-log h-1/6 w-full">
            <p class="text-center text-2xl font-bold">LOGIN</p>
        </div>
        
        <div class="tos-log h-3/6 w-full">
            <p class="text-sm text-center text-[#EEEEEF]">Gunakan akun google yang terhubung dengan akun youtube mu</p>

            <div class="flex mt-10 w-full h-1/6  items-center">
                <input onchange={onTOSChange} type="checkbox" name="tos" id="tos" class="h-4 w-4">
                <p class=" ml-2 text-[#999999] text-sm">Saya sudah membaca & menyetujui <a href="/tos" class="text-white underline">Terms and Service</a> dari Comdel</p>
            </div>

            <div class="flex w-full h-1/6  items-center">
                <input onchange={onPPChange} type="checkbox" name="tos" id="tos" class="h-4 w-4">
                <p class=" ml-2 text-[#999999] text-sm">Saya sudah membaca & menyetujui <a href="/tos" class="text-white underline">Privacy Policy</a> dari Comdel</p>
            </div>
        </div>

        <div class="body-log h-2/6 w-full flex items-center justify-center">
            <button onclick={onButtonClicked} class="h-3/6 w-full border border-[#444444] hover:cursor-pointer flex items-center p-5 rounded-full">
                <img src="/google.png" alt="" class="h-9 w-9">
                <p class="text-center w-full">Login using google account</p>
            </button>
        </div>

    </div>
</div>