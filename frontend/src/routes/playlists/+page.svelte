<script lang="ts">
    import { Playlists, type Playlist } from "$lib/api/playlist";
    import Tabbar from "$lib/tabbar.svelte";
    import { Block, Navbar } from "konsta/svelte";
    import { onMount } from "svelte";

    let playlists: Playlist[] = [];

    onMount(async () => {
        playlists = await Playlists();
    });
</script>

<svelte:head>
    <title>Playlists - yt-archive</title>
</svelte:head>

<Navbar medium transparent title="Playlists" />

<div class="grid gap-1 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
    {#each playlists as p}
        <a href={`/playlists/${p.ID}/`}>
            <Block strong inset class="!my-1 md:!my-4">
                <img
                    class="mb-2"
                    src={`/api/videos/${p.ThumbnailVideo}/${p.Thumbnail}`}
                    alt={p.Title}
                />

                <h2
                    class={"mb-2 text-sm font-bold whitespace-nowrap break-all truncate"}
                >
                    {p.Title}
                </h2>

                <a href={`/channels/${p.Owner}`} class="flex items-center">
                    <img
                        src={`/api/channels/${p.Owner}/${p.OwnerThumbnail}`}
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
</div>

<Tabbar location="playlists"></Tabbar>
