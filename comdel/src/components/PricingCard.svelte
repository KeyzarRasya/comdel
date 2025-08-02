<script>
    import Check from "./Check.svelte";

    var {planName, description, pricePerMonth, imgSvg, checkLists, planType} = $props()

    
</script>

<div class="h-3/6 md:h-5/6 w-5/6 md:w-2/6 bg-gradient-to-bl from-[#222222] to-[#333333] rounded-lg flex flex-col p-3">

    <div class="title-pricing h-10/12 md:h-6/12   border-white">
        <img src={imgSvg} alt="" class="h-12 w-12">
        <p class="text-[#EEEEEE] text-2xl mt-5 font-extrabold">{planName}</p>
        <p class="text-xs text-[#DDDDDD] mt-3 font-sans">{description}</p>

        <p class="text-4xl mt-5 font-extrabold">
        <span class="text-white">
               {pricePerMonth}
        </span>
        <span class="text-[#777] font-bold text-base">
              / Bulan
        </span>
        </p>
    </div>

    <a href="http://localhost:8080/payment/pay/?plan={planType}" class="w-full h-12 bg-gradient-to-br from-sky-700 via-sky-600 to-sky-500 hover:from-sky-500 hover:via-sky-600 hover:to-sky-700 hover:cursor-pointer text-white flex justify-center items-center font-bold rounded">Dapatkan Sekarang</a>

    <div class="flex items-center gap-4 mt-3">
        <div class="flex-grow border-t border-white"></div>
        <span class="text-white font-semibold text-xs">Fitur</span>
        <div class="flex-grow border-t border-white"></div>
    </div>

    <div class="specification-pricing">
        <div class="mt-2"></div>
        {#each checkLists as c}
            <Check text={c} />
        {/each}
    </div>
</div>            