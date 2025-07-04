<script>
    // @ts-ignore
    let {name, comments, profileUrl, publishedDate} = $props();

    // @ts-ignore
    function selisihHariDariHariIni(dateString) {
        const inputDate = new Date(dateString);
        const today = new Date();

        // Hilangkan jam agar perbandingan hanya berdasar hari
        inputDate.setHours(0, 0, 0, 0);
        today.setHours(0, 0, 0, 0);

        const msPerHari = 24 * 60 * 60 * 1000;
        // @ts-ignore
        const selisih = Math.floor((today - inputDate) / msPerHari);

        if (selisih === 0) return 'hari ini';
        if (selisih > 0 && selisih < 7) return `${selisih} hari yang lalu`;
        if (selisih >= 7 && selisih < 30) return `${Math.floor(selisih / 7)} minggu yang lalu`;
        if (selisih >= 30 && selisih < 365) return `${Math.floor(selisih / 30)} bulan yang lalu`;
        if (selisih >= 365) return `${Math.floor(selisih / 365)} tahun yang lalu`
        return `dalam ${Math.abs(selisih)} hari`;
    }

</script>

<div class="h-fit w-full  p-2 flex justify-between mt-1">
    <img src={profileUrl} alt="" class="h-10 w-10 bg-white rounded-full">
    <div class="w-11/12 ml-2">
        <div class=" flex items-center">
            <p class="text-white text-xs font-bold">{name}</p>
            <p class="text-[#999] text-xs ml-2">{selisihHariDariHariIni(publishedDate)}</p>
        </div>
        <p class="text-[#BBB] text-sm text-justify pr-2 mt-1">{comments}</p>
    </div>
</div>