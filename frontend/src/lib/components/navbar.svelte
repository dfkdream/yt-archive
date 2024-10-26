<script lang="ts">
    import { Navbar, Segmented, Button } from "konsta/svelte";

    export let location: "home" | "channels" | "playlists" | "tasks" | null =
        null;

    export let title = "yt-archive";

    export let small = false;

    let locations = [
        { name: "Home", isActive: location == "home", location: "/" },
        {
            name: "Channels",
            isActive: location == "channels",
            location: "/channels",
        },
        {
            name: "Playlists",
            isActive: location == "playlists",
            location: "/playlists",
        },
        { name: "Tasks", isActive: location == "tasks", location: "/tasks" },
    ];
</script>

<Navbar medium={!small} transparent={!small} {title}>
    <Segmented slot="subnavbar" strong>
        {#each locations as l}
            <Button
                segmentedStrong
                segmentedActive={l.isActive}
                href={l.location}
            >
                {l.name}
            </Button>
        {/each}
    </Segmented>
</Navbar>
