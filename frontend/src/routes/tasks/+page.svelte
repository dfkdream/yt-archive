<script lang="ts">
    import {
        checkID,
        SubmitTask,
        Tasks,
        type TaskRequest,
    } from "$lib/api/tasks";
    import TaskCard from "$lib/task_card.svelte";
    import { BlockTitle, List, ListButton, ListInput } from "konsta/svelte";
    import Navbar from "$lib/navbar.svelte";

    export let data;

    let taskRequest: TaskRequest = {
        Type: 0,
        ID: "",
    };

    let validationError = "";

    function validateID() {
        [validationError, taskRequest.ID, taskRequest.Type] = checkID(
            taskRequest.ID,
        );
    }

    let tasks = data.tasks;

    $: queuedTasks = tasks.filter((t) => t.Status == 0);
    $: completedTasks = tasks.filter((t) => t.Status > 1);

    async function submitTask() {
        if (validationError) return;

        if (!taskRequest.ID) return;

        let id = await SubmitTask(taskRequest);

        taskRequest.Type = 0;
        taskRequest.ID = "";

        tasks = await Tasks();

        location.hash = id;
    }
</script>

<svelte:head>
    <title>Tasks - yt-archive</title>
</svelte:head>

<Navbar title="Tasks" location="tasks" />

<BlockTitle>New Task</BlockTitle>
<List strong inset>
    <ListInput
        type="text"
        bind:value={taskRequest.ID}
        error={validationError}
        onChange={validateID}
        placeholder="URL/ID"
    />
    <ListInput type="select" bind:value={taskRequest.Type} dropdown>
        <option value={0}>Video</option>
        <option value={1}>Playlist</option>
    </ListInput>
    <ListButton onClick={submitTask}>Submit Task</ListButton>
</List>

<BlockTitle>Running Tasks</BlockTitle>
{#each tasks.filter((t) => t.Status == 1) as t}
    <TaskCard task={t} />
{/each}

<BlockTitle>Queued Tasks ({queuedTasks.length})</BlockTitle>
{#each queuedTasks as t}
    <TaskCard task={t} />
{/each}

<BlockTitle>Completed Tasks ({completedTasks.length})</BlockTitle>
{#each completedTasks as t}
    <TaskCard task={t} />
{/each}
