<script lang="ts">
    import Navbar from "$lib/components/navbar.svelte";
    import Grid from "$lib/components/grid.svelte";
    import VideoCard from "$lib/components/video_card.svelte";

    export let data;

    let videos = data.videos;
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
