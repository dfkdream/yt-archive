<script lang="ts">
    import { Badge, Block } from "konsta/svelte";
    import type { Task } from "./api/tasks";

    export let task: Task;
</script>

<Block strong inset id={task.ID}>
    <div class="flex justify-between mb-1">
        <h1 class="font-bold truncate">{task.ID}</h1>
        <Badge>{task.Priority}</Badge>
    </div>
    <h2 class="my-1">{task.Type}</h2>
    <p class="my-1">{task.Description}</p>
    <details class="mt-1">
        <summary>Show Payload</summary>
        <pre class="overflow-x-scroll">{task.Payload}</pre>
    </details>
</Block>