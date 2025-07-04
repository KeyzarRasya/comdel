

<script lang="ts">
    import VideoCard from "../../components/VideoCard.svelte";
    import SidebarButton from "../../components/SidebarButton.svelte";
    import Content from "../../components/Content.svelte";
    import Notification from "../../components/Notification.svelte";
    import ModalPopup from "../../components/ModalPopup.svelte";
	import PremiumModal from "../../components/PremiumModal.svelte";

    let {data} = $props();

    if (data.token === null) {
        console.log(data.userInfo)
    }

    let user = $state(data.userInfo.data)

    var konten: Number = $state(1);
    var hideModal: boolean = $state(true)
    var hidePremiumModal: boolean = $state(true)

    const fetchItemReload = async() => {
        const response = await fetch("/api/user");
        const result = await response.json();

        user = result.data
    }

    console.log(user)
    

    function setKonten(num: Number) {
        konten = num;
        console.log(konten);
    }

    const onTambahVideo = () => {
        if (user.premiumPlan == "NONE") {
            hidePremiumModal = false;
        }else {
            hideModal = false;
        }

    }

    const handleModalClose = () => [
        hideModal = true
    ]

    const itemAdded = () => {
        fetchItemReload();
    }

    const premiumButtonClicked = () => {
        hidePremiumModal = false;
    }

    const premiumButtonClosed = () => {
        hidePremiumModal = true;
    }

</script>

<div class="h-screen w-screen bg-[#222222] flex">
    {#if !hidePremiumModal}
        <PremiumModal onCancel={premiumButtonClosed}/>
    {/if}
    <ModalPopup isHide={hideModal} onClose={handleModalClose} onItemAdded={itemAdded}/>
    <div class=" h-full w-1/6 border-r-[0.5px] border-[#333333]">
        <div class="h-2/6 w-full border-b border-[#333333] flex flex-col justify-evenly items-center">
            <div 
                class="h-40 w-40 bg-white rounded-full"
                style="background-image: url('{user.picture}'); background-position: center; background-size: cover;">

            </div>
            <p class="font-sans font-bold text-white">{user.givenName}</p>
        </div>
        <div class="border-b border-[#333333] h-3/6 w-full flex flex-col">
            <SidebarButton imgSrc="/videos.svg" text="Content" onClick={() => setKonten(1)} highlight={konten == 1}/>
            <SidebarButton imgSrc="/analytics.svg" text="Analytics" onClick={() => setKonten(2)} highlight={konten == 2}/>
            <SidebarButton imgSrc="/time.svg" text="Scheduler" onClick={() => setKonten(3)} highlight={konten == 3}/>
            <SidebarButton imgSrc="/setting-2.svg" text="Setting" onClick={() => setKonten(4)} highlight={konten == 4}/>
        </div>
        <div class="h-1/6 w-full flex flex-col">
            <button class="hover:cursor-pointer h-12 flex items-center p-5 text-[#EEEEEE] font-sans hover:bg-[#111111]" >Feedback</button>
            <button class="hover:cursor-pointer h-12 flex items-center p-5 text-[#EEEEEE] font-sans hover:bg-[#111111]">Logout</button>
        </div>
    </div>
    <div class="h-full w-5/6 flex flex-col overflow-hidden">
        <!-- Header -->
        <div class="header-button h-15 w-full flex items-center justify-end p-5">
            <div class="badge w-1/6 h-10 flex justify-center items-center">
                {#if user.premiumPlan == "NONE"}    
                    <button onclick={premiumButtonClicked} class="bg-[#FF0000] p-2 w-30 flex items-center justify-evenly rounded-lg text-white font-bold">
                        <img src="/premium.svg" alt="" class="w-6 h-6">
                        Premium
                    </button>
                {/if}
            </div>
            <div class=" flex w-1/6 justify-between">
                <button class="text-white text-sm font-bold p-2 border-2 border-[#999999] rounded-lg w-34 hover:cursor-pointer" onclick={onTambahVideo}>
                    + Tambah Video
                </button>
                <div
                    class="h-10 w-10  rounded-full"
                    style="background-image: url('{user.picture}'); background-position: center; background-size: cover;">
                </div>

            </div>
        </div>
    
        <!-- Scrollable Grid -->
        {#if konten === 1}
            <Content videos={user.videos}/>
        {:else if konten == 3}
            <div class="w-full h-full p-5 flex flex-col">
                <label for="scheduler" class="text-white">Pilih interval waktu untuk AI Kami mendeteksi sebuah komentar</label>
                <p class="text-[#888888] text-sm">semakin sering maka akan semakin baik</p>
                <select
                name="scheduler"
                id="scheduler"
                class="mt-5 w-40 text-black bg-white border border-white focus:ring focus:ring-blue-500"
                >
                    <option class="text-black bg-white" value="1days">1 Hari</option>
                    <option class="text-black bg-white" value="12hours">12 Jam</option>
                    <option class="text-black bg-white" value="6hours">6 Jam</option>
                    <option class="text-black bg-white" value="3hours">1 Jam</option>
                </select>

                <div class="h-10 mt-10">
                    <button class="p-2 bg-[#FF0000] hover:bg-[#E00000] hover:cursor-pointer text-sm text-white font-bold rounded-lg">Save changes</button>
                </div>
            </div>
        {/if}
    </div>
    
</div>

<style>
</style>