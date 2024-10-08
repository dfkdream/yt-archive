<script lang="ts">
    import type dashjs from "dashjs";
    import { onMount } from "svelte";

    let audioElement: HTMLAudioElement;

    export let manifest: string;

    export let poster: string | undefined | null = undefined;
    export let controls: boolean | undefined | null = undefined;
    export let playsinline: boolean | undefined | null = undefined;
    export let loop: boolean | undefined | null = undefined;
    export let autoplay: boolean = false;

    export let bufferLength: number = 0;

    export let startTime: number = 0;
    export let currentTime: number = 0;

    export let isPlaying: boolean = false;

    let player: dashjs.MediaPlayerClass | null = null;

    onMount(async () => {
        const { MediaPlayer } = await import("dashjs");
        player = MediaPlayer().create();
        player.initialize(audioElement, manifest, autoplay);

        player.on("streamInitialized", () => {
            player?.seek(startTime);
        });

        player.on("metricsChanged", () => {
            bufferLength = player?.getBufferLength("audio") || 0;
        });

        player.on("playbackPlaying", () => {
            isPlaying = true;
        });

        player.on("playbackPaused", () => {
            isPlaying = false;
        });
    });
</script>

<div class="relative" {...$$restProps}>
    <img src={poster} alt="poster" class={"object-contain " + $$props.class} />
    <audio
        bind:this={audioElement}
        {controls}
        {playsinline}
        loop={loop ? loop : undefined}
        bind:currentTime
        class="w-full absolute bottom-0 z-50"
    >
        <slot />
    </audio>
</div>
