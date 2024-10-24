export interface Video {
    ID: string;
    Title: string;
    Description: string;
    Timestamp: Date;
    Duration: string;
    Owner: string;
    Thumbnail: string;
    OwnerThumbnail: string;
}

export async function Videos(f = fetch) {
    let resp = await f("/api/videos");
    let json: Video[] = await resp.json();
    json = json.map((v) => {
        v.Timestamp = new Date(v.Timestamp);
        return v;
    });
    return json;
}

export async function VideoInfo(id: string, f = fetch) {
    let resp = await f(`/api/videos/${id}`);
    let json: Video = await resp.json();
    json.Timestamp = new Date(json.Timestamp);
    return json;
}

export function mapTimestamp<T extends Video>(v: T[]) {
    return v.map((v) => {
        v.Timestamp = new Date(v.Timestamp);
        return v;
    });
}
