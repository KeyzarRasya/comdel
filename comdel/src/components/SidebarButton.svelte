<script>
    var {imgSrc, text, onClick, highlight=false} = $props();
</script>

<button class="hover:cursor-pointer h-12 flex items-center p-5 text-[#EEEEEE] font-sans hover:bg-[#111111]" onclick={onClick}>
    <img src={imgSrc} class="w-5 h-5" alt=""> 
    {#if highlight}
        <p class="ml-3 font-bold">{text}</p>
    {:else}
        <p class="ml-3">{text}</p>
    {/if}
</button>