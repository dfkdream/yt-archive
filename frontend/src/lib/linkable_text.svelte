<script lang="ts">
    import { onMount } from "svelte";

    export let renderExternalLinks = false;

    let div: HTMLDivElement;

    onMount(() => {
        let text = div.innerText;

        text = text
            .replaceAll("&", "&amp;")
            .replaceAll("<", "&lt;")
            .replaceAll(">", "&gt;")
            .replaceAll('"', "&quot;")
            .replaceAll("'", "&#x27;");

        if (renderExternalLinks) {
            text = text.replaceAll(
                /\b((https?|ftp):\/\/[-\w@:%_\+.~#?&//=]*[-\w@:%_\+~#?&//=])\b/g,
                '<a href="$1">$1</a>',
            );
        }

        text = text.replaceAll(/(@\w+)/g, '<a href="/channels/$1">$1</a>');
        const matches = text.matchAll(/\b([0-9]{1,2}):([0-5][0-9])\b/g);
        for (const m of matches) {
            const sec = parseInt(m[1]) * 60 + parseInt(m[2]);
            text = text.replace(m[0], `<a href="?t=${sec}">${m[0]}</a>`);
        }

        text = text.replaceAll("\n", "<br>");

        div.innerHTML = text;
    });
</script>

<div {...$$restProps} bind:this={div}><pre><slot /></pre></div>
