<script>
    import VideoCard from "./VideoCard.svelte";

    const {videos} = $props();
</script>

<div class="w-full grid grid-cols-3 justify-items-center p-5 overflow-y-auto" style="flex: 1;">
    {#each videos as video}
        <VideoCard src={video.thumbnail} title={video.title} uploadTime="" vid={video.id}/>
    {/each}
</div>