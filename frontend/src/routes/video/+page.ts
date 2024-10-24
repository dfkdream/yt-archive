import { isValidVideoID } from "$lib/api/validation.js";
import { VideoInfo } from "$lib/api/video.js";
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

    return {
        id,
        video,
        startTime: parseInt(url.searchParams.get("t") || "0"),
    };
}
