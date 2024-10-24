import { Tasks } from "$lib/api/tasks";

export async function load() {
    const tasks = await Tasks();

    return {
        tasks,
    };
}
