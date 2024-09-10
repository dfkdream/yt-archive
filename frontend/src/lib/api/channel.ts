import { mapTimestamp, type Video } from "./video";

export interface Channel {
    ID: string
    Title: string
    Description: string
    Thumbnail: string
}

export async function Channels(){
    let resp = await fetch("/api/channels");
    let json: Channel[] = await resp.json();
    return json
}

export interface ChannelVideos extends Channel{
    Videos: Video[]
}


export async function ChannelVideos(id: string){
    let resp = await fetch(`/api/channels/${id}`);
    let json: ChannelVideos = await resp.json();
    json.Videos = mapTimestamp(json.Videos)
    return json
}