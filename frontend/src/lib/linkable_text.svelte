<script lang="ts">
    export let renderExternalLinks = false;
    export let text;
    export let videoId;

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

    text = text.replaceAll(/(@\w+)/g, '<a href="/channels?id=$1">$1</a>');

    const matches = text.matchAll(/\b([0-9]{1,2}):([0-5][0-9])\b/g);
    for (const m of matches) {
        const sec = parseInt(m[1]) * 60 + parseInt(m[2]);

        text = text.replace(
            m[0],
            `<a href="?id=${videoId}&t=${sec}" data-sveltekit-preload-data="false">${m[0]}</a>`,
        );
    }

    text = text.replaceAll("\n", "<br>");
</script>

<div {...$$restProps}>{@html text}</div>
