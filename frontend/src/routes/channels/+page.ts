import { Channels } from "$lib/api/channel";

export async function load({ fetch }) {
    const channels = await Channels(fetch);

    return {
        channels,
    };
}
