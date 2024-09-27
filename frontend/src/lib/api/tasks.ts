export const typeVideo = 0;
export const typePlaylist = 1;

export type TaskRequest = {
    Type: 0 | 1;
    ID: string;
};

export async function SubmitTask(t: TaskRequest) {
    let resp = await fetch("/api/tasks", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(t),
    });

    let taskId: string = await resp.json();
    return taskId;
}

export function checkID(id: string): [string, string, 0 | 1] {
    try {
        let url = new URL(id);

        if (
            url.host == "www.youtube.com" ||
            url.host == "m.youtube.com" ||
            url.host == "youtube.com"
        ) {
            if (url.pathname.startsWith("/shorts")) {
                id = url.pathname.slice(8);
                console.log(id);
            }

            if (url.pathname.startsWith("/live")) {
                id = url.pathname.slice(6);
            }

            if (url.pathname.startsWith("/watch")) {
                let list = url.searchParams.get("list");

                if (list) {
                    id = list;
                } else {
                    id = url.searchParams.get("v") ?? "";
                }
            }

            if (url.pathname.startsWith("/playlist")) {
                id = url.searchParams.get("list") ?? "";
            }
        }

        if (url.host == "youtu.be") {
            id = url.pathname.slice(1);
        }
    } catch {}

    let type: 0 | 1;

    if (/^[A-Za-z0-9_-]{11}$/.test(id)) {
        type = typeVideo;
    } else if (/^PL[A-Za-z0-9_-]{32}$/.test(id)) {
        type = typePlaylist;
    } else {
        return ["Invalid ID/URL", id, 0];
    }

    return ["", id, type];
}

export type Task = {
    ID: string;
    Status: number;
    Priority: number;
    Type: string;
    Description: string;
    Payload: string;
};

export async function Tasks() {
    let resp = await fetch("/api/tasks");
    let task: Task[] = await resp.json();
    return task;
}
