<script lang="ts">
    import { ChannelVideos } from "$lib/api/channel.js";
    import ChannelCard from "$lib/channel_card.svelte";
    import Grid from "$lib/grid.svelte";
    import Navbar from "$lib/navbar.svelte";
    import VideoCard from "$lib/video_card.svelte";
    import { onMount } from "svelte";

    export let data;

    export let channelVideos: ChannelVideos;

    onMount(async () => {
        channelVideos = await ChannelVideos(data.id);
    });
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
