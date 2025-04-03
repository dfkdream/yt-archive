<script lang="ts">
    import { Block } from "konsta/svelte";
    import Navbar from "$lib/components/navbar.svelte";
    import Grid from "$lib/components/grid.svelte";

    export let data;

    let playlists = data.playlists;
</script>

<svelte:head>
    <title>Playlists - yt-archive</title>
</svelte:head>

<Navbar title="Playlists" location="playlists" />

<Grid>
    {#each playlists as p}
        <a href={`/playlist?id=${p.ID}`}>
            <Block strong inset class="!my-1 md:!my-4">
                <img
                    class="mb-2"
                    src={`/api/thumbnails/${p.Thumbnail}`}
                    alt={p.Title}
                />

                <h2
                    class={"mb-2 text-sm font-bold whitespace-nowrap break-all truncate"}
                >
                    {p.Title}
                </h2>

                <a href={`/channel?id=${p.Owner}`} class="flex items-center">
                    <img
                        loading="lazy"
                        src={`/api/thumbnails/${p.OwnerThumbnail}`}
                        alt={p.Owner}
                        width="45px"
                        class="rounded-full"
                    />
                    <span class="text-sm mx-2">
                        {p.Owner}
                        <br />
                        {p.Timestamp.toLocaleString([], {
                            dateStyle: "medium",
                            timeStyle: "medium",
                        })}
                    </span>
                </a>
            </Block>
        </a>
    {/each}
</Grid>
