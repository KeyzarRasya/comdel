<script>
    var {src, title, uploadTime, vid} = $props();
</script>

<a href="/content/{vid}" class="w-11/12 h-70 p-2">
    <img src={src} alt="" class=" rounded-xl ">
    <p class="mt-2 text-white text-sm">{title}</p>
</a>