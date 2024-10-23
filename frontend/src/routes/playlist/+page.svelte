<script lang="ts">
    import {
        PlaylistVideos,
        UpdateIndex,
        type IndexedVideo,
    } from "$lib/api/playlist.js";
    import Navbar from "$lib/navbar.svelte";
    import {
        Block,
        List,
        ListItem,
        ListInput,
        Dialog,
        DialogButton,
    } from "konsta/svelte";
    import { onMount } from "svelte";
    import VideoCard from "$lib/video_card.svelte";
    import Grid from "$lib/grid.svelte";

    export let data;

    export let playlistVideos: PlaylistVideos | null = null;

    onMount(async () => {
        playlistVideos = await PlaylistVideos(data.id);
    });

    let dialogOpened = false;
    let dialogVideo: IndexedVideo | null = null;
    let newIndex: number | undefined = undefined;
</script>

<svelte:head>
    {#if playlistVideos}
        <title>{playlistVideos.Title} - yt-archive</title>
    {:else}
        <title>{data.id} - yt-archive</title>
    {/if}
</svelte:head>

<Navbar title={(playlistVideos && playlistVideos.Title) || data.id} />

<Dialog
    opened={dialogOpened}
    onBackdropClick={() => {
        dialogOpened = false;
    }}
>
    <span class="text-sm" slot="title">{dialogVideo?.Title}</span>
    <List nested class="-mx-4">
        <ListItem title="Index" after={dialogVideo?.Index.toString()} />
        <ListInput
            type="number"
            placeholder="New Index"
            bind:value={newIndex}
        />
    </List>
    <DialogButton
        slot="buttons"
        onClick={async () => {
            if (
                playlistVideos != null &&
                dialogVideo != null &&
                newIndex != undefined
            ) {
                await UpdateIndex(playlistVideos.ID, dialogVideo.ID, newIndex);
                playlistVideos = await PlaylistVideos(data.id);
            }
            dialogOpened = false;
        }}>Update Index</DialogButton
    >
</Dialog>

{#if playlistVideos}
    <Block strong inset class="!my-4">
        {#if playlistVideos.Description}
            <pre
                class="font-sans overflow-y-scroll">{playlistVideos.Description.trim()}</pre>
            <br />
        {/if}
        <a
            href={`/channel?id=${playlistVideos.Owner}`}
            class="flex items-center"
        >
            <img
                src={`/api/thumbnails/${playlistVideos.OwnerThumbnail}`}
                alt={playlistVideos.Owner}
                width="45px"
                class="rounded-full"
            />
            <span class="text-sm mx-2">
                {playlistVideos.Owner}
                <br />
                {playlistVideos.Timestamp.toLocaleString([], {
                    dateStyle: "medium",
                    timeStyle: "medium",
                })}
            </span>
        </a>
    </Block>

    <Grid>
        {#each playlistVideos.Videos as v}
            <VideoCard video={v} showChannel showPoster needStyle>
                <button
                    class="bg-black text-white h-8 w-8 p-1 rounded opacity-50 hover:opacity-70 absolute top-2 right-2 z-20"
                    title="Update index"
                    on:click={() => {
                        dialogVideo = v;
                        newIndex = undefined;
                        dialogOpened = true;
                    }}
                >
                    ⋮
                </button>
            </VideoCard>
        {/each}
    </Grid>
{/if}
