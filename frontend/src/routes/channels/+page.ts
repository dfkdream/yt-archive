import { Channels } from "$lib/api/channel";

export async function load() {
    const channels = await Channels();

    return {
        channels,
    };
}
