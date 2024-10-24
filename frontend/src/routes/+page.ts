import { Videos } from "$lib/api/video.js";

export async function load() {
    const videos = await Videos();

    return {
        videos,
    };
}
