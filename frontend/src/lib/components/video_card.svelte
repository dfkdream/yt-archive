<script lang="ts">
    import { Block } from "konsta/svelte";
    import type { Video } from "../api/video";

    export let showPoster = false;
    export let showChannel = false;
    export let fullTitle = false;
    export let needStyle = false;
    export let video: Video;
    export let listID: string | null = null;

    $: href = `/video?id=${video.ID}` + (listID ? `&list=${listID}` : "");
</script>

<Block strong inset class={needStyle ? "!my-1 md:!my-4" : ""}>
    <slot />
    {#if showPoster}
        <a {href}>
            <div class="relative">
                <img
                    class="mb-2 aspect-video bg-gray-50"
                    loading="lazy"
                    src={`/api/thumbnails/${video.Thumbnail}`}
                    alt={video.Title}
                />
                <span
                    class="bg-black text-white p-1 rounded opacity-70 absolute bottom-1 right-1"
                    >{video.Duration}</span
                >
            </div>
        </a>
    {/if}

    <a {href}>
        <h2
            class={"mb-2 text-sm font-bold" +
                (fullTitle ? "" : " whitespace-nowrap break-all truncate")}
        >
            {video.Title}
        </h2>
    </a>

    {#if showChannel}
        <a href={`/channel?id=${video.Owner.ID}`} class="flex items-center">
            <img
                loading="lazy"
                src={`/api/thumbnails/${video.Owner.Thumbnail}`}
                alt={video.Owner.Title}
                width="45px"
                class="rounded-full"
            />
            <span class="text-sm mx-2">
                {video.Owner.ID}
                <br />
                {video.Timestamp.toLocaleString([], {
                    dateStyle: "medium",
                    timeStyle: "medium",
                })}
            </span>
        </a>
    {:else}
        <div class="text-sm">
            {video.Timestamp.toLocaleString([], {
                dateStyle: "medium",
                timeStyle: "medium",
            })}
        </div>
    {/if}
</Block>
