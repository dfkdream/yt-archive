<script lang="ts">
    import type { Video } from "./api/video";
    import VideoCard from "./video_card.svelte";

    export let videos: Video[];
    export let showChannel = false;
</script>

<div class="grid gap-1 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
    {#each videos as v}
        <VideoCard video={v} {showChannel} showPoster needStyle />
    {/each}
</div>
