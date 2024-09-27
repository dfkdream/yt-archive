import { precacheAndRoute } from "workbox-precaching";

precacheAndRoute(self.__WB_MANIFEST, {
    directoryIndex: "index.html",
    cleanURLs: true,
});
