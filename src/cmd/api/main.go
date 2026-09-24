// Command api runs the HTTP server that exposes the media pipeline.
//
// It must stay thin: load config, build dependencies, start the server, and
// handle OS signals for graceful shutdown. No media-processing logic here.
package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "api: error:", err)
		os.Exit(1)
	}
}

// run holds the real (testable) server logic. Right now it does nothing.
//
// TODO (story 11+): construct the app service and HTTP handlers, start the
// server, and register graceful shutdown (story 12).
func run() error {
	fmt.Println("api: not implemented yet")
	return nil
}
