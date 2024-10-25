<script lang="ts">
    export let renderExternalLinks = false;
    export let text: string;
    export let videoId: string;

    function render(input: string) {
        input = input
            .replaceAll("&", "&amp;")
            .replaceAll("<", "&lt;")
            .replaceAll(">", "&gt;")
            .replaceAll('"', "&quot;")
            .replaceAll("'", "&#x27;");

        if (renderExternalLinks) {
            input = input.replaceAll(
                /\b((https?|ftp):\/\/[-\w@:%_\+.~#?&//=]*[-\w@:%_\+~#?&//=])\b/g,
                '<a href="$1">$1</a>',
            );
        }

        input = input.replaceAll(/(@\w+)/g, '<a href="/channels?id=$1">$1</a>');

        const matches = input.matchAll(/\b([0-9]{1,2}):([0-5][0-9])\b/g);
        for (const m of matches) {
            const sec = parseInt(m[1]) * 60 + parseInt(m[2]);

            input = input.replace(
                m[0],
                `<a href="?id=${videoId}&t=${sec}" data-sveltekit-preload-data="false">${m[0]}</a>`,
            );
        }

        text = input.replaceAll("\n", "<br>");
    }

    $: render(text);
</script>

<div {...$$restProps}>{@html text}</div>
