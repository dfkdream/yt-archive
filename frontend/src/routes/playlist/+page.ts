import { PlaylistVideos } from "$lib/api/playlist.js";
import { isValidPlaylistID } from "$lib/api/validation.js";
import { error } from "@sveltejs/kit";

export async function load({ url, fetch }) {
    const id = url.searchParams.get("id");
    if (!id) {
        error(400, "required paramter id not provided");
    }

    if (!isValidPlaylistID(id)) {
        error(400, "invalid parameter id");
    }

    const playlistVideos = await PlaylistVideos(id, fetch);

    return {
        id,
        playlistVideos,
    };
}
