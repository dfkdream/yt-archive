<script lang="ts">
    import Navbar from "$lib/components/navbar.svelte";
    import ChannelCard from "$lib/components/channel_card.svelte";

    export let data;

    let channels = data.channels;
</script>

<svelte:head>
    <title>Channels - yt-archive</title>
</svelte:head>

<Navbar title="Channels" location="channels" />

{#each channels as c}
    <a href={`/channel?id=${c.ID}`}>
        <ChannelCard channel={c} />
    </a>
{/each}
