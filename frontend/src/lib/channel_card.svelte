<script lang="ts">
    import { Block } from "konsta/svelte";
    import type { Channel } from "./api/channel";

    export let channel: Channel;
</script>

<Block strong inset class="flex items-center !my-4">
    <img
        src={`/api/thumbnails/${channel.Thumbnail}`}
        alt={channel.Title}
        width="45px"
        class="rounded-full"
    />
    <div class="mx-2 overflow-hidden">
        <h2 class="text-lg font-bold">{channel.Title}</h2>
        <pre
            class="font-sans overflow-x-scroll">{channel.Description.trim()}</pre>
    </div>
</Block>
