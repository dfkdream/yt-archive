<script lang="ts">
    import { Block } from "konsta/svelte";
    import type { Video } from "./api/video";

    export let showPoster = false;
    export let showChannel = false;
    export let fullTitle = false;
    export let needStyle = false;
    export let video: Video;
</script>

<Block strong inset class={needStyle ? "!my-1 md:!my-4" : ""}>
    {#if showPoster}
        <div class="relative">
            <img
                class="mb-2"
                src={`/api/videos/${video.ID}/${video.Thumbnail}`}
                alt={video.Title}
            />
            <span
                class="bg-black text-white p-1 rounded opacity-70 absolute bottom-1 right-1"
                >{video.Duration}</span
            >
        </div>
    {/if}

    <h2
        class={"mb-2 text-sm font-bold" +
            (fullTitle ? "" : " whitespace-nowrap break-all truncate")}
    >
        {video.Title}
    </h2>

    {#if showChannel}
        <a href={`/channels/${video.Owner}`} class="flex items-center">
            <img
                src={`/api/channels/${video.Owner}/${video.OwnerThumbnail}`}
                alt={video.Owner}
                width="45px"
                class="rounded-full"
            />
            <span class="text-sm mx-2">
                {video.Owner}
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
