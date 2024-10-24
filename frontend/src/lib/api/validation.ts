export function isValidVideoID(id: string) {
    return /^[A-Za-z0-9_-]{11}$/.test(id);
}

export function isValidPlaylistID(id: string) {
    return /^PL[A-Za-z0-9_-]{32}$/.test(id);
}
