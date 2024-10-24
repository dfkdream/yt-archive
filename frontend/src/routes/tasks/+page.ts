import { Tasks } from "$lib/api/tasks";

export async function load({ fetch }) {
    const tasks = await Tasks(fetch);

    return {
        tasks,
    };
}
