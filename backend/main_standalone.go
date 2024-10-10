//go:build standalone

package main

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed dist/*
var distFS embed.FS

func main() {
	log.Println("Running standalone build of yt-archive.")

	subFS, err := fs.Sub(distFS, "dist")
	if err != nil {
		log.Fatal(err)
	}

	entrypoint(subFS)
}
