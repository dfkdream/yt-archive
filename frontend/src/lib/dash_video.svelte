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

    export let videoQuality: number = 0;
    export let bufferLength: number = 0;
    export let videoBitrateList: dashjs.BitrateInfo[] | null = null;

    let player: dashjs.MediaPlayerClass | null = null;

    onMount(async () => {
        const { MediaPlayer } = await import("dashjs");
        player = MediaPlayer().create();
        player.initialize(videoElement, manifest, false);

        player.on("metricsChanged", () => {
            if (!player) return;

            try {
                videoQuality = player.getQualityFor("video");
                bufferLength = player.getBufferLength("video");
                videoBitrateList = player.getBitrateInfoListFor("video");
            } catch {}
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
    {...$$restProps}
>
    <slot />
</video>
