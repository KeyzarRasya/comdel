<script>
	import Check from "../components/Check.svelte";
	import PricingCard from "../components/PricingCard.svelte";

    console.log(import.meta.env)

</script>
<div>
    <div class="h-screen w-screen bg-[#111111] text-white">
        <div class="h-1/12 w-full flex justify-between p-5">
            <p class="w-2/6 h-full flex justify-center items-center text-bla font-sans font-extrabold text-2xl">
                Comdel
            </p>
            <div class=" w-2/6 h-full text-white text-sm flex justify-evenly items-center">
                <a class="border-b" href="/">Cara kerja</a>
                <a class="border-b" href="/">Model AI</a>
                <a class="border-b" href="/">Keamanan Data</a>
            </div>
            <div class="w-2/6 h-full flex justify-center items-center">
                <a href="/" class=" justify-center items-baseline rounded-lg text-base p-2 bg-blue-600">Berlangganan</a>
            </div>
        </div>
        <div class=" h-11/12 w-full flex justify-center items-center">
            <div class="w-3/6 h-5/6 flex items-center justify-center flex-col p-5">
                <p class="w-4/6 text-5xl text-center font-semibold">Hapus Komentar <br> <s class="text-[#555]">Judi Online</s> Di kolom komentarmu Secara Otomatis</p>
                <p class="text-center mt-7 text-lg">Dibantu dengan teknologi <span class="font-extrabold">Kecerdasan Buatan (AI)</span> dan <span class="font-extrabold">Google Youtube API</span> untuk menciptakan lingkungan digital yang lebih aman</p>
                <a href="/" class="mt-7 text-white font-bold p-3 rounded-lg bg-blue-600 hover:cursor-pointer hover:bg-blue-800">Coba Sekarang</a>
            </div>
        </div>

    </div>

    <!-- 2nd Section -->
    <div class="h-screen w-screen bg-[#111111] text-white">

        <div class="h-1/12 w-full">
            <p class="text-center text-4xl text-[#888] font-extrabold font-mono">KENAPA HARUS COMDEL</p>
        </div>

        <div class=" h-11/12 w-full flex justify-between p-5">
            <div class="h-full w-2/6 flex flex-col justify-center">
                <h1 class="text-3xl font-sans font-bold text-[#EEEEEE]">User Interface Friendly</h1>
                <h1 class=" text-lg font-sans text-[#EEEEEE] mt-5">Antarmuka dibuat menyerupai halaman <span class="font-extrabold">YouTube</span> sehingga terlihat familiar bagi pengguna</h1>
            </div>
    
            <div class="h-full w-4/6 flex justify-center items-center">
                <img src="/banner.png" alt="" class="w-5/6 h-4/6 bg-white rounded-lg shadow-lg shadow-white">
            </div>
        </div>
    </div>

    <div class="h-fit w-screen bg-[#111111] text-[#DDDDDD]">
        <div class="h-80 w-full flex justify-evenly">
            <div class="h-5/6 w-3/12 p-5">
                <img src="/youtube2.svg" alt="" class="h-20 w-20">
                <p class="text-xl font-extrabold text-white mt-4">API Resmi YouTube</p>
                <p class="mt-4">Menggunakan Layanan resmi <span class="font-extrabold">YouTube API</span> dari <span class="font-extrabold">Google</span> sehingga menjadi keselamatan data anda</p>
            </div>
            
            <div class="h-5/6 w-3/12 p-5">
                <img src="/brain.svg" alt="" class="h-20 w-20">
                <p class="text-xl font-extrabold text-white mt-4">Menggunakan AI</p>
                <p class="mt-4">Dilengkapi dengan <span class="font-extrabold">Artificial Intelligence</span> untuk mendeteksi Komentar</p>
            </div>

            <div class="h-5/6 w-3/12 p-5">
                <img src="/time.svg" alt="" class="h-20 w-20">
                <p class="text-xl font-extrabold text-white mt-4">Penghapusan Terjadwal</p>
                <p class="mt-4">Komentar terdeteksi akan dihapus secara terjadwal, tidak perlu khawatir untuk check ulang</p>
            </div>
        </div>
    </div>

    <div class="h-fit w-screen bg-[#111111] flex flex-col items-center justify-center ">
        <div class="h-1/12 w-full">
            <p class="text-center text-4xl text-[#888] font-extrabold font-mono">BIAYA LANGGANAN</p>
        </div>

        <div class="card-group h-150 w-4/6 flex justify-around items-center">

            <PricingCard
                planName={"Newbie Plan"}
                description={"Cocok untuk pemula yang sering upload seminggu satu video"}
                pricePerMonth={"Rp.15.000"}
                imgSvg={"/medal.svg"}
                checkLists={["Upload satu video per minggu", "Diperiksa menggunakan Model AI", "Deteksi Manual"]}
            />

            <PricingCard 
                planName={"Creator Plan"} 
                description={"Cocok untuk kamu seorang kreator sejati"}
                pricePerMonth={"Rp.20.000"}
                imgSvg={"/crown.svg"}
                checkLists={["Upload video sebebasnya", "Diperiksa dengan Model AI", "Deteksi Otomatis", "Bebas atur jadwal deteksi"]}/>
        </div>
    </div>

</div>