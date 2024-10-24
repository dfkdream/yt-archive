import { Playlists } from "$lib/api/playlist";

export async function load({ fetch }) {
    const playlists = await Playlists(fetch);

    return {
        playlists,
    };
}
