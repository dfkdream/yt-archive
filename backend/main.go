//go:build !standalone

package main

import (
	"os"
)

func main() {
	entrypoint(os.DirFS("dist"))
}
