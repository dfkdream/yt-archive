import { VideoInfo } from "$lib/api/video.js";
import { error } from "@sveltejs/kit";

export async function load({ url }) {
    const id = url.searchParams.get("id");
    if (!id) {
        error(400, "required paramter id not provided");
    }

    const video = await VideoInfo(id);

    return {
        id,
        video,
        startTime: parseInt(url.searchParams.get("t") || "0"),
    };
}
