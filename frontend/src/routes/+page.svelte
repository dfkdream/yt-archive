<script lang="ts">
    import { Navbar } from 'konsta/svelte';
    import Tabbar from '$lib/tabbar.svelte';
    import { type Video, Videos } from '$lib/api/video';
    import { onMount } from 'svelte';
    import VideoGrid from '$lib/video_grid.svelte';

    let videos: Video[] = [];

    onMount(async ()=>{
        videos = await Videos();
    })
</script>

<svelte:head>
    <title>Home - yt-archive</title>
</svelte:head>

<Navbar medium transparent title="yt-archive" />
<VideoGrid videos={videos} showChannel/>
<Tabbar location="home"/>
