package main

import (
	"flag"
	"yt-archive/tasks"
)

func main() {
	p := flag.String("path", ".", "directory to rebuild mpd")
	o := flag.String("output", "manifest.mpd", "output filename")
	h := flag.Bool("help", false, "print this help message")

	flag.Parse()

	if *h {
		flag.Usage()
		return
	}

	tasks.BuildManifest(*p, *o)
}
