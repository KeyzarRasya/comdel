<script lang="ts">
	import Comment from '../../../components/Comment.svelte';

    let {data} = $props();
    let isComment = $state(true);
    let isHaveComment = $state(true);
    let video = data.result.data;

    let detectedCount = 0;
    let notDetectedCount = 0;

    let commentLists = [
        {
            name:"KeyzarRasya",
            comments:"lorem ipsum banget nihh"
        },
        {
            name:"Oppenheimer",
            comments:"Lorem ipsum dolor sit amet, consectetur adipisicing elit. Id, modi perferendis vel cum praesentium excepturi consectetur molestiae error, cumque minima, autem ut voluptatum quas odit provident eum amet! Sit, aperiam."
        }
    ]

    let deletedComment = [
        {
            name:"BujangLapuk",
            comments:"main du putrapetir77"
        }
    ]

    let list: Array<any> = video.comments;

    list.forEach(element => {
        if (element.isDetected) {
            detectedCount++;
        } else {
            notDetectedCount++;
        }
    });

    console.log(detectedCount)

    if (list == null) {
        isHaveComment = false;
    }

    const onCommentClick = () => {
        isComment = true;
    }

    const onDeletedCommentClick = () => {
        isComment = false;
    }


</script>
  


<div class="h-screen w-screen bg-[#222222] flex flex-col">
    <div class="h-1/12 w-full border-b border-[#444444] flex">
        <div class="on-left w-1/12 h-full flex justify-center items-center">
            <a href="/dashboard" class="h-4/6 w-5/6 hover:cursor-pointer flex justify-center items-center">
                <img src="/back2.svg" alt="back" class="w-full h-4/6">
            </a>
        </div>

        <div class="on-right">

        </div>
    </div>

    <div class="flex h-11/12">
        <!-- Left section: 4/6 width -->
        <div class="w-3/6 h-full border-r border-[#444444] p-3">
            <img src={video.thumbnail} alt="thumbnail" class="thumbnail w-full h-8/12 bg-white rounded-xl border border-[#444]">
            <p class="mt-5 text-white text-base font-bold text-center">{video.title}</p>
            
            <div class="h-3/12 w-full flex justify-around items-center p-5">
                <div class="h-4/6 w-2/6 flex flex-col justify-around items-center text-white">
                    <p class="text-4xl font-bold">{notDetectedCount}</p>
                    <p class="text-sm">Komentar</p>
                </div>
                <div class="h-4/6 w-2/6 flex flex-col justify-around items-center text-white">
                    <p class="text-4xl font-bold">{detectedCount}</p>
                    <p class="text-sm">Komentar Terdeteksi</p>
                </div>
            </div>
        </div>
    
        <!-- Right section: 2/6 width -->
        <div class="side-comment-bar w-3/6 h-full border-b border-[#333333]">
            <!-- Comment Tab Headers -->
            <div class="headers h-20 border-b border-[#444444] flex">
                {#if isComment}
                    <button onclick={onCommentClick} class="hover:cursor-pointer w-3/6 h-full flex flex-col justify-center p-2 border-b-blue-500 border-b items-center">
                    <p class="text-white font-bold">Comments</p>
                    <p class="text-[#999999] text-xs mt-1">List of all your comments</p>
                </button>
                <button onclick={onDeletedCommentClick} class="hover:cursor-pointer w-3/6 h-full flex flex-col justify-center p-2 items-center">
                    <p class="text-white font-bold">Deleted Comments</p>
                    <p class="text-[#999999] text-xs mt-1">The comment that we deleted</p>
                </button>

                {:else}
                    <button onclick={onCommentClick} class="hover:cursor-pointer w-3/6 h-full flex flex-col justify-center p-2  items-center">
                        <p class="text-white font-bold">Comments</p>
                        <p class="text-[#999999] text-xs mt-1">List of all your comments</p>
                    </button>
                    <button onclick={onDeletedCommentClick} class="hover:cursor-pointer w-3/6 h-full flex flex-col justify-center p-2 border-b-blue-500 border-b items-center">
                        <p class="text-white font-bold">Deleted Comments</p>
                        <p class="text-[#999999] text-xs mt-1">The comment that we deleted</p>
                    </button>
                {/if}
            </div>

            <!-- Comment List -->
            <div class="comment-list mt-2 overflow-y-auto h-[calc(100%-5rem)] pr-2">
                {#each list as c}
                    {#if (isComment && !c.isDetected) || (!isComment && c.isDetected)}
                        <Comment
                            name={c.displayName}
                            comments={c.textDisplay}
                            profileUrl={c.profileUrl}
                            publishedDate={c.publishedAt}
                        />
                    {/if}
                {/each}
            </div>

        </div>
    </div>
</div>

<style>
    .thumbnail{
        background-image: url("/thumbnail.png");
        background-position: center;
        background-size: cover;
    }
</style>