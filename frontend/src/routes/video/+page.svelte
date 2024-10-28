<script lang="ts">
    import { goto } from "$app/navigation";
    import type { PlaylistVideos } from "$lib/api/playlist.js";
    import type { Video } from "$lib/api/video.js";
    import DashAudio from "$lib/components/dash_audio.svelte";
    import DashVideo from "$lib/components/dash_video.svelte";
    import LinkableText from "$lib/components/linkable_text.svelte";
    import Navbar from "$lib/components/navbar.svelte";
    import VideoCard from "$lib/components/video_card.svelte";
    import type dashjs from "dashjs";
    import { List, ListItem, Block, Toggle, BlockTitle } from "konsta/svelte";

    let mediaClass = "m-auto w-full sticky top-0 z-50 max-h-[60vh] bg-black";

    export let data;

    $: video = data.video;
    $: manifest = `/api/videos/${data.id}/${data.id}.mpd`;
    $: poster = `/api/thumbnails/${data.id}.webp`;

    let bufferLength = 0;
    let videoQuality = 0;
    let videoBitrateList: dashjs.BitrateInfo[] | null = null;

    let startTime = data.startTime;
    let currentTime = startTime;

    let loop = false;
    let radioMode = false;
    let autoplay = false;
    let loopPlaylist = false;

    let isPlaying = false;

    let videoBefore: Video | null = null;
    let videoAfter: Video | null = null;

    function getAdjacentVideos(
        video: Video,
        playlist: PlaylistVideos | null,
        loopPlaylist: boolean,
    ) {
        if (!playlist) return;

        const len = playlist.Videos.length;

        videoBefore = null;
        videoAfter = null;

        playlist.Videos.forEach((v, i, a) => {
            if (v.ID === video.ID) {
                if (i != 0 || loopPlaylist) {
                    videoBefore = a[(len + i - 1) % len];
                }

                if (i != len - 1 || loopPlaylist) {
                    videoAfter = a[(i + 1) % len];
                }
            }
        });
    }

    $: getAdjacentVideos(data.video, data.playlist, loopPlaylist);

    function nextVideo() {
        if (!autoplay) return;
        if (!videoAfter) return;
        if (!data.playlist) return;
        isPlaying = true;
        goto(`?id=${videoAfter.ID}&list=${data.playlist.ID}`);
    }

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
        onPlaybackEnded={nextVideo}
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
        onPlaybackEnded={nextVideo}
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
    {#if data.playlist}
        <ListItem title="Autoplay">
            <Toggle slot="after" bind:checked={autoplay} />
        </ListItem>
        <ListItem title="Loop Playlist">
            <Toggle slot="after" bind:checked={loopPlaylist} />
        </ListItem>
    {/if}
</List>

{#if data.playlist && videoAfter}
    <BlockTitle>Next</BlockTitle>
    <VideoCard video={videoAfter} listID={data.playlist.ID} showChannel />
{/if}

{#if data.playlist && videoBefore}
    <BlockTitle>Previous</BlockTitle>
    <VideoCard video={videoBefore} listID={data.playlist.ID} showChannel />
{/if}

<List strong inset>
    {#if !radioMode}
        <ListItem title="Quality" after={bitrateString} />
    {/if}
    <ListItem title="Buffer length" after={bufferLength} />
</List>
