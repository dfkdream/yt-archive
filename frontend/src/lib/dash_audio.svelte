<script lang="ts">
    import type dashjs from "dashjs";
    import { onMount } from "svelte";

    let audioElement: HTMLAudioElement;

    export let manifest: string;

    export let poster: string | undefined | null = undefined;
    export let controls: boolean | undefined | null = undefined;
    export let playsinline: boolean | undefined | null = undefined;
    export let loop: boolean | undefined | null = undefined;

    export let bufferLength: number = 0;

    let player: dashjs.MediaPlayerClass | null = null;

    onMount(async () => {
        const { MediaPlayer } = await import("dashjs");
        player = MediaPlayer().create();
        player.initialize(audioElement, manifest, false);

        player.on("metricsChanged", () => {
            if (!player) return;

            try {
                bufferLength = player.getBufferLength("audio");
            } catch {}
        });
    });
</script>

<div {...$$restProps}>
    <img src={poster} alt="poster" />
    <audio
        bind:this={audioElement}
        {controls}
        {playsinline}
        loop={loop ? loop : undefined}
        class="w-full absolute bottom-0"
    >
        <slot />
    </audio>
</div>
