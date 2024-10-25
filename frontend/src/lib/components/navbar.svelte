<script lang="ts">
    import { goto } from "$app/navigation";
    import { Navbar, Segmented, SegmentedButton } from "konsta/svelte";

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
            <SegmentedButton
                strong
                active={l.isActive}
                onClick={() => {
                    goto(l.location);
                }}>{l.name}</SegmentedButton
            >
        {/each}
    </Segmented>
</Navbar>
