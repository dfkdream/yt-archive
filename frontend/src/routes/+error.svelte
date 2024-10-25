<script lang="ts">
    import { page } from "$app/stores";
    import Navbar from "$lib/components/navbar.svelte";
</script>

<Navbar />

<h1 class="font-serif text-9xl text-center mt-4">{$page.status}</h1>
<pre class="text-center">{$page.error?.message}</pre>
