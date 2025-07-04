<script>
	import PricingCard from "./PricingCard.svelte";

    var {onCancel} = $props()

</script>
<div class=" h-screen w-screen bg-black opacity-80 fixed z-10"></div>
    <div class=" fixed inset-0 flex justify-center items-center z-20">
        <button onclick={onCancel}><img src="/cancel.svg" alt="" class="fixed top-10 right-10 h-10 w-10 hover:cursor-pointer"></button>
        <div class="card-group h-150 w-4/6 flex justify-around items-center z-30">

            <PricingCard
                planName={"Newbie Plan"}
                description={"Cocok untuk pemula yang sering upload seminggu satu video"}
                pricePerMonth={"Rp.15.000"}
                imgSvg={"/medal.svg"}
                checkLists={["Upload satu video per minggu", "Diperiksa menggunakan Model AI", "Deteksi Manual"]}
                planType={"NEWBIE"}
            />

            <PricingCard 
                planName={"Creator Plan"} 
                description={"Cocok untuk kamu seorang kreator sejati"}
                pricePerMonth={"Rp.20.000"}
                imgSvg={"/crown.svg"}
                checkLists={["Upload video sebebasnya", "Diperiksa dengan Model AI", "Deteksi Otomatis", "Bebas atur jadwal deteksi"]}
                planType={"CREATOR"}
                />
        </div>
    </div>
