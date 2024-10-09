<script lang="ts">
    import { onMount } from "svelte";
    import { type Channel, Channels } from "$lib/api/channel";
    import Navbar from "$lib/navbar.svelte";
    import ChannelCard from "$lib/channel_card.svelte";

    let channels: Channel[] = [];

    onMount(async () => {
        channels = await Channels();
    });
</script>

<svelte:head>
    <title>Channels - yt-archive</title>
</svelte:head>

<Navbar title="Channels" location="channels" />

{#each channels as c}
    <a href={`/channels/${c.ID}`}>
        <ChannelCard channel={c} />
    </a>
{/each}
