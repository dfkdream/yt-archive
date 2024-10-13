<script lang="ts">
    import { PlaylistVideos } from "$lib/api/playlist.js";
    import Navbar from "$lib/navbar.svelte";
    import { Block } from "konsta/svelte";
    import { onMount } from "svelte";
    import VideoCard from "$lib/video_card.svelte";
    import Grid from "$lib/grid.svelte";

    export let data;

    export let playlistVideos: PlaylistVideos;

    onMount(async () => {
        playlistVideos = await PlaylistVideos(data.id);
    });
</script>

<svelte:head>
    {#if playlistVideos}
        <title>{playlistVideos.Title} - yt-archive</title>
    {:else}
        <title>{data.id} - yt-archive</title>
    {/if}
</svelte:head>

<Navbar title={(playlistVideos && playlistVideos.Title) || data.id} />

{#if playlistVideos}
    <Block strong inset class="!my-4">
        {#if playlistVideos.Description}
            <pre
                class="font-sans overflow-y-scroll">{playlistVideos.Description.trim()}</pre>
            <br />
        {/if}
        <a href={`/channels/${playlistVideos.Owner}`} class="flex items-center">
            <img
                src={`/api/channels/${playlistVideos.Owner}/${playlistVideos.OwnerThumbnail}`}
                alt={playlistVideos.Owner}
                width="45px"
                class="rounded-full"
            />
            <span class="text-sm mx-2">
                {playlistVideos.Owner}
                <br />
                {playlistVideos.Timestamp.toLocaleString([], {
                    dateStyle: "medium",
                    timeStyle: "medium",
                })}
            </span>
        </a>
    </Block>

    <Grid>
        {#each playlistVideos.Videos as v}
            <VideoCard video={v} showChannel showPoster needStyle />
        {/each}
    </Grid>
{/if}
