<script lang="ts">
	import Check from "../components/Check.svelte";
	import PricingCard from "../components/PricingCard.svelte";

    console.log(import.meta.env)
    let openIndex: number | null = null;

    const faqs = [
        {
            question:"Apa itu Comdel?",
            answer:"Comdel singkatan 'delete comment' merupakan sebuah website yang membantu konten kreator di indonesia untuk menghapus komentar yang terindikasi mempromosikan website Judi Online secara otomatis dan terjadwal menggunakan bantuan AI."
        },

        {
            question:"Data apa saja yang diakses oleh Comdel?",
            answer:"Data yang kami kumpulkan adalah profil youtube seperti username dan foto profil dan juga list video supaya kami bisa memverifikasi bahwa video tersebut milik anda"
        },

        {
            question:"Apakah ketika sudah membeli dapat refund?",
            answer:"Service yang sudah dibeli tidak dapat dikembalikan kembali"
        },

        {
            question:"Apakah tidak ada cara lain untuk login selain menggunakan Google?",
            answer:"Tidak, login menggunakan akun google merupakan satu satunya cara paling aman untuk autentikasi dan juga akun google dapat terhubung langsung ke channel youtube milik anda sendiri"
        }
    ]

    function toggleFAQ(index:Number) {
		openIndex = openIndex === index ? null : index;
	}

</script>
<div class="bg-[#111111]">
    <div class="h-screen w-full bg-[#111111] text-white">
        <div class="h-1/12 w-full flex justify-between p-5">
            <p class="w-6/6 md:w-2/6 h-full flex justify-center items-center text-bla font-sans font-extrabold text-2xl">
                Comdel
            </p>
            <div class=" hidden md:w-2/6 h-full text-white text-sm md:flex justify-evenly items-center">
                <a class="border-b" href="/">Cara kerja</a>
                <a class="border-b" href="/">Model AI</a>
                <a class="border-b" href="/">Keamanan Data</a>
            </div>
            <div class="hidden md:w-2/6 h-full md:flex justify-center items-center">
                <a href="/" class=" justify-center items-baseline rounded-lg text-base p-2 bg-blue-600">Berlangganan</a>
            </div>
        </div>
        <div class=" h-11/12 w-full flex justify-center items-center">
            <div class="w-5/6   md:w-3/6 h-5/6 flex items-center justify-center flex-col p-5">
                <p class="w-4/6 text-4xl md:text-5xl text-center font-semibold">Hapus Komentar <br> <s class="text-[#555]">Judi Online</s> Di kolom komentarmu Secara Otomatis</p>
                <p class="text-center mt-7 text-md md:text-lg">Dibantu dengan teknologi <span class="font-extrabold">Kecerdasan Buatan (AI)</span> dan <span class="font-extrabold">Google Youtube API</span> untuk menciptakan lingkungan digital yang lebih aman</p>
                <a href="/" class="mt-7 text-white font-bold p-3 rounded-lg bg-blue-600 hover:cursor-pointer hover:bg-blue-800">Coba Sekarang</a>
            </div>
        </div>

    </div>

    <!-- 2nd Section -->
    <div class="h-screen w-full bg-[#111111] text-white">

        <div class="h-1/12 w-full">
            <p class="text-center text-2xl md:text-4xl text-[#888] font-extrabold font-mono">KENAPA HARUS COMDEL</p>
        </div>

        <div class=" h-11/12 w-full flex flex-col md:flex-row items-center justify-center md:justify-between p-2 md:p-5">
            <div class="h-1/6 md:h-full w-5/6 md:w-2/6 flex flex-col justify-center">
                <h1 class="text-3xl font-sans font-bold text-[#EEEEEE]">User Interface Friendly</h1>
                <h1 class=" text-lg font-sans text-[#EEEEEE] mt-5">Antarmuka dibuat menyerupai halaman <span class="font-extrabold">YouTube</span> sehingga terlihat familiar bagi pengguna</h1>
            </div>
    
            <div class="h-3/6 md:h-full w-full md:w-4/6 flex justify-center items-center">
                <img src="/banner.png" alt="" class="w-5/6 h-4/6 bg-white rounded-lg shadow-lg shadow-white">
            </div>
        </div>
    </div>

    <div class="h-fit w-full bg-[#111111] text-[#DDDDDD]">
        <div class="h-full md:h-80 w-full flex flex-col md:flex-row justify-evenly items-center">
            <div class="h-5/6 w-5/6 md:w-3/12 p-5">
                <img src="/youtube2.svg" alt="" class="h-20 w-20">
                <p class="text-xl font-extrabold text-white mt-4">API Resmi YouTube</p>
                <p class="mt-4">Menggunakan Layanan resmi <span class="font-extrabold">YouTube API</span> dari <span class="font-extrabold">Google</span> sehingga menjadi keselamatan data anda</p>
            </div>
            
            <div class="h-5/6 w-5/6 md:w-3/12 p-5">
                <img src="/brain.svg" alt="" class="h-20 w-20">
                <p class="text-xl font-extrabold text-white mt-4">Menggunakan AI</p>
                <p class="mt-4">Dilengkapi dengan <span class="font-extrabold">Artificial Intelligence</span> untuk mendeteksi Komentar</p>
            </div>

            <div class="h-5/6 w-5/6 md:w-3/12 p-5">
                <img src="/time.svg" alt="" class="h-20 w-20">
                <p class="text-xl font-extrabold text-white mt-4">Penghapusan Terjadwal</p>
                <p class="mt-4">Komentar terdeteksi akan dihapus secara terjadwal, tidak perlu khawatir untuk check ulang</p>
            </div>
        </div>
    </div>

    <div class="h-fit w-full bg-[#111111] flex flex-col items-center justify-center mt-10 md:mt-0">
        <div class="h-1/12 w-full">
            <p class="text-center text-3xl md:text-4xl text-[#888] font-extrabold font-mono">BIAYA LANGGANAN</p>
        </div>

        <div class="card-group h-250 md:h-150 w-4/6 flex flex-col md:flex-row justify-start md:justify-evenly gap-5 items-center mt-5">

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

    <div class="faq h-fit w-full text-white mt-5 flex flex-col justify-center items-center">
        <h1 class="text-3xl md:text-4xl text-[#888] text-center font-mmono font-extrabold">FAQ</h1>
        <div class="w-3/6 mx-auto">
            {#each faqs as faq, index}
                <div class="mb-4 border-b border-gray-600">
                <button
                    class="w-full text-left py-4 font-semibold"
                    on:click={() => toggleFAQ(index)}
                >
                    {faq.question}
                </button>
                <div class={`faq-answer px-4 pb-4 text-sm text-gray-300 ${openIndex === index ? 'block' : 'hidden'}`}>
                    {faq.answer}
                </div>
                </div>
            {/each}
        </div>
    </div>

    <div class="footer w-full bg-black text-[#DDD] py-8 px-4">
        <div class="grid justify-items-center grid-cols-2 sm:grid-cols-2 md:grid-cols-4">
            <div>
                <h3 class="text-lg font-semibold mb-2">Layanan</h3>
                <ul class="space-y-1 text-sm">
                    <li><a href="/tos" class="hover:underline">Terms of Service</a></li>
                    <li><a href="/privacy" class="hover:underline">Privacy Policy</a></li>
                    <li><a href="/help" class="hover:underline">Help Center</a></li>
                </ul>
            </div>

            <div>
                <h3 class="text-lg font-semibold mb-2">Media Sosial</h3>
                <ul class="space-y-i text-sm">
                    <li><a href="/instagram" class="hover:underline">Instagram</a></li>
                </ul>
            </div>

            <div>
                <h3 class="text-lg font-semibold mb-2">Tentang</h3>
                <ul class="space-y-1 text-sm">
                    <li><a href="/about" class="hover:underline">Tentang Kami</a></li>
                    <li><a href="/career" class="hover:underline">Karir</a></li>
                    <li><a href="/contact" class="hover:underline">Kontak</a></li>
                </ul>
            </div>
            
            <div>
                <h3 class="text-lg font-semibold mb-2">Alamat</h3>
                <p class="text-sm">
                    Jl. Damai No. 123<br />
                    Jakarta, Indonesia<br />
                    support@email.com
                </p>
            </div>
        </div>

        <div class="mt-8 text-center text-xs text-gray-500">
            &copy; {new Date().getFullYear()} NamaPerusahaan. All rights reserved.
        </div>
    </div>


</div>

<style>
  .faq-answer {
    transition: all 0.3s ease;
  }
</style>
