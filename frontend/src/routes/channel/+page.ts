import { ChannelVideos } from "$lib/api/channel.js";
import { error } from "@sveltejs/kit";

export async function load({ url, fetch }) {
    const id = url.searchParams.get("id");
    if (!id) {
        error(400, "required paramter id not provided");
    }

    const channelVideos = await ChannelVideos(id, fetch);

    return {
        id,
        channelVideos,
    };
}
