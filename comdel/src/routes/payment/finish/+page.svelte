<script>
// @ts-nocheck

    export let data;

    console.log(data)
</script>

<div class=" h-screen w-screen bg-[#222222] flex justify-center items-center">
    <div class="h-5/6 w-3/6 flex justify-center items-center flex-col">
        {#if data.transactionStatus === "capture"}    
            <img src="/success.svg" alt="" class="h-2/6 w-2/6">
            <p class="text-white font-extrabold text-3xl">Transaksi Berhasil</p>
        {:else if data.transactionStatus === "pending"}
            <img src="/time.svg" alt="" class="h-2/6 w-2/6">
            <p class="text-white font-extrabold text-3xl">Menunggu Pembayaran</p>
        {:else}
            <img src="/cancel.svg" alt="" class="h-2/6 w-2/6">
            <p class="text-white font-extrabold text-3xl">Transaksi Gagal</p>
        {/if}
        <p class="text-white font-bold text-sm mt-5">Order Id : {data.orderId}</p>
        <a href="/dashboard" class="mt-10 text-white p-3 rounded-lg font-bold bg-blue-600">Menuju Dashboard</a>
    </div>
</div>