<script lang="ts">
    import { type Video, VideoInfo } from "$lib/api/video";
    import DashAudio from "$lib/dash_audio.svelte";
    import DashVideo from "$lib/dash_video.svelte";
    import LinkableText from "$lib/linkable_text.svelte";
    import Navbar from "$lib/navbar.svelte";
    import VideoCard from "$lib/video_card.svelte";
    import type dashjs from "dashjs";
    import { List, ListItem, Block, Toggle } from "konsta/svelte";

    let mediaClass = "m-auto w-full sticky top-0 z-50 max-h-[60vh] bg-black";

    export let data;

    let video = data.video;

    let manifest = `/api/videos/${data.id}/${data.id}.mpd`;
    let poster = `/api/thumbnails/${data.id}.webp`;

    let bufferLength = 0;
    let videoQuality = 0;
    let videoBitrateList: dashjs.BitrateInfo[] | null = null;

    let startTime = data.startTime;
    let currentTime = startTime;

    let loop = false;
    let radioMode = false;

    let isPlaying = false;

    let bitrateString = "N/A";
    function getBitrateString(
        list: dashjs.BitrateInfo[] | null,
        quality: number,
    ): string {
        if (!list) {
            return "N/A";
        }

        if (list.length - 1 < quality) {
            return "N/A";
        }

        let info = list[quality];
        return `${info.width}x${info.height}`;
    }

    $: bitrateString = getBitrateString(videoBitrateList, videoQuality);
    $: startTime = data.startTime;
</script>

<svelte:head>
    <title>{video.Title} - yt-archive</title>
</svelte:head>

<Navbar small />

{#if radioMode}
    <DashAudio
        {manifest}
        {poster}
        controls
        {loop}
        bind:bufferLength
        class={mediaClass}
        {startTime}
        bind:currentTime
        bind:isPlaying
        autoplay={isPlaying}
    />
{:else}
    <DashVideo
        {manifest}
        {poster}
        controls
        playsinline
        {loop}
        class={mediaClass}
        bind:videoQuality
        bind:videoBitrateList
        bind:bufferLength
        {startTime}
        bind:currentTime
        bind:isPlaying
        autoplay={isPlaying}
    />
{/if}

<VideoCard {video} showChannel fullTitle />
<Block strong inset>
    <LinkableText
        class="overflow-x-scroll text-nowrap"
        text={video.Description.trim()}
        videoId={video.ID}
    />
</Block>

<List strong inset>
    <ListItem title="Loop">
        <Toggle slot="after" bind:checked={loop} />
    </ListItem>
    <ListItem title="Radio Mode">
        <Toggle
            slot="after"
            bind:checked={radioMode}
            onChange={() => {
                startTime = currentTime;
            }}
        ></Toggle>
    </ListItem>
</List>

<List strong inset>
    {#if !radioMode}
        <ListItem title="Quality" after={bitrateString} />
    {/if}
    <ListItem title="Buffer length" after={bufferLength} />
</List>
