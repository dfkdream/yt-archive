import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { SvelteKitPWA } from '@vite-pwa/sveltekit';

export default defineConfig({
	server: {
		proxy: {
			'^/api/.*': 'http://localhost:8080'
		}
	},
	plugins: [
		sveltekit(),
		SvelteKitPWA({
			strategies: "injectManifest",
			srcDir: "src",
			filename: "service-worker.js",
			injectRegister: null,
			manifest: {
				name: "yt-archive",
				short_name: "yt-archive",
				description: "yt-archive",
				background_color: "rgb(247 247 248)",
				theme_color: "rgb(247 247 248)",
				lang: "en-US",
				icons: [
					{
						src: "favicon.png",
						sizes: "128x128",
						type: "image/png",
						purpose: "any maskable"
					}
				]
			},
		})
	]
});
