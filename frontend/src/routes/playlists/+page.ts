import { Playlists } from "$lib/api/playlist";

export async function load() {
    const playlists = await Playlists();

    return {
        playlists,
    };
}
