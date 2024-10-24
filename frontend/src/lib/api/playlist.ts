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

export async function Playlists(f = fetch) {
    let resp = await f("/api/playlists");
    let json: Playlist[] = await resp.json();
    return json.map((v) => {
        v.Timestamp = new Date(v.Timestamp);
        return v;
    });
}

export interface IndexedVideo extends Video {
    Index: number;
}

export interface PlaylistVideos extends Playlist {
    Videos: IndexedVideo[];
}

export async function PlaylistVideos(id: string, f = fetch) {
    let resp = await f(`/api/playlists/${id}`);
    let json: PlaylistVideos = await resp.json();
    json.Timestamp = new Date(json.Timestamp);
    json.Videos = mapTimestamp(json.Videos);
    return json;
}

export async function UpdateIndex(
    pid: string,
    vid: string,
    i: number,
    f = fetch,
) {
    let resp = await f(`/api/playlists/${pid}/video/${vid}/index`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: i.toString(),
    });

    let newIndex: number = await resp.json();
    return newIndex;
}
