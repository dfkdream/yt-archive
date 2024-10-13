export function load({ params, url }) {
    return {
        id: params.id,
        startTime: parseInt(url.searchParams.get('t') || "0")
    };
}
