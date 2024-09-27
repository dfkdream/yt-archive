<script lang="ts">
    import Tabbar from "$lib/tabbar.svelte";
    import { onMount } from "svelte";
    import { type Channel, Channels } from "$lib/api/channel";
    import { Navbar } from "konsta/svelte";
    import ChannelCard from "$lib/channel_card.svelte";

    let channels: Channel[] = [];

    onMount(async () => {
        channels = await Channels();
    });
</script>

<svelte:head>
    <title>Channels - yt-archive</title>
</svelte:head>

<Navbar medium transparent title="Channels" />

{#each channels as c}
    <a href={`/channels/${c.ID}`}>
        <ChannelCard channel={c} />
    </a>
{/each}

<Tabbar location="channels" />
