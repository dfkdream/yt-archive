<script lang="ts">
    import { Tabbar, TabbarLink } from "konsta/svelte";

    export let location: "home" | "channels" | "playlists" | "tasks" | null = null;
</script>

<Tabbar class="max-w-5xl bottom-0 fixed">
    <TabbarLink active={location!="home"} label="Home" href="/" />
    <TabbarLink active={location!="channels"} label="Channels" href="/channels/" />
    <TabbarLink active={location!="playlists"} label="Playlists" href="/playlists/" />
    <TabbarLink active={location!="tasks"} label="Tasks" href="/tasks/" />
</Tabbar>