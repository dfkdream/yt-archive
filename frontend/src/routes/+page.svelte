<script lang="ts">
    import Navbar from "$lib/navbar.svelte";
    import { type Video, Videos } from "$lib/api/video";
    import { onMount } from "svelte";
    import Grid from "$lib/grid.svelte";
    import VideoCard from "$lib/video_card.svelte";

    let videos: Video[] = [];

    onMount(async () => {
        videos = await Videos();
    });
</script>

<svelte:head>
    <title>Home - yt-archive</title>
</svelte:head>

<Navbar location="home" />
<Grid>
    {#each videos as v}
        <VideoCard video={v} showPoster showChannel needStyle />
    {/each}
</Grid>
