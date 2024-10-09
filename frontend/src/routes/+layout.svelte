<script>
    import "../app.css";
    import { App, Page } from "konsta/svelte";
    import { pwaInfo } from "virtual:pwa-info";

    $: webManifestLink = pwaInfo ? pwaInfo.webManifest.linkTag : "";
</script>

<svelte:head>
    {@html webManifestLink}
</svelte:head>

<App theme="ios">
    <Page class="max-w-5xl m-auto">
        <slot />
    </Page>
</App>
