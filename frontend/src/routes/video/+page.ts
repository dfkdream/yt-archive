import { error } from "@sveltejs/kit";

export function load({ url }) {
    const id = url.searchParams.get("id");
    if (!id) {
        error(400, "required paramter id not provided");
    }

    return {
        id,
        startTime: parseInt(url.searchParams.get("t") || "0"),
    };
}
