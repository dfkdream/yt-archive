import { PlaylistVideos } from "$lib/api/playlist.js";
import { isValidPlaylistID, isValidVideoID } from "$lib/api/validation.js";
import { VideoInfo, type Video } from "$lib/api/video.js";
import { error } from "@sveltejs/kit";

export async function load({ url, fetch }) {
    const id = url.searchParams.get("id");
    if (!id) {
        error(400, "required paramter id not provided");
    }

    if (!isValidVideoID(id)) {
        error(400, "invalid parameter id");
    }

    const video = await VideoInfo(id, fetch);

    const playlistID = url.searchParams.get("list") || "";
    let videoBefore: Video | null = null;
    let videoAfter: Video | null = null;

    if (isValidPlaylistID(playlistID)) {
        const playlist = await PlaylistVideos(playlistID, fetch);

        let index = -1;
        for (let i = 0; i < playlist.Videos.length; i++) {
            if (playlist.Videos[i].ID === id) {
                index = i;
                break;
            }
        }

        if (index != -1) {
            if (index != 0) {
                videoBefore = playlist.Videos[index - 1];
            }

            if (index != playlist.Videos.length - 1) {
                videoAfter = playlist.Videos[index + 1];
            }
        }
    }

    return {
        id,
        playlistID,
        video,
        videoBefore,
        videoAfter,
        startTime: parseInt(url.searchParams.get("t") || "0"),
    };
}
