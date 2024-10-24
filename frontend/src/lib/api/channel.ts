import { mapTimestamp, type Video } from "./video";

export interface Channel {
    ID: string;
    Title: string;
    Description: string;
    Thumbnail: string;
}

export async function Channels(f = fetch) {
    let resp = await f("/api/channels");
    let json: Channel[] = await resp.json();
    return json;
}

export interface ChannelVideos extends Channel {
    Videos: Video[];
}

export async function ChannelVideos(id: string, f = fetch) {
    let resp = await f(`/api/channels/${id}`);
    let json: ChannelVideos = await resp.json();
    json.Videos = mapTimestamp(json.Videos);
    return json;
}
