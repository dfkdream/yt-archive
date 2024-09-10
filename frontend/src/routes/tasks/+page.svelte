<script lang="ts">
    import { checkID, SubmitTask, Tasks, type Task, type TaskRequest } from "$lib/api/tasks";
    import Tabbar from "$lib/tabbar.svelte";
    import TaskCard from "$lib/task_card.svelte";
    import { BlockTitle, List, ListButton, ListInput, Navbar } from "konsta/svelte";
    import { onMount } from "svelte";

    let taskRequest: TaskRequest = {
        Type: 0,
        ID: ""
    };

    let validationError = "";

    function validateID(){
        [validationError, taskRequest.ID, taskRequest.Type] = checkID(taskRequest.ID)
    }

    let tasks: Task[] = [];

    async function submitTask(){
        if (validationError) return;

        if (!taskRequest.ID) return;

        let id = await SubmitTask(taskRequest);

        taskRequest.Type = 0;
        taskRequest.ID = "";

        tasks = await Tasks();

        location.hash = id;
    }

    onMount(async ()=>{
        tasks = await Tasks();
    })

</script>

<svelte:head>
    <title>Tasks - yt-archive</title>
</svelte:head>

<Navbar medium transparent title="Tasks" />

<BlockTitle>New Task</BlockTitle>
<List strong inset>
    <ListInput type="text" bind:value={taskRequest.ID} error={validationError} onChange={validateID} placeholder="URL/ID" />
    <ListInput type="select" bind:value={taskRequest.Type} dropdown>
        <option value={0}>Video</option>
        <option value={1}>Playlist</option>
    </ListInput>
    <ListButton onClick={submitTask}>Submit Task</ListButton>
</List>

<BlockTitle>Running Tasks</BlockTitle>
{#each tasks.filter(t=>t.Status==1) as t}
<TaskCard task={t} />
{/each}

<BlockTitle>Queued Tasks</BlockTitle>
{#each tasks.filter(t=>t.Status==0) as t}
<TaskCard task={t} />
{/each}

<BlockTitle>Completed Tasks</BlockTitle>
{#each tasks.filter(t=>t.Status>1) as t}
<TaskCard task={t} />
{/each}

<Tabbar location="tasks" />