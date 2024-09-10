<script lang="ts">
    import type dashjs from "dashjs";
    import { onMount } from "svelte";

    let videoElement: HTMLVideoElement;

    export let manifest: string;

	export let height: number | string | undefined | null = undefined;
	export let width: number | string | undefined | null = undefined;
	export let poster: string | undefined | null = undefined;
	export let controls: boolean | undefined | null = undefined;
    export let autoplay: boolean | undefined | null = undefined;
	export let playsinline: boolean | undefined | null = undefined;
	export let disablepictureinpicture: boolean | undefined | null = undefined;
	export let disableremoteplayback: boolean | undefined | null = undefined;

    export let pip = false;

    export let videoQuality: number = 0;
    export let bufferLength: number = 0;
    export let videoBitrateList: dashjs.BitrateInfo[] | null = null;

    let player: dashjs.MediaPlayerClass | null = null;

    onMount(async ()=>{
        const {MediaPlayer} = await import('dashjs');
        player = MediaPlayer().create();
        player.initialize(videoElement, manifest, autoplay?true:false);

        player.on("metricsChanged", ()=>{
            if (!player) return;

            try{
                videoQuality = player.getQualityFor("video")
                bufferLength = player.getBufferLength("video")
                videoBitrateList = player.getBitrateInfoListFor("video")
            }catch{
            }
        });
    });
</script>

<video bind:this={videoElement}
    {height} {width} {poster} {controls} {autoplay} {playsinline} 
    {disablepictureinpicture} {disableremoteplayback} {...$$restProps}>
    <slot />
</video>