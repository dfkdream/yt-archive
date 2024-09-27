import { mapTimestamp, type Video } from "./video";

export interface Playlist {
    ID: string;
    Title: string;
    Description: string;
    Timestamp: Date;
    Owner: string;
    OwnerThumbnail: string;
    ThumbnailVideo: string;
    Thumbnail: string;
}

export async function Playlists() {
    let resp = await fetch("/api/playlists");
    let json: Playlist[] = await resp.json();
    return json.map((v) => {
        v.Timestamp = new Date(v.Timestamp);
        return v;
    });
}

export interface PlaylistVideos extends Playlist {
    Videos: Video[];
}

export async function PlaylistVideos(id: string) {
    let resp = await fetch(`/api/playlists/${id}`);
    let json: PlaylistVideos = await resp.json();
    json.Timestamp = new Date(json.Timestamp);
    json.Videos = mapTimestamp(json.Videos);
    return json;
}
