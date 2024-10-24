import { PlaylistVideos } from "$lib/api/playlist.js";
import { error } from "@sveltejs/kit";

export async function load({ url }) {
    const id = url.searchParams.get("id");
    if (!id) {
        error(400, "required paramter id not provided");
    }

    const playlistVideos = await PlaylistVideos(id);

    return {
        id,
        playlistVideos,
    };
}
