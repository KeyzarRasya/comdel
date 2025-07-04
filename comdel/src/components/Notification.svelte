<script>
    import { fly, fade } from "svelte/transition";
    import { onDestroy } from 'svelte';

    let isVisible = true; // Start visible to trigger the animation
    export let isSuccess = true;
    export let status = '';
    export let message = '';
    export let callback = () => {};

    const closeNotification = () => {
        isVisible = false;
        callback()
    };

    // Auto-hide after 5 seconds
    const timer = setTimeout(() => {
        isVisible = false;
    }, 3000);

    onDestroy(() => clearTimeout(timer));
</script>

    {#if isVisible}
        <div
            out:fly={{y: 200, duration: 1000 }}
            in:fly={{y: -200, duration: 1000 }}
            class="h-fit w-65 bg-[#333333] fixed right-2 top-[20%] rounded-lg p-2 flex border-b-2 z-50"
            class:border-green-500={isSuccess}
            class:border-red-500={!isSuccess}
        >
            <div class="w-5/6 h-full">
                <p class:text-green-400={isSuccess} class:text-red-400={!isSuccess} class="font-bold">
                    {status.toUpperCase()}
                </p>
                <p class="text-[#EEEEEE] text-xs mt-1">{message}</p>
            </div>
            <button class="w-1/6 flex justify-center items-center" on:click={closeNotification}>
                <img src="/cancel.svg" alt="Close" class="w-5 h-5">
            </button>
        </div>
    {/if}
