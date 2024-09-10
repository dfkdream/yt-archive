<script lang="ts">
    import { Block } from "konsta/svelte";
    import type { Channel } from "./api/channel";

    export let channel: Channel
</script>

<Block strong inset class="flex items-center !my-4">
    <img src={`/api/channels/${channel.ID}/${channel.Thumbnail}`} alt={channel.Title} width="45px" class="rounded-full"/>
    <span class="mx-2">
        <h2 class="text-lg font-bold">{channel.Title}</h2>
        <pre class="font-sans overflow-y-scroll">{channel.Description.trim()}</pre>
    </span>
</Block>