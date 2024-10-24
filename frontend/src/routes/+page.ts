import { Videos } from "$lib/api/video.js";

export async function load({ fetch }) {
    const videos = await Videos(fetch);

    return {
        videos,
    };
}
