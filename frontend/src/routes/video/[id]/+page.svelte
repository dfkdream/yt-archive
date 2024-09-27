<script lang="ts">
    import { type Video, VideoInfo } from '$lib/api/video';
    import DashAudio from '$lib/dash_audio.svelte';
    import DashVideo from '$lib/dash_video.svelte';
    import Tabbar from '$lib/tabbar.svelte';
    import VideoCard from '$lib/video_card.svelte';
    import type dashjs from 'dashjs';
    import { List, ListItem, Block, Toggle } from 'konsta/svelte';
    import { onMount } from 'svelte';

    let mediaClass = "m-auto w-full sticky top-0 z-50 max-h-[60vh] bg-black";

    export let data;

    let video: Video;
    onMount(async ()=>{
        video = await VideoInfo(data.id);
    })

    let manifest = `/api/videos/${data.id}/${data.id}.mpd`
    let poster = `/api/videos/${data.id}/${data.id}.webp`

    let bufferLength = 0;
    let videoQuality = 0;
    let videoBitrateList: dashjs.BitrateInfo[] | null = null;

    let loop = false;
    let radioMode = false;

    let bitrateString = "N/A";
    function getBitrateString(list: dashjs.BitrateInfo[] | null, quality: number): string{
        if (!list){
            return "N/A";
        }

        if (list.length-1 < quality){
            return "N/A";
        }

        let info = list[quality];
        return `${info.width}x${info.height}`;
    }

    $:bitrateString = getBitrateString(videoBitrateList, videoQuality);
</script>

<svelte:head>
    {#if video}
    <title>{video.Title} - yt-archive</title>
    {:else}
    <title>{data.id} - yt-archive</title>
    {/if}
</svelte:head>

{#if radioMode}
<DashAudio 
    {manifest} 
    {poster} 
    controls 
    {loop} 
    bind:bufferLength 
    class={mediaClass}
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
/>
{/if}

{#if video}
<VideoCard video={video} showChannel fullTitle />
<Block strong inset>
    <pre class="font-sans overflow-x-scroll">{video.Description.trim()}</pre>
</Block>
{/if}

<List strong inset>
    <ListItem title="Loop">
        <Toggle
            slot="after"
            bind:checked={loop}
        />
    </ListItem>
    <ListItem title="Radio Mode">
        <Toggle
            slot="after"
            bind:checked={radioMode}
        />
    </ListItem>
</List>

<List strong inset>
    {#if !radioMode}
    <ListItem title="Quality" after={bitrateString} />
    {/if}
    <ListItem title="Buffer length" after={bufferLength} />
</List>

<Tabbar />