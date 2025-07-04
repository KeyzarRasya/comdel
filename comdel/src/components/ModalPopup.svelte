<script lang="ts">
	import Notification from "./Notification.svelte";


    /* State */
    var modalContent = $state(1);
    var isLinkValid = $state(false)
    var link = $state("");
    var scheduler = $state("24")
    var isOwnershipValid = $state(false);
    var isVerifyClicked = $state(false);
    var verifyClicked = $state(0)
    var messageCheck = $state("")
    
    var uploadResult = $state({
        status:400,
        message:"Failed to fetch",
        data:null
    });
    var doneFetch = $state(false);

    
    let strategy = $state("AUTO");

    var {isHide, onClose, onItemAdded} = $props();

    const isYoutubeLink = (url: string) =>
  /^https?:\/\/(www\.)?(youtube\.com|youtu\.be)\/.+$/.test(url);

    const onBack = () => {
        modalContent--;

        if (modalContent === 1) {
            isOwnershipValid = false
            isVerifyClicked = false;
            verifyClicked = 0;
        }
    }

    const onNext = () => {
        modalContent++;
    }

    const onCheckOwnership = async () => {
        const response = await fetch(`http://backend:8080/videos/ownership/?vid=${link}`, {
            method:"GET",
            credentials:"include"
        })

        const result = await response.json()

        console.log(result)
        
        if (result.status === 200) {
            isOwnershipValid = true;
        } else {
            isOwnershipValid = false
        }
        isVerifyClicked = true;
        verifyClicked++;
        messageCheck = result.message
        console.log(result)
    }

    const onInputLink = (e) => {
        link = e.target.value;
        isLinkValid = isYoutubeLink(link)
    }

    const onSchedulerChange = (e) => {
        scheduler = e.target.value
    }

    const onSubmitLink = async(e) => {

        const response = await fetch(`http://backend:8080/videos/upload/?vid=${link}&st=${strategy}&sc=${scheduler}`, {
            method:"POST",
            credentials:"include"
        })

        const result = await response.json();
        console.log(result);
        
        uploadResult = result;
        doneFetch = true;

        if (onItemAdded) onItemAdded();
        isHide = true;

    }
    
</script>

{#key verifyClicked}
    {#if isVerifyClicked}
        <Notification
            status={isOwnershipValid ? "Success" : "Failed"}
            isSuccess={isOwnershipValid}
            message={messageCheck}
        />
    {/if}
{/key}


{#if doneFetch}
    <Notification status={uploadResult.status === 200 ? "Success" : "Failed"} isSuccess={uploadResult.status === 200? true : false} message={uploadResult.message}/>
{/if}
{#if !isHide}
    <div class="h-screen w-screen bg-black opacity-50 fixed z-10"></div>
    <div class="fixed inset-0 flex justify-center items-center z-20">
        
        <!-- Konten modal -->
        
        <!-- Insert Video Link -->
        {#if modalContent === 1}
            <div class="modal h-3/6 w-2/6 bg-[#222222] p-5 rounded-lg">
                <p class="text-white text-xl font-bold">Tambah Video</p>
                <div class="h-3/6 w-full flex items-center">
                    <div class="relative z-0 w-full group">
                    <input type="text" name="floating_email" id="floating_email" class="block py-2.5 px-0 w-full text-sm text-gray-900 bg-transparent border-0 border-b-2 border-gray-300 appearance-none dark:text-white dark:border-gray-600 dark:focus:border-blue-500 focus:outline-none focus:ring-0 focus:border-blue-600 peer" placeholder=" " required oninput={onInputLink} value={link}/>
                    <label for="floating_email" class="peer-focus:font-medium absolute text-sm text-gray-500 dark:text-gray-400 duration-300 transform -translate-y-6 scale-75 top-3 -z-10 origin-[0] peer-focus:start-0 rtl:peer-focus:translate-x-1/4 rtl:peer-focus:left-auto peer-focus:text-blue-600 peer-focus:dark:text-blue-500 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-6">Link Video Youtube</label>
                </div>
            </div>
    
                <div class=" h-2/6 flex justify-end items-end">
                    {#if isOwnershipValid}    
                        <div class="w-3/6 h-3/6  flex items-end justify-between">   
                            <button class="hover:cursor-pointer bg-white font-bold text-blue-600 p-2 w-25 border-2 border-blue-800 rounded-full" onclick={onClose}>Cancel</button>
                            <button class="hover:cursor-pointer  p-2 w-25 text-white font-bold rounded-full"
                            class:bg-blue-900={!isLinkValid}
                            class:bg-blue-600={isLinkValid}
                            class:hover:cursor-pointer={isLinkValid}
                            disabled={!isLinkValid}
                            onclick={onNext}>Next</button>
                        </div>
                    {:else}
                        <div class="w-3/6 h-3/6  flex items-end justify-between">   
                            <button class="hover:cursor-pointer bg-white font-bold text-blue-600 p-2 w-25 border-2 border-blue-800 rounded-full" onclick={onClose}>Cancel</button>
                            <button class="hover:cursor-pointer  p-2 w-25 text-white font-bold rounded-full"
                            class:bg-blue-900={!isLinkValid}
                            class:bg-blue-600={isLinkValid}
                            class:hover:cursor-pointer={isLinkValid}
                            disabled={!isLinkValid}
                            onclick={onCheckOwnership}>Verify</button>
                        </div>
                    {/if}
                </div>
            </div>


                {:else if modalContent === 2}
                <div class="modal h-4/6 w-2/6 bg-[#222222] p-5 rounded-lg">
                    <div class="w-full h-full flex flex-col">
                        <p class="text-white text-xl font-bold">Confirmation</p>
                        <div class="h-3/6 w-5/6 bg-white mt-5 self-center rounded-md">
        
                        </div>
                        <p class="w-5/6 self-center mt-3 text-[#DDDDDD] text-sm">Tutorial pointer pemula</p>

                        <div class=" h-2/6 flex justify-end items-end">
                            <div class="w-3/6 h-3/6  flex items-end justify-between">   
                                <button class="hover:cursor-pointer bg-white font-bold text-blue-600 p-2 w-25 border-2 border-blue-800 rounded-full" onclick={onBack}>Back</button>
                                <button class="hover:cursor-pointer bg-blue-600 p-2 w-30 text-white font-bold rounded-full" onclick={onNext}>Konfirmasi</button>
                            </div>
                        </div>
                    </div>
                </div>

                {:else if modalContent === 3}
                <div class="modal h-4/6 w-2/6 bg-[#222222] p-5 rounded-lg">
                    <p class="text-xl text-white font-bold">Scheduler</p>

                    <div class="w-full h-4/6 flex flex-col justify-evenly">
                        <div>
                            <label for="scheduler" class="text-white text-base">Pilih interval waktu untuk AI kami mendeteksi sebuah komentar<span class="text-[#FF0000]">*</span></label>
                            <p class="text-[#888888] text-xs">Anda bisa mengubah ini kapan saja</p>
                            <select
                            name="scheduler"
                            id="scheduler"
                            class="mt-5 w-40 text-white bg-transparent  border border-blue rounded-md focus:ring focus:ring-blue-500"
                            onchange={onSchedulerChange}
                            >
                                <option class="text-black bg-white" value="24">1 Hari</option>
                                <option class="text-black bg-white" value="12">12 Jam</option>
                                <option class="text-black bg-white" value="6">6 Jam</option>
                                <option class="text-black bg-white" value="3">3 Jam</option>
                            </select>
                        </div>

                        <div class="flex flex-col">
                            <p class="text-white text-base">Pilih strategi untuk menghapus komentar<span class="text-[#FF0000]">*</span></p>
                            <div class="flex mt-5">
                                
                                <div class="w-3/6">
                                    <input type="radio" name="strategy" id="AUTO" value="AUTO" bind:group={strategy}>
                                    <label class="text-white ml-4" for="otomatis">Otomatis</label>
                                </div>
                                <div class="w-3/6">
                                    <input type="radio" name="strategy" id="MANUAL" value="MANUAL" bind:group={strategy}>
                                    <label class="text-white ml-4" for="manual">Manual</label>
                                </div>
                            </div>
                        </div>

                    </div>

                    <div class=" h-1/6 flex justify-end items-end">
                        <div class="w-3/6 h-3/6  flex items-end justify-between">   
                            <button class="hover:cursor-pointer bg-white font-bold text-blue-600 p-2 w-25 border-2 border-blue-800 rounded-full" onclick={onBack}>Back</button>
                            <button class="hover:cursor-pointer bg-blue-600 p-2 w-25 text-white font-bold rounded-full" onclick={onSubmitLink}>Selesai</button>
                        </div>
                    </div>

                    
                </div>



             {/if}

    </div>
{/if}