<script lang="ts">
    import ChannelCard from "$lib/components/channel_card.svelte";
    import Grid from "$lib/components/grid.svelte";
    import Navbar from "$lib/components/navbar.svelte";
    import VideoCard from "$lib/components/video_card.svelte";

    export let data;

    export let channelVideos = data.channelVideos;
</script>

<svelte:head>
    {#if channelVideos}
        <title>{channelVideos.Title} - yt-archive</title>
    {:else}
        <title>{data.id} - yt-archive</title>
    {/if}
</svelte:head>

<Navbar title={(channelVideos && channelVideos.ID) || data.id} />
{#if channelVideos}
    <ChannelCard channel={channelVideos} />
    <Grid>
        {#each channelVideos.Videos as v}
            <VideoCard video={v} needStyle showPoster />
        {/each}
    </Grid>
{/if}
