<script lang="ts">
    import type dashjs from "dashjs";
    import { onMount } from "svelte";

    let videoElement: HTMLVideoElement;

    export let manifest: string;

    export let height: number | string | undefined | null = undefined;
    export let width: number | string | undefined | null = undefined;
    export let poster: string | undefined | null = undefined;
    export let controls: boolean | undefined | null = undefined;
    export let playsinline: boolean | undefined | null = undefined;
    export let loop: boolean | undefined | null = undefined;
    export let autoplay: boolean = false;

    export let videoQuality: number = 0;
    export let bufferLength: number = 0;
    export let videoBitrateList: dashjs.BitrateInfo[] | null = null;

    export let startTime: number = 0;
    export let currentTime: number = 0;

    export let isPlaying: boolean = false;

    let player: dashjs.MediaPlayerClass | null = null;

    onMount(async () => {
        const { MediaPlayer } = await import("dashjs");
        player = MediaPlayer().create();
        player.initialize(videoElement, manifest, autoplay);

        player.on("streamInitialized", () => {
            videoBitrateList = player?.getBitrateInfoListFor("video") || null;
            player?.seek(startTime);
        });

        player.on("metricsChanged", () => {
            videoQuality = player?.getQualityFor("video") || 0;
            bufferLength = player?.getBufferLength("video") || 0;
        });

        player.on("playbackPlaying", () => {
            isPlaying = true;
        });

        player.on("playbackPaused", () => {
            isPlaying = false;
        });
    });
</script>

<!-- svelte-ignore a11y-media-has-caption -->
<video
    bind:this={videoElement}
    {height}
    {width}
    {poster}
    {controls}
    {playsinline}
    loop={loop ? loop : undefined}
    bind:currentTime
    {...$$restProps}
>
    <slot />
</video>
